package access

import (
	"bytes"
	"encoding/json"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/usc-core/client"
	"github.com/tv42/compound"
)

// compound index types
type ssb struct {
	A string
	B string
	C []byte
}

func MakeBuckets(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Indexes"))
		_, err = tx.CreateBucketIfNotExists([]byte("Channels"))
		_, err = tx.CreateBucketIfNotExists([]byte("Judges"))
		_, err = tx.CreateBucketIfNotExists([]byte("MyAccounts"))
		_, err = tx.CreateBucketIfNotExists([]byte("TheirAccounts"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func SetJudge(tx *bolt.Tx, jd *core.Judge) error {
	b, err := json.Marshal(jd)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("Judges")).Put(jd.Pubkey, b)
	if err != nil {
		return err
	}

	return nil
}

func SetMyAccount(tx *bolt.Tx, ma *core.MyAccount) error {
	b, err := json.Marshal(ma)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("MyAccounts")).Put([]byte(ma.Pubkey), b)
	if err != nil {
		return err
	}

	// Relations

	b, err = json.Marshal(ma.Judge)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("Judges")).Put(ma.Judge.Pubkey, b)
	if err != nil {
		return err
	}

	return nil
}

func GetMyAccount(tx *bolt.Tx, key []byte) (*core.MyAccount, error) {
	ma := &core.MyAccount{}
	err := json.Unmarshal(tx.Bucket([]byte("MyAccounts")).Get(key), ma)
	if err != nil {
		return nil, err
	}
	err = PopulateMyAccount(tx, ma)
	if err != nil {
		return nil, err
	}
	return ma, nil
}

func PopulateMyAccount(tx *bolt.Tx, ma *core.MyAccount) error {
	jd := &core.Judge{}
	err := json.Unmarshal(tx.Bucket([]byte("Judges")).Get([]byte(ma.Judge.Pubkey)), jd)
	if err != nil {
		return err
	}
	ma.Judge = jd

	return nil
}

func SetTheirAccount(tx *bolt.Tx, ta *core.TheirAccount) error {
	b, err := json.Marshal(ta)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("TheirAccounts")).Put([]byte(ta.Pubkey), b)
	if err != nil {
		return err
	}

	// Relations

	b, err = json.Marshal(ta.Judge)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("Judges")).Put(ta.Judge.Pubkey, b)
	if err != nil {
		return err
	}

	return nil
}

func GetTheirAccount(tx *bolt.Tx, key []byte) (*core.TheirAccount, error) {
	ta := &core.TheirAccount{}
	err := json.Unmarshal(tx.Bucket([]byte("TheirAccounts")).Get(key), ta)
	if err != nil {
		return nil, err
	}
	err = PopulateTheirAccount(tx, ta)
	if err != nil {
		return nil, err
	}
	return ta, nil
}

func PopulateTheirAccount(tx *bolt.Tx, ta *core.TheirAccount) error {
	jd := &core.Judge{}
	err := json.Unmarshal(tx.Bucket([]byte("Judges")).Get([]byte(ta.Judge.Pubkey)), jd)
	if err != nil {
		return err
	}

	ta.Judge = jd

	return nil
}

func SetChannel(tx *bolt.Tx, ch *core.Channel) error {
	b, err := json.Marshal(ch)
	if err != nil {
		return err
	}
	err = tx.Bucket([]byte("Channels")).Put([]byte(ch.ChannelId), b)
	if err != nil {
		return err
	}

	// Relations

	// Escrow Provider

	b, err = json.Marshal(ch.Judge)
	if err != nil {
		return err
	}

	tx.Bucket([]byte("Judges")).Put(ch.Judge.Pubkey, b)

	// My Account

	b, err = json.Marshal(ch.MyAccount)
	if err != nil {
		return err
	}

	tx.Bucket([]byte("MyAccounts")).Put(ch.MyAccount.Pubkey, b)

	// Their Account

	b, err = json.Marshal(ch.TheirAccount)
	if err != nil {
		return err
	}

	tx.Bucket([]byte("TheirAccounts")).Put(ch.TheirAccount.Pubkey, b)

	// Indexes

	// Escrow Provider Pubkey

	err = tx.Bucket([]byte("Indexes")).Put(compound.Key(ssb{
		"Judge",
		"Pubkey",
		ch.Judge.Pubkey}), []byte(ch.ChannelId))
	if err != nil {
		return err
	}

	return nil
}

func GetChannel(tx *bolt.Tx, key string) (*core.Channel, error) {
	b := tx.Bucket([]byte("Channels")).Get([]byte(key))
	if bytes.Compare(b, []byte{}) == 0 {
		return nil, nil
	}

	ch := &core.Channel{}
	err := json.Unmarshal(b, ch)
	if err != nil {
		return nil, err
	}
	err = PopulateChannel(tx, ch)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func PopulateChannel(tx *bolt.Tx, ch *core.Channel) error {
	ma := &core.MyAccount{}
	err := json.Unmarshal(tx.Bucket([]byte("MyAccounts")).Get([]byte(ch.MyAccount.Pubkey)), ma)
	if err != nil {
		return err
	}
	err = PopulateMyAccount(tx, ma)
	if err != nil {
		return err
	}

	ta := &core.TheirAccount{}
	err = json.Unmarshal(tx.Bucket([]byte("TheirAccounts")).Get([]byte(ch.TheirAccount.Pubkey)), ta)
	if err != nil {
		return err
	}
	err = PopulateTheirAccount(tx, ta)
	if err != nil {
		return err
	}

	jd := &core.Judge{}
	err = json.Unmarshal(tx.Bucket([]byte("Judges")).Get([]byte(ch.Judge.Pubkey)), jd)
	if err != nil {
		return err
	}

	ch.MyAccount = ma
	ch.TheirAccount = ta
	ch.Judge = jd

	return nil
}

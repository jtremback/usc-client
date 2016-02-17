package access

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/usc-core/peer"
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
		_, err = tx.CreateBucketIfNotExists([]byte("Accounts"))
		_, err = tx.CreateBucketIfNotExists([]byte("Counterparties"))
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

func GetJudge(tx *bolt.Tx, key []byte) (*core.Judge, error) {
	jd := &core.Judge{}
	err := json.Unmarshal(tx.Bucket([]byte("Accounts")).Get(key), jd)
	if err != nil {
		return nil, errors.New("database error")
	}
	if jd == nil {
		return nil, errors.New("judge not found")
	}

	return jd, nil
}

func SetAccount(tx *bolt.Tx, acct *core.Account) error {
	b, err := json.Marshal(acct)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("Accounts")).Put([]byte(acct.Pubkey), b)
	if err != nil {
		return err
	}

	// Relations

	b, err = json.Marshal(acct.Judge)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("Judges")).Put(acct.Judge.Pubkey, b)
	if err != nil {
		return err
	}

	return nil
}

func GetAccount(tx *bolt.Tx, key []byte) (*core.Account, error) {
	acct := &core.Account{}
	err := json.Unmarshal(tx.Bucket([]byte("Accounts")).Get(key), acct)
	if err != nil {
		return nil, errors.New("database error")
	}
	if acct == nil {
		return nil, errors.New("account not found")
	}
	err = PopulateAccount(tx, acct)
	if err != nil {
		return nil, errors.New("database error")
	}
	return acct, nil
}

func PopulateAccount(tx *bolt.Tx, acct *core.Account) error {
	jd := &core.Judge{}
	err := json.Unmarshal(tx.Bucket([]byte("Judges")).Get([]byte(acct.Judge.Pubkey)), jd)
	if err != nil {
		return err
	}
	acct.Judge = jd

	return nil
}

func SetCounterparty(tx *bolt.Tx, cpt *core.Counterparty) error {
	b, err := json.Marshal(cpt)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("Counterparties")).Put([]byte(cpt.Pubkey), b)
	if err != nil {
		return err
	}

	// Relations

	b, err = json.Marshal(cpt.Judge)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("Judges")).Put(cpt.Judge.Pubkey, b)
	if err != nil {
		return err
	}

	return nil
}

func GetCounterparty(tx *bolt.Tx, key []byte) (*core.Counterparty, error) {
	cpt := &core.Counterparty{}
	err := json.Unmarshal(tx.Bucket([]byte("Counterparties")).Get(key), cpt)
	if err != nil {
		return nil, errors.New("database error")
	}
	if cpt == nil {
		return nil, errors.New("account not found")
	}
	err = PopulateCounterparty(tx, cpt)
	if err != nil {
		return nil, errors.New("database error")
	}
	return cpt, nil
}

func PopulateCounterparty(tx *bolt.Tx, cpt *core.Counterparty) error {
	jd := &core.Judge{}
	err := json.Unmarshal(tx.Bucket([]byte("Judges")).Get([]byte(cpt.Judge.Pubkey)), jd)
	if err != nil {
		return err
	}

	cpt.Judge = jd

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

	b, err = json.Marshal(ch.Account)
	if err != nil {
		return err
	}

	tx.Bucket([]byte("Accounts")).Put(ch.Account.Pubkey, b)

	// Their Account

	b, err = json.Marshal(ch.Counterparty)
	if err != nil {
		return err
	}

	tx.Bucket([]byte("Counterparties")).Put(ch.Counterparty.Pubkey, b)

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
	ch := &core.Channel{}
	err := json.Unmarshal(tx.Bucket([]byte("Channels")).Get([]byte(key)), ch)
	if err != nil {
		return nil, errors.New("database error")
	}
	if ch == nil {
		return nil, errors.New("channel not found")
	}
	err = PopulateChannel(tx, ch)
	if err != nil {
		return nil, errors.New("database error")
	}
	return ch, nil
}

func GetProposedChannels(tx *bolt.Tx) ([]*core.Channel, error) {
	var err error
	chs := []*core.Channel{}
	i := 0
	err = tx.Bucket([]byte("Channels")).ForEach(func(k, v []byte) error {
		ch := &core.Channel{}
		err = json.Unmarshal(v, ch)
		if err != nil {
			return err
		}
		if ch.Phase == core.PENDING_OPEN && len(ch.OpeningTxEnvelope.Signatures) == 1 {
			chs[i] = ch
			i++
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("database error")
	}
	return chs, nil
}

func GetProposedUpdateTxs(tx *bolt.Tx) ([]*core.Channel, error) {
	var err error
	chs := []*core.Channel{}
	i := 0
	err = tx.Bucket([]byte("Channels")).ForEach(func(k, v []byte) error {
		ch := &core.Channel{}
		err = json.Unmarshal(v, ch)
		if err != nil {
			return err
		}
		if bytes.Compare(ch.ProposedUpdateTxEnvelope.Signatures[ch.Me], []byte{}) == 0 {
			chs[i] = ch
			i++
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("database error")
	}
	return chs, nil
}

func PopulateChannel(tx *bolt.Tx, ch *core.Channel) error {
	acct := &core.Account{}
	err := json.Unmarshal(tx.Bucket([]byte("Accounts")).Get([]byte(ch.Account.Pubkey)), acct)
	if err != nil {
		return err
	}
	err = PopulateAccount(tx, acct)
	if err != nil {
		return err
	}

	cpt := &core.Counterparty{}
	err = json.Unmarshal(tx.Bucket([]byte("Counterparties")).Get([]byte(ch.Counterparty.Pubkey)), cpt)
	if err != nil {
		return err
	}
	err = PopulateCounterparty(tx, cpt)
	if err != nil {
		return err
	}

	jd := &core.Judge{}
	err = json.Unmarshal(tx.Bucket([]byte("Judges")).Get([]byte(ch.Judge.Pubkey)), jd)
	if err != nil {
		return err
	}

	ch.Account = acct
	ch.Counterparty = cpt
	ch.Judge = jd

	return nil
}

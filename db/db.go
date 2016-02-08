package db

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/upc-core/wallet"
	"github.com/tv42/compound"
)

// compound index types
type ssb struct {
	b string
	c string
	d []byte
}

func makeBuckets(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Channels"))
		_, err = tx.CreateBucketIfNotExists([]byte("EscrowProviders"))
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

func setEscrowProvider(tx *bolt.Tx, ep *core.EscrowProvider) error {
	b, err := json.Marshal(ep)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("EscrowProviders")).Put(ep.Pubkey, b)
	if err != nil {
		return err
	}

	return nil
}

func setMyAccount(tx *bolt.Tx, ma *core.MyAccount) error {
	b, err := json.Marshal(ma)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("MyAccounts")).Put([]byte(ma.Pubkey), b)
	if err != nil {
		return err
	}

	return nil
}

func populateMyAccount(tx *bolt.Tx, ma *core.MyAccount) error {
	var ep *core.EscrowProvider
	err := json.Unmarshal(tx.Bucket([]byte("EscrowProviders")).Get([]byte(ma.EscrowProvider.Pubkey)), ep)
	if err != nil {
		return err
	}

	ma.EscrowProvider = ep

	return nil
}

func setTheirAccount(tx *bolt.Tx, ta *core.TheirAccount) error {
	b, err := json.Marshal(ta)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("TheirAccounts")).Put([]byte(ta.Pubkey), b)
	if err != nil {
		return err
	}

	return nil

}

func populateTheirAccount(tx *bolt.Tx, ta *core.TheirAccount) error {
	var ep *core.EscrowProvider
	err := json.Unmarshal(tx.Bucket([]byte("EscrowProviders")).Get([]byte(ta.EscrowProvider.Pubkey)), ep)
	if err != nil {
		return err
	}

	ta.EscrowProvider = ep

	return nil
}

func setChannel(tx *bolt.Tx, ch *core.Channel) error {
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

	b, err = json.Marshal(ch.EscrowProvider)
	if err != nil {
		return err
	}

	tx.Bucket([]byte("EscrowProviders")).Put(ch.EscrowProvider.Pubkey, b)

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
		"EscrowProvider",
		"Pubkey",
		ch.EscrowProvider.Pubkey}), []byte(ch.ChannelId))
	if err != nil {
		return err
	}

	return nil
}

func populateChannel(tx *bolt.Tx, ch *core.Channel) error {
	var ma *core.MyAccount
	err := json.Unmarshal(tx.Bucket([]byte("MyAccounts")).Get([]byte(ch.MyAccount.Pubkey)), ma)
	if err != nil {
		return err
	}
	err = populateMyAccount(tx, ma)
	if err != nil {
		return err
	}

	var ta *core.TheirAccount
	err = json.Unmarshal(tx.Bucket([]byte("TheirAccounts")).Get([]byte(ch.TheirAccount.Pubkey)), ta)
	if err != nil {
		return err
	}
	err = populateTheirAccount(tx, ta)
	if err != nil {
		return err
	}

	var ep *core.EscrowProvider
	err = json.Unmarshal(tx.Bucket([]byte("EscrowProviders")).Get([]byte(ch.EscrowProvider.Pubkey)), ep)
	if err != nil {
		return err
	}

	ch.MyAccount = ma
	ch.TheirAccount = ta
	ch.EscrowProvider = ep

	return nil
}

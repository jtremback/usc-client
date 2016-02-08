package db

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/tv42/compound"
)

func setEscrowProvider(db *bolt.DB, ep *core.EscrowProvider) error {
	err := db.Update(func(tx *bolt.Tx) error {

		b, err := json.Marshal(ep)
		if err != nil {
			return err
		}

		err = tx.Bucket([]byte("EscrowProviders")).Put([]byte(ep.Pubkey), b)
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

func setMyAccount(db *bolt.DB, ma *core.MyAccount) error {
	err := db.Update(func(tx *bolt.Tx) error {

		b, err := json.Marshal(ma)
		if err != nil {
			return err
		}

		err = tx.Bucket([]byte("EscrowProviders")).Put([]byte(ma.Pubkey), b)
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

func populateMyAccount(db *bolt.DB, ma *core.MyAccount) error {
	err := db.View(func(tx *bolt.Tx) error {

		var ep *core.EscrowProvider
		err := json.Unmarshal(tx.Bucket([]byte("EscrowProviders")).Get([]byte(ma.EscrowProvider.Pubkey)), ep)
		if err != nil {
			return err
		}

		ma.EscrowProvider = ep

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func setTheirAccount(db *bolt.DB, ta *core.TheirAccount) error {
	err := db.Update(func(tx *bolt.Tx) error {

		b, err := json.Marshal(ta)
		if err != nil {
			return err
		}

		err = tx.Bucket([]byte("EscrowProviders")).Put([]byte(ta.Pubkey), b)
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

func populateTheirAccount(db *bolt.DB, ta *core.MyAccount) error {
	err := db.View(func(tx *bolt.Tx) error {

		var ep *core.EscrowProvider
		err := json.Unmarshal(tx.Bucket([]byte("EscrowProviders")).Get([]byte(ta.EscrowProvider.Pubkey)), ep)
		if err != nil {
			return err
		}

		ta.EscrowProvider = ep

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func setChannel(db *bolt.DB, ch *core.Channel) error {
	err := db.Update(func(tx *bolt.Tx) error {

		b, err := json.Marshal(ch)
		if err != nil {
			return err
		}

		err = tx.Bucket([]byte("Channels")).Put([]byte(ch.ChannelId), b)
		if err != nil {
			return err
		}

		// Relations

		b, err = json.Marshal(ch.EscrowProvider)
		if err != nil {
			return err
		}

		tx.Bucket([]byte("EscrowProviders")).Put(ch.EscrowProvider.Pubkey, b)

		b, err = json.Marshal(ch.MyAccount)
		if err != nil {
			return err
		}

		tx.Bucket([]byte("MyAccounts")).Put(ch.MyAccount.Pubkey, b)

		b, err = json.Marshal(ch.TheirAccount)
		if err != nil {
			return err
		}

		tx.Bucket([]byte("TheirAccounts")).Put(ch.TheirAccount.Pubkey, b)

		// Indexes

		err = tx.Bucket([]byte("Indexes")).Put(compound.Key(ssb{
			"EscrowProvider",
			"Pubkey",
			ch.EscrowProvider.Pubkey}), []byte(ch.ChannelId))
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

func populateChannel(db *bolt.DB, ch *core.Channel) error {
	err := db.View(func(tx *bolt.Tx) error {
		var ma *core.MyAccount
		err := json.Unmarshal(tx.Bucket([]byte("MyAccounts")).Get([]byte(ch.MyAccount.Pubkey)), ma)
		if err != nil {
			return err
		}

		var ta *core.TheirAccount
		err = json.Unmarshal(tx.Bucket([]byte("TheirAccounts")).Get([]byte(ch.TheirAccount.Pubkey)), ta)
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
	})
	if err != nil {
		return err
	}
	return nil
}

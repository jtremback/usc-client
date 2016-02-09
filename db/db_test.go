package db

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/upc-core/wallet"
	"github.com/jtremback/upc-core/wire"
)

func TestSetEscrowProvider(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = makeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ep := &core.EscrowProvider{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}
	ep2 := &core.EscrowProvider{}

	db.Update(func(tx *bolt.Tx) error {
		err := setEscrowProvider(tx, ep)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		err := json.Unmarshal(tx.Bucket([]byte("EscrowProviders")).Get(ep.Pubkey), ep2)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	if !reflect.DeepEqual(ep, ep2) {
		t.Fatal("structs not equal :(")
	}
}

func TestSetMyAccount(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = makeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ma := &core.MyAccount{
		Name:    "boogie",
		Privkey: []byte{30, 30, 30},
		Pubkey:  []byte{40, 40, 40},
		EscrowProvider: &core.EscrowProvider{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := setMyAccount(tx, ma)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		ma2 := &core.MyAccount{}
		json.Unmarshal(tx.Bucket([]byte("MyAccounts")).Get(ma.Pubkey), ma2)
		if !reflect.DeepEqual(ma, ma2) {
			t.Fatal("MyAccount incorrect")
		}

		fromDB := tx.Bucket([]byte("EscrowProviders")).Get(ma.EscrowProvider.Pubkey)
		ep := &core.EscrowProvider{}
		json.Unmarshal(fromDB, ep)

		if !reflect.DeepEqual(ma.EscrowProvider, ep) {
			t.Fatal("EscrowProvider incorrect", ma.EscrowProvider, ep, string(tx.Bucket([]byte("EscrowProviders")).Get(ma.EscrowProvider.Pubkey)))
		}
		return nil
	})
}

func TestPopulateMyAccount(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = makeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ma := &core.MyAccount{
		Name:    "boogie",
		Privkey: []byte{30, 30, 30},
		Pubkey:  []byte{40, 40, 40},
		EscrowProvider: &core.EscrowProvider{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	ep := &core.EscrowProvider{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}

	db.Update(func(tx *bolt.Tx) error {
		err := setMyAccount(tx, ma)
		if err != nil {
			t.Fatal(err)
		}

		err = setEscrowProvider(tx, ep)
		if err != nil {
			t.Fatal(err)
		}

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		err := populateMyAccount(tx, ma)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(ma.EscrowProvider, ep) {
			t.Fatal("EscrowProvider incorrect", ma.EscrowProvider, ep)
		}
		return nil
	})
}

func TestSetTheirAccount(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = makeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ta := &core.TheirAccount{
		Name:   "boogie",
		Pubkey: []byte{40, 40, 40},
		EscrowProvider: &core.EscrowProvider{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := setTheirAccount(tx, ta)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		ta2 := &core.TheirAccount{}
		json.Unmarshal(tx.Bucket([]byte("TheirAccounts")).Get(ta.Pubkey), ta2)
		if !reflect.DeepEqual(ta, ta2) {
			t.Fatal("TheirAccount incorrect")
		}

		fromDB := tx.Bucket([]byte("EscrowProviders")).Get(ta.EscrowProvider.Pubkey)
		ep := &core.EscrowProvider{}
		json.Unmarshal(fromDB, ep)

		if !reflect.DeepEqual(ta.EscrowProvider, ep) {
			t.Fatal("EscrowProvider incorrect", ta.EscrowProvider, ep, string(tx.Bucket([]byte("EscrowProviders")).Get(ta.EscrowProvider.Pubkey)))
		}
		return nil
	})
}

func TestPopulateTheirAccount(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = makeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ta := &core.TheirAccount{
		Name:   "boogie",
		Pubkey: []byte{40, 40, 40},
		EscrowProvider: &core.EscrowProvider{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	ep := &core.EscrowProvider{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}

	db.Update(func(tx *bolt.Tx) error {
		err := setTheirAccount(tx, ta)
		if err != nil {
			t.Fatal(err)
		}

		err = setEscrowProvider(tx, ep)
		if err != nil {
			t.Fatal(err)
		}

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		err := populateTheirAccount(tx, ta)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(ta.EscrowProvider, ep) {
			t.Fatal("EscrowProvider incorrect", ta.EscrowProvider, ep)
		}
		return nil
	})
}

func TestSetChannel(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = makeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ch := &core.Channel{
		ChannelId: "xyz23",
		Phase:     2,

		OpeningTx:         &wire.OpeningTx{},
		OpeningTxEnvelope: &wire.Envelope{},

		LastUpdateTx:         &wire.UpdateTx{},
		LastUpdateTxEnvelope: &wire.Envelope{},

		LastFullUpdateTx:         &wire.UpdateTx{},
		LastFullUpdateTxEnvelope: &wire.Envelope{},

		Me:           0,
		Fulfillments: [][]byte{[]byte{80, 80}},

		EscrowProvider: &core.EscrowProvider{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},

		MyAccount: &core.MyAccount{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Privkey: []byte{40, 40, 40},
			EscrowProvider: &core.EscrowProvider{
				Name:    "wrong",
				Pubkey:  []byte{40, 40, 40},
				Address: "stoops.com:3004",
			},
		},

		TheirAccount: &core.TheirAccount{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
			EscrowProvider: &core.EscrowProvider{
				Name:    "wrong",
				Pubkey:  []byte{40, 40, 40},
				Address: "stoops.com:3004",
			},
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := setChannel(tx, ch)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		ch2 := &core.Channel{}
		json.Unmarshal(tx.Bucket([]byte("Channels")).Get([]byte(ch.ChannelId)), ch2)
		if !reflect.DeepEqual(ch, ch2) {
			t.Fatal("Channel incorrect")
		}
		epJson := tx.Bucket([]byte("EscrowProviders")).Get(ch.EscrowProvider.Pubkey)
		ep := &core.EscrowProvider{}
		json.Unmarshal(epJson, ep)

		if !reflect.DeepEqual(ch.EscrowProvider, ep) {
			t.Fatal("EscrowProvider incorrect")
		}

		maJson := tx.Bucket([]byte("MyAccounts")).Get(ch.MyAccount.Pubkey)
		ma := &core.MyAccount{}
		json.Unmarshal(maJson, ma)

		if !reflect.DeepEqual(ch.MyAccount, ma) {
			t.Fatal("MyAccount incorrect")
		}

		taJson := tx.Bucket([]byte("TheirAccounts")).Get(ch.TheirAccount.Pubkey)
		ta := &core.TheirAccount{}
		json.Unmarshal(taJson, ta)

		if !reflect.DeepEqual(ch.TheirAccount, ta) {
			t.Fatal("TheirAccount incorrect")
		}
		return nil
	})
}

func TestPopulateChannel(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = makeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ch := &core.Channel{
		ChannelId: "xyz23",
		Phase:     2,

		OpeningTx:         &wire.OpeningTx{},
		OpeningTxEnvelope: &wire.Envelope{},

		LastUpdateTx:         &wire.UpdateTx{},
		LastUpdateTxEnvelope: &wire.Envelope{},

		LastFullUpdateTx:         &wire.UpdateTx{},
		LastFullUpdateTxEnvelope: &wire.Envelope{},

		Me:           0,
		Fulfillments: [][]byte{[]byte{80, 80}},

		EscrowProvider: &core.EscrowProvider{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},

		MyAccount: &core.MyAccount{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Privkey: []byte{40, 40, 40},
			EscrowProvider: &core.EscrowProvider{
				Name:    "wrong",
				Pubkey:  []byte{40, 40, 40},
				Address: "stoops.com:3004",
			},
		},

		TheirAccount: &core.TheirAccount{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
			EscrowProvider: &core.EscrowProvider{
				Name:    "wrong",
				Pubkey:  []byte{40, 40, 40},
				Address: "stoops.com:3004",
			},
		},
	}

	ep := &core.EscrowProvider{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}

	ma := &core.MyAccount{
		Name:    "crow",
		Pubkey:  []byte{40, 40, 40},
		Privkey: []byte{40, 40, 40},
		EscrowProvider: &core.EscrowProvider{
			Name:    "joe",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	ta := &core.TheirAccount{
		Name:    "flerb",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
		EscrowProvider: &core.EscrowProvider{
			Name:    "joe",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := setChannel(tx, ch)
		if err != nil {
			t.Fatal(err)
		}

		err = setMyAccount(tx, ma)
		if err != nil {
			t.Fatal(err)
		}

		err = setEscrowProvider(tx, ep)
		if err != nil {
			t.Fatal(err)
		}

		err = setTheirAccount(tx, ta)
		if err != nil {
			t.Fatal(err)
		}

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		err = populateChannel(tx, ch)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(ch.EscrowProvider, ep) {
			t.Fatal("EscrowProvider incorrect", ta.EscrowProvider, ep)
		}

		if !reflect.DeepEqual(ch.MyAccount, ma) {
			t.Fatal("MyAccount incorrect", ch.MyAccount, ma)
		}

		if !reflect.DeepEqual(ch.TheirAccount, ta) {
			t.Fatal("TheirAccount incorrect", ch.TheirAccount, ta)
		}

		return nil
	})
}

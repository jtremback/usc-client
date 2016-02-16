package access

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/usc-core/client"
	"github.com/jtremback/usc-core/wire"
)

func TestSetJudge(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ju := &core.Judge{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}
	ju2 := &core.Judge{}

	db.Update(func(tx *bolt.Tx) error {
		err := SetJudge(tx, ju)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		fmt.Println(string(tx.Bucket([]byte("Judges")).Get(ju.Pubkey)))
		err := json.Unmarshal(tx.Bucket([]byte("Judges")).Get(ju.Pubkey), ju2)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	if !reflect.DeepEqual(ju, ju2) {
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

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ma := &core.MyAccount{
		Name:    "boogie",
		Privkey: []byte{30, 30, 30},
		Pubkey:  []byte{40, 40, 40},
		Judge: &core.Judge{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := SetMyAccount(tx, ma)
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

		fromDB := tx.Bucket([]byte("Judges")).Get(ma.Judge.Pubkey)
		ju := &core.Judge{}
		json.Unmarshal(fromDB, ju)

		if !reflect.DeepEqual(ma.Judge, ju) {
			t.Fatal("Judge incorrect", ma.Judge, ju, string(tx.Bucket([]byte("Judges")).Get(ma.Judge.Pubkey)))
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

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ma := &core.MyAccount{
		Name:    "boogie",
		Privkey: []byte{30, 30, 30},
		Pubkey:  []byte{40, 40, 40},
		Judge: &core.Judge{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	ju := &core.Judge{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}

	db.Update(func(tx *bolt.Tx) error {
		err := SetMyAccount(tx, ma)
		if err != nil {
			t.Fatal(err)
		}

		err = SetJudge(tx, ju)
		if err != nil {
			t.Fatal(err)
		}

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		err := PopulateMyAccount(tx, ma)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(ma.Judge, ju) {
			t.Fatal("Judge incorrect", ma.Judge, ju)
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

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ta := &core.TheirAccount{
		Name:   "boogie",
		Pubkey: []byte{40, 40, 40},
		Judge: &core.Judge{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := SetTheirAccount(tx, ta)
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

		fromDB := tx.Bucket([]byte("Judges")).Get(ta.Judge.Pubkey)
		ju := &core.Judge{}
		json.Unmarshal(fromDB, ju)

		if !reflect.DeepEqual(ta.Judge, ju) {
			t.Fatal("Judge incorrect", ta.Judge, ju, string(tx.Bucket([]byte("Judges")).Get(ta.Judge.Pubkey)))
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

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ta := &core.TheirAccount{
		Name:   "boogie",
		Pubkey: []byte{40, 40, 40},
		Judge: &core.Judge{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	ju := &core.Judge{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}

	db.Update(func(tx *bolt.Tx) error {
		err := SetTheirAccount(tx, ta)
		if err != nil {
			t.Fatal(err)
		}

		err = SetJudge(tx, ju)
		if err != nil {
			t.Fatal(err)
		}

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		err := PopulateTheirAccount(tx, ta)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(ta.Judge, ju) {
			t.Fatal("Judge incorrect", ta.Judge, ju)
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

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ch := &core.Channel{
		ChannelId: "xyz23",
		Phase:     2,

		OpeningTx:         &wire.OpeningTx{},
		OpeningTxEnvelope: &wire.Envelope{},

		ProposedUpdateTx:         &wire.UpdateTx{},
		ProposedUpdateTxEnvelope: &wire.Envelope{},

		LastFullUpdateTx:         &wire.UpdateTx{},
		LastFullUpdateTxEnvelope: &wire.Envelope{},

		Me:           0,
		Fulfillments: [][]byte{[]byte{80, 80}},

		Judge: &core.Judge{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},

		MyAccount: &core.MyAccount{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Privkey: []byte{40, 40, 40},
			Judge: &core.Judge{
				Name:    "wrong",
				Pubkey:  []byte{40, 40, 40},
				Address: "stoops.com:3004",
			},
		},

		TheirAccount: &core.TheirAccount{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
			Judge: &core.Judge{
				Name:    "wrong",
				Pubkey:  []byte{40, 40, 40},
				Address: "stoops.com:3004",
			},
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := SetChannel(tx, ch)
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
		juJson := tx.Bucket([]byte("Judges")).Get(ch.Judge.Pubkey)
		ju := &core.Judge{}
		json.Unmarshal(juJson, ju)

		if !reflect.DeepEqual(ch.Judge, ju) {
			t.Fatal("Judge incorrect")
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

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ch := &core.Channel{
		ChannelId: "xyz23",
		Phase:     2,

		OpeningTx:         &wire.OpeningTx{},
		OpeningTxEnvelope: &wire.Envelope{},

		ProposedUpdateTx:         &wire.UpdateTx{},
		ProposedUpdateTxEnvelope: &wire.Envelope{},

		LastFullUpdateTx:         &wire.UpdateTx{},
		LastFullUpdateTxEnvelope: &wire.Envelope{},

		Me:           0,
		Fulfillments: [][]byte{[]byte{80, 80}},

		Judge: &core.Judge{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},

		MyAccount: &core.MyAccount{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Privkey: []byte{40, 40, 40},
			Judge: &core.Judge{
				Name:    "wrong",
				Pubkey:  []byte{40, 40, 40},
				Address: "stoops.com:3004",
			},
		},

		TheirAccount: &core.TheirAccount{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
			Judge: &core.Judge{
				Name:    "wrong",
				Pubkey:  []byte{40, 40, 40},
				Address: "stoops.com:3004",
			},
		},
	}

	ju := &core.Judge{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}

	ma := &core.MyAccount{
		Name:    "crow",
		Pubkey:  []byte{40, 40, 40},
		Privkey: []byte{40, 40, 40},
		Judge: &core.Judge{
			Name:    "joe",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	ta := &core.TheirAccount{
		Name:    "flerb",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
		Judge: &core.Judge{
			Name:    "joe",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := SetChannel(tx, ch)
		if err != nil {
			t.Fatal(err)
		}

		err = SetMyAccount(tx, ma)
		if err != nil {
			t.Fatal(err)
		}

		err = SetJudge(tx, ju)
		if err != nil {
			t.Fatal(err)
		}

		err = SetTheirAccount(tx, ta)
		if err != nil {
			t.Fatal(err)
		}

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		err = PopulateChannel(tx, ch)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(ch.Judge, ju) {
			t.Fatal("Judge incorrect", ch.Judge, ju)
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

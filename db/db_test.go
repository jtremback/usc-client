package db

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/upc-core/wallet"
)

func TestEscrowProvider(t *testing.T) {
	ep := &core.EscrowProvider{
		Name:    "joe",
		Privkey: []byte{30, 30, 30},
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}
	ep2 := &core.EscrowProvider{}

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

	db.Update(func(tx *bolt.Tx) error {
		setEscrowProvider(tx, ep)
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

}
func TestPopulateMyAccount(t *testing.T) {

}
func TestSetTheirAccount(t *testing.T) {

}
func TestPopulateTheirAccount(t *testing.T) {

}
func TestSetChannel(t *testing.T) {

}
func TestPopulateChannel(t *testing.T) {

}

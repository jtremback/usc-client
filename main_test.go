package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/jtremback/usc-client/access"
)

func TestAddJudge(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")
	err = access.MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	app := &api{db}
	req, err := http.NewRequest("GET", "http://localhost:3004", strings.NewReader(`{
    "name": "joe",
    "pubkey": "R5lVVs82M80i5OpR369StJqaHS61Ld+PzTCfS+0zyAA=",
    "address": "http://localhost:3401"
  }`))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	app.addJudge(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	db.View(func(tx *bolt.Tx) error {
		ref := []byte(`{"Name":"joe","Pubkey":"R5lVVs82M80i5OpR369StJqaHS61Ld+PzTCfS+0zyAA=","Address":"http://localhost:3401"}`)

		fromDB := tx.Bucket([]byte("Judges")).Get([]byte{40, 40, 40})

		if bytes.Compare(fromDB, ref) != 0 {
			t.Fatal("saved data not correct", string(fromDB), string(ref))
		}
		return nil
	})
}

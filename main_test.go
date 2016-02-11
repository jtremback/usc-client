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
    "name": "sffcu",
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
		ref := []byte(`{"Name":"sffcu","Pubkey":"R5lVVs82M80i5OpR369StJqaHS61Ld+PzTCfS+0zyAA=","Address":"http://localhost:3401"}`)

		fromDB := tx.Bucket([]byte("Judges")).Get([]byte{71, 153, 85, 86, 207, 54, 51, 205, 34, 228, 234, 81, 223, 175, 82, 180, 154, 154, 29, 46, 181, 45, 223, 143, 205, 48, 159, 75, 237, 51, 200, 0})

		if bytes.Compare(fromDB, ref) != 0 {
			t.Fatal("saved data not correct", string(fromDB), "shib", string(ref))
		}
		return nil
	})
}

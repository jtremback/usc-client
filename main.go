package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"database/sql"

	"github.com/jinzhu/gorm"
	core "github.com/jtremback/upc-core/wallet"
	_ "github.com/mattn/go-sqlite3"
)

type api struct {
	db *sql.DB
}

func main() {
	go http.ListenAndServe(":8120", nil)

	ticker := time.NewTicker(time.Second * 30)
	testTicker(time.Now())
	for t := range ticker.C {
		testTicker(t)
	}
}

func addAccountsTable(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	tx.Exec(`CREATE TABLE IF NOT EXISTS my_accounts (
	  name    				TEXT PRIMARY KEY,
	  pubkey  				BLOB,
	  privkey 				BLOB,
	  address 				TEXT,
	  escrow_provider TEXT
	);`)

	tx.Exec(`CREATE TABLE IF NOT EXISTS escrow_providers (
	  name 		TEXT PRIMARY KEY,
	  pubkey 	BLOB,
	  address TEXT
	);`)

	tx.Exec(`CREATE TABLE IF NOT EXISTS channels (
		ChannelId 							 TEXT PRIMARY KEY
		Phase

		OpeningTx         			 BLOB
		OpeningTxEnvelope 			 BLOB

		LastUpdateTx         		 BLOB
		LastUpdateTxEnvelope 		 BLOB

		LastFullUpdateTx         BLOB
		LastFullUpdateTxEnvelope BLOB

		EscrowProvider 					 TEXT
		Accounts 								 TEXT

		Me           						 INT
		Fulfillments 						 TEXT
	);`)
}

func testTicker(t time.Time) {
	fmt.Println("Tick at", t)
}

// func api() {
// 	http.HandleFunc("/test/", testHandler)

// 	http.HandleFunc("/v1/channels/", viewChannels)
// 	http.HandleFunc("/v1/channels/new/", newChannel)

// 	http.HandleFunc("/v1/accounts/", viewAccounts)
// 	http.HandleFunc("/v1/accounts/new/", newAccount)

// 	http.HandleFunc("/v1/escrow_providers/", viewEscrowProviders)
// 	http.HandleFunc("/v1/escrow_providers/new/", addEscrowProvider)

// 	http.HandleFunc("/v1/peers/", viewPeers)
// 	http.HandleFunc("/v1/peers/new/", addPeer)
// }

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println()
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func viewChannels(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func newChannel(w http.ResponseWriter, r *http.Request) {

}

func viewAccounts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func (a *api) newAccount(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		fmt.Println("no body")
		return
	}

	var reqData struct {
		Name           string
		Address        string
		EscrowProvider string
	}
	err := json.NewDecoder(r.Body).Decode(reqData)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		panic(err)
	}

	ep := &core.EscrowProvider{}
	db.First(ep, reqData.Name)

	acct, err := core.NewAccount(reqData.Name, reqData.Address, ep)
	if err != nil {
		panic(err)
	}

	db.NewRecord(acct)
	db.Create(&acct)

	json.NewEncoder(w).Encode(acct)
}

func viewEscrowProviders(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func addEscrowProvider(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		fmt.Println("no body")
		return
	}

	var ep core.EscrowProvider
	err := json.NewDecoder(r.Body).Decode(&ep)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		panic(err)
	}

	db.NewRecord(ep)
	db.Create(&ep)

	json.NewEncoder(w).Encode(ep)
}

func viewPeers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func addPeer(w http.ResponseWriter, r *http.Request) {
	// if r.Body == nil {
	// 	fmt.Println("no body")
	// 	return
	// }
	// dec := json.NewDecoder(r.Body)

	// var ep core.EscrowProvider
	// if err := dec.Decode(&ep); err != nil {
	// 	panic(err)
	// }

	// bytes, err := json.Marshal(ep)

	// db, err := bolt.Open(".", 0600, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// db.Update(func(tx *bolt.Tx) error {
	// 	indexes := tx.Bucket([]byte("Indexes"))
	// 	escrowProviders := tx.Bucket([]byte("EscrowProviders"))
	// 	err := escrowProviders.Put([]byte(ep.Name), bytes)
	// 	err = indexes.Put(makeKey("EscrowProviders", "Pubkey", string(ep.Pubkey)), ep.Name)
	// 	return err
	// })
}

func makeKey(s ...string) []byte {
	return []byte(strings.Join(append(s), "/"))
}

func (a *api) fail(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")

	data := struct {
		Error string
	}{Error: msg}

	resp, _ := json.Marshal(data)
	w.WriteHeader(status)
	w.Write(resp)
}

func (a *api) ok(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.fail(w, "oops something evil has happened", 500)
		return
	}
	w.Write(resp)
}

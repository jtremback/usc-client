package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"
	"github.com/tv42/compound"

	core "github.com/jtremback/upc-core/wallet"
	_ "github.com/mattn/go-sqlite3"
)

const (
	escrowProviders int = 1
	channels
)

type api struct {
	db *bolt.DB
}

// compound index types
type ssb struct {
	b string
	c string
	d []byte
}

func main() {
	go http.ListenAndServe(":8120", nil)

	ticker := time.NewTicker(time.Second * 30)
	testTicker(time.Now())
	for t := range ticker.C {
		testTicker(t)
	}
}

// func (a *api) getChannels(w http.ResponseWriter, r *http.Request) {
// 	rows, err := a.db.Query(`SELECT
// 		ChannelId,
// 		Phase,

// 		OpeningTx,
// 		OpeningTxEnvelope,

// 		LastUpdateTx,
// 		LastUpdateTxEnvelope,

// 		LastFullUpdateTx,
// 		LastFullUpdateTxEnvelope,

// 		EscrowProvider,
// 		Accounts,

// 		Me,
// 		Fulfillments
// 	FROM channels`)

// 	if err != nil {
// 		a.fail(w, err.Error(), 500)
// 		return
// 	}
// 	defer rows.Close()

// 	var channels []*core.Channel
// 	for rows.Next() {
// 		ch := core.Channel{}
// 		var accounts []byte
// 		var fulfillments []byte
// 		var openingTx []byte
// 		var openingTxEnvelope []byte
// 		var lastUpdateTx []byte
// 		var lastUpdateTxEnvelope []byte
// 		var lastFullUpdateTx []byte
// 		var lastFullUpdateTxEnvelope []byte
// 		var escrowProvider []byte
// 		err := rows.Scan(
// 			&ch.ChannelId,
// 			&ch.Phase,
// 			&openingTx,
// 			&openingTxEnvelope,
// 			&lastUpdateTx,
// 			&lastUpdateTxEnvelope,
// 			&lastFullUpdateTx,
// 			&lastFullUpdateTxEnvelope,
// 			&escrowProvider,
// 			&accounts,
// 			&ch.Me,
// 			&fulfillments,
// 		)
// 		if err != nil {
// 			a.fail(w, err.Error(), 500)
// 			return
// 		}

// 		err = json.Unmarshal(openingTx, ch.OpeningTx)
// 		err = json.Unmarshal(openingTx, ch.OpeningTx)

// 		var accts []string
// 		json.Unmarshal(accounts, accts)

// 		rows, err := a.db.Query(`SELECT * FROM accounts WHERE name IN ($1,$2)
// 														 INNER JOIN escrow_providers
// 														 ON accounts.EscrowProvider = escrow_providers.name`, accts[0], accts[1])
// 		if err != nil {
// 			a.fail(w, err.Error(), 500)
// 			return
// 		}
// 		for rows.Next() {
// 			acct := core.Account{}
// 			err := rows.Scan(acct.Name, acct.Pubkey, acct.Privkey, acct.Address)
// 			if err != nil {
// 				a.fail(w, err.Error(), 500)
// 				return
// 			}
// 		}

// 		channels = append(channels, &ch)
// 	}
// 	if rows.Err() != nil {
// 		a.fail(w, rows.Err().Error(), 500)
// 		return
// 	}

// 	data := struct {
// 		channels []*core.Channel
// 	}{channels}

// 	a.ok(w, data)
// }

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

// func addEscrowProvider(w http.ResponseWriter, r *http.Request) {
// 	if r.Body == nil {
// 		fmt.Println("no body")
// 		return
// 	}

// 	var ep core.EscrowProvider
// 	err := json.NewDecoder(r.Body).Decode(&ep)
// 	if err != nil {
// 		panic(err)
// 	}

// 	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
// 	if err != nil {
// 		panic(err)
// 	}

// 	db.NewRecord(ep)
// 	db.Create(&ep)

// 	json.NewEncoder(w).Encode(ep)
// }

func viewPeers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

// func addEscrowProvider(w http.ResponseWriter, r *http.Request) {
// 	if r.Body == nil {
// 		fmt.Println("no body")
// 		return
// 	}
// 	dec := json.NewDecoder(r.Body)

// 	var ep core.EscrowProvider
// 	if err := dec.Decode(&ep); err != nil {
// 		panic(err)
// 	}

// 	bytes, err := json.Marshal(ep)

// 	db, err := bolt.Open(".", 0600, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	db.Update(func(tx *bolt.Tx) error {
// 		indexes := tx.Bucket([]byte("Indexes"))
// 		escrowProviders := tx.Bucket([]byte("EscrowProviders"))
// 		err := escrowProviders.Put([]byte(ep.Name), bytes)
// 		err = indexes.Put(makeKey("EscrowProviders", "Pubkey", string(ep.Pubkey)), []byte(ep.Name))
// 		return err
// 	})
// }

func populateChannel(db *bolt.DB, ch *core.Channel, ChannelID string) error {
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

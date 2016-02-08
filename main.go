package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
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

func testTicker(t time.Time) {
	fmt.Println("Tick at", t)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func viewChannels(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func newChannel(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func viewAccounts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func (a *api) newAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func viewEscrowProviders(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func addEscrowProvider(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func viewPeers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

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

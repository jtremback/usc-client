package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jtremback/usc-client/access"
	core "github.com/jtremback/usc-core/client"
)

type api struct {
	db *bolt.DB
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

func (a *api) testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func (a *api) viewChannels(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func (a *api) newChannel(w http.ResponseWriter, r *http.Request) {

}

func (a *api) viewAccounts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func (a *api) newAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func (a *api) viewJudges(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
}

func (a *api) addJudge(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		a.fail(w, "no body", 500)
		return
	}

	jd := &core.Judge{}
	err := json.NewDecoder(r.Body).Decode(jd)
	if err != nil {
		panic(err)
	}

	a.db.Update(func(tx *bolt.Tx) error {
		access.SetJudge(tx, jd)
		return nil
	})

	a.send(w, "ok")
}

func (a *api) addMyAccount(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		a.fail(w, "no body", 500)
		return
	}

	req := &struct {
		Data        *core.MyAccount
		JudgePubkey []byte
	}{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		a.fail(w, "body parsing error", 500)
	}

	if len(req.JudgePubkey) == 0 {
		a.fail(w, "missing judge_pubkey", 500)
	}

	a.db.Update(func(tx *bolt.Tx) error {
		jd := &core.Judge{}
		err := json.Unmarshal(tx.Bucket([]byte("Judges")).Get([]byte(req.JudgePubkey)), jd)
		if err != nil {
			a.fail(w, "missing judge_pubkey", 500)
		}
		req.Data.Judge = jd
		access.SetMyAccount(tx, req.Data)
		return nil
	})

	a.send(w, "ok")
}

func (a *api) viewPeers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "one", "two")
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

func (a *api) send(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.fail(w, "oops something evil has happened", 500)
		return
	}
	w.Write(resp)
}

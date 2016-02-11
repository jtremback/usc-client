package caller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/jtremback/usc-client/access"
	core "github.com/jtremback/usc-core/client"
	"github.com/jtremback/usc-core/wire"
)

type api struct {
	db *bolt.DB
}

func send(ev *wire.Envelope, address string) error {
	b, err := proto.Marshal(ev)

	resp, err := http.Post(address, "application/octet-stream", bytes.NewReader(b))
	if err != nil {
		return errors.New("counterparty unresponsive")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("counterparty error")
	}

	return nil
}

func (a *api) newChannel(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		a.fail(w, "no body", 500)
		return
	}

	req := &struct {
		State              []byte
		MyAccountPubkey    []byte
		TheirAccountPubkey []byte
		HoldPeriod         uint32
	}{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		a.fail(w, "body parsing error", 500)
	}

	err = newChannel(a.db, req.State, req.MyAccountPubkey, req.TheirAccountPubkey, req.HoldPeriod)
	if err != nil {
		a.fail(w, err.Error(), 500)
	}
}

func newChannel(db *bolt.DB, state []byte, mpk []byte, tpk []byte, hold uint32) error {
	var err error
	ta := &core.TheirAccount{}
	ma := &core.MyAccount{}
	err = db.View(func(tx *bolt.Tx) error {
		ma, err = access.GetMyAccount(tx, mpk)
		if err != nil {
			return errors.New("database error")
		}
		if ma == nil {
			return errors.New("channel not found")
		}

		ta, err = access.GetTheirAccount(tx, tpk)
		if err != nil {
			return errors.New("database error")
		}
		if ta == nil {
			return errors.New("channel not found")
		}

		return nil
	})
	if err != nil {
		return err
	}

	otx, err := ma.NewOpeningTx(ta, state, hold)
	if err != nil {
		return errors.New("server error")
	}

	ev, err := ma.SignOpeningTx(otx)
	if err != nil {
		return errors.New("server error")
	}

	err = send(ev, ta.Address)
	if err != nil {
		return err
	}

	return nil
}

func (a *api) sendUpdateTx(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		a.fail(w, "no body", 500)
		return
	}

	req := &struct {
		State     []byte
		ChannelId []byte
		Fast      bool
	}{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		a.fail(w, "body parsing error", 500)
	}

	err = sendUpdateTx(a.db, req.State, req.ChannelId, req.Fast)
	if err != nil {
		a.fail(w, err.Error(), 500)
	}
}

func sendUpdateTx(db *bolt.DB, state []byte, chId []byte, fast bool) error {
	ch := &core.Channel{}
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		ch, err = access.GetChannel(tx, chId)
		if err != nil {
			return errors.New("database error")
		}
		if ch == nil {
			return errors.New("channel not found")
		}

		return nil
	})
	if err != nil {
		return err
	}

	utx, err := ch.NewUpdateTx(state, fast)
	if err != nil {
		return errors.New("server error")
	}

	ev, err := ch.SignUpdateTx(utx)
	if err != nil {
		return errors.New("server error")
	}

	err = send(ev, ch.TheirAccount.Address)
	if err != nil {
		return err
	}

	return nil
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
			a.fail(w, "database error", 500)
		}
		req.Data.Judge = jd
		access.SetMyAccount(tx, req.Data)
		return nil
	})

	a.send(w, "ok")
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

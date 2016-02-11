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

type Api struct {
	db *bolt.DB
}

func (a *Api) MountRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/new_channel", a.newChannel)
	mux.HandleFunc("/send_update_tx", a.sendUpdateTx)
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

func (a *Api) newChannel(w http.ResponseWriter, r *http.Request) {
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

func (a *Api) sendUpdateTx(w http.ResponseWriter, r *http.Request) {
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

func (a *Api) addJudge(w http.ResponseWriter, r *http.Request) {
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

func (a *Api) fail(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")

	data := struct {
		Error string
	}{Error: msg}

	resp, _ := json.Marshal(data)
	w.WriteHeader(status)
	w.Write(resp)
}

func (a *Api) send(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.fail(w, "oops something evil has happened", 500)
		return
	}
	w.Write(resp)
}

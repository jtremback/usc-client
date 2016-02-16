package counterparty

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/jtremback/usc-core/wire"
)

type Api struct {
	db            *bolt.DB
	callerAddress string
}

func (a *Api) MountRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/add_channel", a.addChannel)
}

func send(ev *wire.Envelope, address string) error {
	b, err := proto.Marshal(ev)

	resp, err := http.Post(address, "application/octet-stream", bytes.NewReader(b))
	if err != nil {
		return errors.New("network error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("counterparty error")
	}

	return nil
}

func (a *Api) addChannel(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		a.fail(w, "no body", 500)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.fail(w, "server error", 500)
	}

	ev := &wire.Envelope{}
	proto.Unmarshal(b, ev)

	err = addChannel(a.db, ev)
	if err != nil {
		a.fail(w, "server error", 500)
	}
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

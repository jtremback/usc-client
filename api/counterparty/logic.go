package counterparty

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/jtremback/upc_temp/wire"
	"github.com/jtremback/usc-client/access"
	core "github.com/jtremback/usc-core/client"
)

func confirmOpeningTx(db *bolt.DB, callerAddress string, ev *wire.Envelope) error {
	var err error
	ev, otx, err = ma.ConfirmOpeningTx(ev)
	if err != nil {
		return errors.New("server error")
	}

	ma := &core.MyAccount{}
	err = db.View(func(tx *bolt.Tx) error {
		ma, err = access.GetMyAccount(tx, mpk)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	b, err := json.Marshal(otx)

	resp, err := http.Post(address+"/confirm_opening_tx", "application/json", bytes.NewReader(b))
	if err != nil {
		return errors.New("network error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("caller error")
	}

	return nil
}

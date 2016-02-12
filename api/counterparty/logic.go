package counterparty

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/jtremback/usc-client/access"
	core "github.com/jtremback/usc-core/client"
	"github.com/jtremback/usc-core/wire"
)

func confirmOpeningTx(db *bolt.DB, callerAddress string, ev *wire.Envelope) error {
	var err error
	otx, err := core.UnmarshallOpeningTx(ev)
	if err != nil {
		return err
	}

	ma := &core.MyAccount{}
	err = db.View(func(tx *bolt.Tx) error {
		ma, err = access.GetMyAccount(tx, otx.Pubkeys[1])
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	ev, err = ma.ConfirmOpeningTx(ev, otx)
	if err != nil {
		return errors.New("server error")
	}

	b, err := json.Marshal(otx)

	resp, err := http.Post(callerAddress+"/confirm_opening_tx", "application/json", bytes.NewReader(b))
	if err != nil {
		return errors.New("network error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("caller error")
	}

	return nil
}

package protocol

import (
	"net/http"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/usc-core/client"
	"github.com/jtremback/usc-core/wire"
)

type api struct {
	db *bolt.DB
}

func (a api) ConfirmOpeningTx(w http.ResponseWriter, r *http.Request) {

}

func confirmOpeningTx(acct *core.MyAccount, ev *wire.Envelope) error {
	var err error
	ev, otx, err := acct.ConfirmOpeningTx(ev)

}

package protocol

import (
	"net/http"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/usc-core/client"
)

type api struct {
	db *bolt.DB
}

func (a api) ConfirmOpeningTx(w http.ResponseWriter, r *http.Request) {
	core.ConfirmOpeningTx()
}

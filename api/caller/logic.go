package caller

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/jtremback/usc-client/access"
	core "github.com/jtremback/usc-core/client"
)

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

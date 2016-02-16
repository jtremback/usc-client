package counterparty

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/jtremback/usc-client/access"
	core "github.com/jtremback/usc-core/client"
	"github.com/jtremback/usc-core/wire"
)

func AddChannel(db *bolt.DB, ev *wire.Envelope) error {
	var err error

	otx := &wire.OpeningTx{}
	err = proto.Unmarshal(ev.Payload, otx)
	if err != nil {
		return err
	}

	ma := &core.MyAccount{}
	ta := &core.TheirAccount{}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = access.GetChannel(tx, otx.ChannelId)
		if err != nil {
			return errors.New("channel already exists")
		}

		ta, err = access.GetTheirAccount(tx, otx.Pubkeys[0])
		if err != nil {
			return err
		}

		ma, err = access.GetMyAccount(tx, otx.Pubkeys[1])
		if err != nil {
			return err
		}

		ch, err := core.NewChannel(ev, ma, ta)

		access.SetChannel(tx, ch)
		if err != nil {
			return errors.New("database error")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func AddUpdateTx(db *bolt.DB, ev *wire.Envelope) error {
	var err error

	utx := &wire.UpdateTx{}
	err = proto.Unmarshal(ev.Payload, utx)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		ch, err := access.GetChannel(tx, utx.ChannelId)
		if err != nil {
			return err
		}

		err = ch.CheckUpdateTx(ev, utx)
		if err != nil {
			return err
		}

		ch.ProposedUpdateTx = utx
		ch.ProposedUpdateTxEnvelope = ev

		access.SetChannel(tx, ch)
		if err != nil {
			return errors.New("database error")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

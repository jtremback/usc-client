package logic

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/jtremback/usc-client/access"
	core "github.com/jtremback/usc-core/peer"
	"github.com/jtremback/usc-core/wire"
)

type Counterparty struct {
	db *bolt.DB
}

func (a *Counterparty) AddChannel(ev *wire.Envelope) error {
	var err error

	otx := &wire.OpeningTx{}
	err = proto.Unmarshal(ev.Payload, otx)
	if err != nil {
		return err
	}

	acct := &core.Account{}
	cpt := &core.Counterparty{}
	err = a.db.Update(func(tx *bolt.Tx) error {
		_, err = access.GetChannel(tx, otx.ChannelId)
		if err != nil {
			return errors.New("channel already exists")
		}

		cpt, err = access.GetCounterparty(tx, otx.Pubkeys[0])
		if err != nil {
			return err
		}

		acct, err = access.GetAccount(tx, otx.Pubkeys[1])
		if err != nil {
			return err
		}

		ch, err := core.NewChannel(ev, otx, acct, cpt)

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

func (a *Counterparty) AddUpdateTx(ev *wire.Envelope) error {
	var err error

	utx := &wire.UpdateTx{}
	err = proto.Unmarshal(ev.Payload, utx)
	if err != nil {
		return err
	}

	err = a.db.Update(func(tx *bolt.Tx) error {
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

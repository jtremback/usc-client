package logic

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/jtremback/usc-client/access"
	"github.com/jtremback/usc-client/clients"
	core "github.com/jtremback/usc-core/client"
	"github.com/jtremback/usc-core/wire"
)

type Caller struct {
	DB           *bolt.DB
	Counterparty *clients.Counterparty
	Judge        *clients.Judge
}

func (a *Caller) ProposeChannel(state []byte, mpk []byte, tpk []byte, hold uint32) error {
	var err error
	ta := &core.TheirAccount{}
	ma := &core.MyAccount{}
	err = a.DB.Update(func(tx *bolt.Tx) error {
		ma, err = access.GetMyAccount(tx, mpk)
		if err != nil {
			return err
		}

		ta, err = access.GetTheirAccount(tx, tpk)
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

		ch, err := core.NewChannel(ev, ma, ta)
		if err != nil {
			return errors.New("server error")
		}

		err = a.Counterparty.Send(ev, ta.Address)
		if err != nil {
			return err
		}

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

func (a *Caller) ConfirmChannel(chID string) error {
	var err error
	ch := &core.Channel{}
	err = a.DB.Update(func(tx *bolt.Tx) error {
		ch, err = access.GetChannel(tx, chID)
		if err != nil {
			return err
		}

		ch.OpeningTxEnvelope = ch.MyAccount.SignEnvelope(ch.OpeningTxEnvelope)

		access.SetChannel(tx, ch)
		if err != nil {
			return errors.New("database error")
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = a.Judge.Send(ch.OpeningTxEnvelope, ch.Judge.Address)
	if err != nil {
		return err
	}

	return nil
}

func (a *Caller) OpenChannel(ev *wire.Envelope) error {
	var err error

	ch := &core.Channel{}
	err = a.DB.Update(func(tx *bolt.Tx) error {
		otx := &wire.OpeningTx{}
		err = proto.Unmarshal(ev.Payload, otx)
		if err != nil {
			return err
		}

		ch, err = access.GetChannel(tx, otx.ChannelId)
		if err != nil {
			return err
		}

		ch.Open(ev, otx)
		if err != nil {
			return err
		}

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

func (a *Caller) SendUpdateTx(state []byte, chID string, fast bool) error {
	var err error
	ch := &core.Channel{}
	err = a.DB.Update(func(tx *bolt.Tx) error {
		ch, err = access.GetChannel(tx, chID)
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

		err = a.Counterparty.Send(ev, ch.TheirAccount.Address)
		if err != nil {
			return err
		}

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

func (a *Caller) ConfirmUpdateTx(chID string) error {
	var err error
	err = a.DB.Update(func(tx *bolt.Tx) error {
		ch, err := access.GetChannel(tx, chID)
		if err != nil {
			return err
		}

		_, err = ch.ConfirmUpdateTx()
		if err != nil {
			return err
		}

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

func (a *Caller) CheckFinalUpdateTx(ev *wire.Envelope) error {
	var err error
	utx := &wire.UpdateTx{}
	err = proto.Unmarshal(ev.Payload, utx)
	if err != nil {
		return err
	}

	err = a.DB.Update(func(tx *bolt.Tx) error {
		ch, err := access.GetChannel(tx, utx.ChannelId)
		if err != nil {
			return err
		}

		ev2, err := ch.CheckFinalUpdateTx(ev, utx)
		if err != nil {
			return err
		}
		if ev2 != nil {
			err = a.Judge.Send(ev2, ch.Judge.Address)
			if err != nil {
				return err
			}
		}

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

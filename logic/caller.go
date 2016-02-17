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
	DB             *bolt.DB
	CounterpartyCl *clients.Counterparty
	JudgeCl        *clients.Judge
}

func (a *Caller) ProposeChannel(state []byte, mpk []byte, tpk []byte, hold uint32) error {
	var err error
	cpt := &core.Counterparty{}
	acct := &core.Account{}
	err = a.DB.Update(func(tx *bolt.Tx) error {
		acct, err = access.GetMyAccount(tx, mpk)
		if err != nil {
			return err
		}

		cpt, err = access.GetTheirAccount(tx, tpk)
		if err != nil {
			return err
		}

		otx, err := acct.NewOpeningTx(cpt, state, hold)
		if err != nil {
			return errors.New("server error")
		}

		ev, err := acct.SignOpeningTx(otx)
		if err != nil {
			return errors.New("server error")
		}

		ch, err := core.NewChannel(ev, acct, cpt)
		if err != nil {
			return errors.New("server error")
		}

		err = a.CounterpartyCl.Send(ev, cpt.Address)
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

		ch.OpeningTxEnvelope = ch.Account.SignEnvelope(ch.OpeningTxEnvelope)

		access.SetChannel(tx, ch)
		if err != nil {
			return errors.New("database error")
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = a.JudgeCl.Send(ch.OpeningTxEnvelope, ch.Judge.Address)
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

		err = a.CounterpartyCl.Send(ev, ch.Counterparty.Address)
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
			err = a.JudgeCl.Send(ev2, ch.Judge.Address)
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

func (a *Caller) AddJudge() {}

func (a *Caller) NewAccount() {}

func (a *Caller) AddCounterparty() {}

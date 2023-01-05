package logic

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

type SimpleContext struct {
	inner       context.Context
	cancelFunc  context.CancelFunc
	cancelled   bool
	transaction sq.Tx
}

func CreateSimpleContext(innerCtx context.Context, cancelFn context.CancelFunc) *SimpleContext {
	return &SimpleContext{
		inner:       innerCtx,
		cancelFunc:  cancelFn,
		cancelled:   false,
		transaction: nil,
	}
}

func (sc *SimpleContext) Deadline() (deadline time.Time, ok bool) {
	return sc.inner.Deadline()
}

func (sc *SimpleContext) Done() <-chan struct{} {
	return sc.inner.Done()
}

func (sc *SimpleContext) Err() error {
	return sc.inner.Err()
}

func (sc *SimpleContext) Value(key any) any {
	return sc.inner.Value(key)
}

func (sc *SimpleContext) Cancel() {
	sc.cancelled = true
	if sc.transaction != nil {
		log.Error().Msg("Rollback transaction")
		err := sc.transaction.Rollback()
		if err != nil {
			panic("failed to rollback transaction: " + err.Error())
		}
		sc.transaction = nil
	}
	sc.cancelFunc()
}

func (sc *SimpleContext) GetOrCreateTransaction(db db.DatabaseImpl) (sq.Tx, error) {
	if sc.cancelled {
		return nil, errors.New("context cancelled")
	}
	if sc.transaction != nil {
		return sc.transaction, nil
	}
	tx, err := db.BeginTx(sc)
	if err != nil {
		return nil, err
	}
	sc.transaction = tx
	return tx, nil
}

func (sc *SimpleContext) CommitTransaction() error {
	if sc.transaction == nil {
		return nil
	}
	err := sc.transaction.Commit()
	if err != nil {
		return err
	}
	sc.transaction = nil
	return nil
}

func (sc *SimpleContext) RollbackTransaction() {
	if sc.transaction == nil {
		return
	}
	err := sc.transaction.Rollback()
	if err != nil {
		log.Err(err).Stack().Msg("Failed to rollback transaction")
		panic(err)
	}
	sc.transaction = nil
	return
}

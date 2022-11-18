package logic

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/db"
	"context"
	"database/sql"
	"errors"
	"time"
)

type AppContext struct {
	inner       context.Context
	cancelFunc  context.CancelFunc
	cancelled   bool
	transaction *sql.Tx
	permissions PermissionSet
}

func CreateAppContext(innerCtx context.Context, cancelFn context.CancelFunc) *AppContext {
	return &AppContext{
		inner:       innerCtx,
		cancelFunc:  cancelFn,
		cancelled:   false,
		transaction: nil,
		permissions: NewEmptyPermissions(),
	}
}

func (ac *AppContext) Deadline() (deadline time.Time, ok bool) {
	return ac.inner.Deadline()
}

func (ac *AppContext) Done() <-chan struct{} {
	return ac.inner.Done()
}

func (ac *AppContext) Err() error {
	return ac.inner.Err()
}

func (ac *AppContext) Value(key any) any {
	return ac.inner.Value(key)
}

func (ac *AppContext) Cancel() {
	ac.cancelled = true
	if ac.transaction != nil {
		err := ac.transaction.Rollback()
		if err != nil {
			panic("failed to rollback transaction: " + err.Error())
		}
		ac.transaction = nil
	}
	ac.cancelFunc()
}

func (ac *AppContext) FinishSuccess(res ginresp.HTTPResponse) ginresp.HTTPResponse {
	if ac.cancelled {
		panic("Cannot finish a cancelled request")
	}
	if ac.transaction != nil {
		err := ac.transaction.Commit()
		if err != nil {
			return ginresp.InternAPIError(500, apierr.COMMIT_FAILED, "Failed to comit changes to DB", err)
		}
		ac.transaction = nil
	}
	return res
}

func (ac *AppContext) GetOrCreateTransaction(db *db.Database) (*sql.Tx, error) {
	if ac.cancelled {
		return nil, errors.New("context cancelled")
	}
	if ac.transaction != nil {
		return ac.transaction, nil
	}
	tx, err := db.BeginTx(ac)
	if err != nil {
		return nil, err
	}
	ac.transaction = tx
	return tx, nil
}

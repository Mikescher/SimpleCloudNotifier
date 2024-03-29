package logic

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

type TxContext interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any

	GetOrCreateTransaction(db db.DatabaseImpl) (sq.Tx, error)
}

type AppContext struct {
	app         *Application
	inner       context.Context
	cancelFunc  context.CancelFunc
	cancelled   bool
	transaction sq.Tx
	permissions models.PermissionSet
	ginContext  *gin.Context
}

func CreateAppContext(app *Application, g *gin.Context, innerCtx context.Context, cancelFn context.CancelFunc) *AppContext {
	return &AppContext{
		app:         app,
		inner:       innerCtx,
		cancelFunc:  cancelFn,
		cancelled:   false,
		transaction: nil,
		permissions: models.NewEmptyPermissions(),
		ginContext:  g,
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
		log.Error().Str("uri", ac.RequestURI()).Msg("Rollback transaction (ctx-cancel)")
		err := ac.transaction.Rollback()
		if err != nil {
			log.Err(err).Stack().Msg("Failed to rollback transaction")
		}
		ac.transaction = nil
	}
	ac.cancelFunc()
}

func (ac *AppContext) RequestURI() string {
	if ac.ginContext != nil && ac.ginContext.Request != nil {
		return ac.ginContext.Request.Method + " :: " + ac.ginContext.Request.RequestURI
	} else {
		return ""
	}
}

func (ac *AppContext) FinishSuccess(res ginresp.HTTPResponse) ginresp.HTTPResponse {
	if ac.cancelled {
		panic("Cannot finish a cancelled request")
	}
	if ac.transaction != nil {
		err := ac.transaction.Commit()
		if err != nil {
			return ginresp.APIError(ac.ginContext, 500, apierr.COMMIT_FAILED, "Failed to comit changes to DB", err)
		}
		ac.transaction = nil
	}
	return res
}

func (ac *AppContext) GetOrCreateTransaction(db db.DatabaseImpl) (sq.Tx, error) {
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

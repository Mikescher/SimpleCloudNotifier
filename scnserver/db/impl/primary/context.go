package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db"
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

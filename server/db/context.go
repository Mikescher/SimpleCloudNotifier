package db

import (
	"blackforestbytes.com/simplecloudnotifier/sq"
	"time"
)

type TxContext interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any

	GetOrCreateTransaction(db *Database) (sq.Tx, error)
}

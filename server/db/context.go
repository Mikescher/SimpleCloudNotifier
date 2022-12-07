package db

import (
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

type TxContext interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any

	GetOrCreateTransaction(db *Database) (sq.Tx, error)
}

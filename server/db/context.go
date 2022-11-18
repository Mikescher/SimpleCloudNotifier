package db

import (
	"database/sql"
	"time"
)

type TxContext interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any

	GetOrCreateTransaction(db *Database) (*sql.Tx, error)
}

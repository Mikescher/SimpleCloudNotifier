package sq

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"sync"
)

type DB interface {
	Exec(ctx context.Context, sql string, prep PP) (sql.Result, error)
	Query(ctx context.Context, sql string, prep PP) (*sqlx.Rows, error)
	Ping(ctx context.Context) error
	BeginTransaction(ctx context.Context, iso sql.IsolationLevel) (Tx, error)
}

type database struct {
	db    *sqlx.DB
	txctr uint16
	lock  sync.Mutex
}

func NewDB(db *sqlx.DB) DB {
	return &database{
		db:    db,
		txctr: 0,
		lock:  sync.Mutex{},
	}
}

func (db *database) Exec(ctx context.Context, sql string, prep PP) (sql.Result, error) {
	log.Debug().Msg(fmt.Sprintf("[SQL-EXEC] %s", fmtSQLPrint(sql)))

	res, err := db.db.NamedExecContext(ctx, sql, prep)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db *database) Query(ctx context.Context, sql string, prep PP) (*sqlx.Rows, error) {
	log.Debug().Msg(fmt.Sprintf("[SQL-QUERY] %s", fmtSQLPrint(sql)))

	rows, err := db.db.NamedQueryContext(ctx, sql, prep)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *database) Ping(ctx context.Context) error {
	log.Debug().Msg("[SQL-PING]")

	err := db.db.PingContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (db *database) BeginTransaction(ctx context.Context, iso sql.IsolationLevel) (Tx, error) {
	db.lock.Lock()
	txid := db.txctr
	db.txctr += 1 // with overflow !
	db.lock.Unlock()

	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%d>-START]", txid))

	xtx, err := db.db.BeginTxx(ctx, &sql.TxOptions{Isolation: iso})
	if err != nil {
		return nil, err
	}

	return NewTransaction(xtx, txid), nil
}

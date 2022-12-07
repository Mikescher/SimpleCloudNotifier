package sq

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Tx interface {
	Rollback() error
	Commit() error
	Exec(ctx context.Context, sql string, prep PP) (sql.Result, error)
	Query(ctx context.Context, sql string, prep PP) (*sqlx.Rows, error)
}

type transaction struct {
	tx *sqlx.Tx
	id uint16
}

func NewTransaction(xtx *sqlx.Tx, txid uint16) Tx {
	return &transaction{
		tx: xtx,
		id: txid,
	}
}

func (tx *transaction) Rollback() error {
	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%d>-ROLLBACK]", tx.id))

	return tx.tx.Rollback()
}

func (tx *transaction) Commit() error {
	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%d>-COMMIT]", tx.id))

	return tx.tx.Commit()
}

func (tx *transaction) Exec(ctx context.Context, sql string, prep PP) (sql.Result, error) {
	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%d>-EXEC] %s", tx.id, fmtSQLPrint(sql)))

	res, err := tx.tx.NamedExecContext(ctx, sql, prep)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (tx *transaction) Query(ctx context.Context, sql string, prep PP) (*sqlx.Rows, error) {
	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%d>-QUERY] %s", tx.id, fmtSQLPrint(sql)))

	rows, err := sqlx.NamedQueryContext(ctx, tx.tx, sql, prep)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

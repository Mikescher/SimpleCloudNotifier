package sq

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Queryable interface {
	Exec(ctx context.Context, sql string, prep PP) (sql.Result, error)
	Query(ctx context.Context, sql string, prep PP) (*sqlx.Rows, error)
}

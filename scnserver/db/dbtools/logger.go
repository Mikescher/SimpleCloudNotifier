package dbtools

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"strings"
)

type DBLogger struct {
	Ident string
}

func (l DBLogger) PrePing(ctx context.Context) error {
	log.Debug().Msg("[SQL-PING]")

	return nil
}

func (l DBLogger) PreTxBegin(ctx context.Context, txid uint16) error {
	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%s|%d>-START]", l.Ident, txid))

	return nil
}

func (l DBLogger) PreTxCommit(txid uint16) error {
	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%s|%d>-COMMIT]", l.Ident, txid))

	return nil
}

func (l DBLogger) PreTxRollback(txid uint16) error {
	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%s|%d>-ROLLBACK]", l.Ident, txid))

	return nil
}

func (l DBLogger) PreQuery(ctx context.Context, txID *uint16, sql *string, params *sq.PP) error {
	if txID == nil {
		log.Debug().Msg(fmt.Sprintf("[SQL<%s>-QUERY] %s", l.Ident, fmtSQLPrint(*sql)))
	} else {
		log.Debug().Msg(fmt.Sprintf("[SQL-TX<%s|%d>-QUERY] %s", l.Ident, *txID, fmtSQLPrint(*sql)))
	}

	return nil
}

func (l DBLogger) PreExec(ctx context.Context, txID *uint16, sql *string, params *sq.PP) error {
	if txID == nil {
		log.Debug().Msg(fmt.Sprintf("[SQL-<%s>-EXEC] %s", l.Ident, fmtSQLPrint(*sql)))
	} else {
		log.Debug().Msg(fmt.Sprintf("[SQL-TX<%s|%d>-EXEC] %s", l.Ident, *txID, fmtSQLPrint(*sql)))
	}

	return nil
}

func (l DBLogger) PostPing(result error) {
	//
}

func (l DBLogger) PostTxBegin(txid uint16, result error) {
	//
}

func (l DBLogger) PostTxCommit(txid uint16, result error) {
	//
}

func (l DBLogger) PostTxRollback(txid uint16, result error) {
	//
}

func (l DBLogger) PostQuery(txID *uint16, sqlOriginal string, sqlReal string, params sq.PP) {
	//
}

func (l DBLogger) PostExec(txID *uint16, sqlOriginal string, sqlReal string, params sq.PP) {
	//
}

func fmtSQLPrint(sql string) string {
	if strings.Contains(sql, ";") && len(sql) > 1024 {
		return "(...multi...)"
	}

	sql = strings.ReplaceAll(sql, "\r", "")
	sql = strings.ReplaceAll(sql, "\n", " ")

	return sql
}

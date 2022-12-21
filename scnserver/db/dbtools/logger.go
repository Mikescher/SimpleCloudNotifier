package dbtools

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"strings"
)

type DBLogger struct{}

func (l DBLogger) PrePing(ctx context.Context) error {
	log.Debug().Msg("[SQL-PING]")

	return nil
}

func (l DBLogger) PreTxBegin(ctx context.Context, txid uint16) error {
	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%d>-START]", txid))

	return nil
}

func (l DBLogger) PreTxCommit(txid uint16) error {
	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%d>-COMMIT]", txid))

	return nil
}

func (l DBLogger) PreTxRollback(txid uint16) error {
	log.Debug().Msg(fmt.Sprintf("[SQL-TX<%d>-ROLLBACK]", txid))

	return nil
}

func (l DBLogger) PreQuery(ctx context.Context, txID *uint16, sql *string, params *sq.PP) error {
	if txID == nil {
		log.Debug().Msg(fmt.Sprintf("[SQL-QUERY] %s", fmtSQLPrint(*sql)))
	} else {
		log.Debug().Msg(fmt.Sprintf("[SQL-TX<%d>-QUERY] %s", *txID, fmtSQLPrint(*sql)))
	}

	return nil
}

func (l DBLogger) PreExec(ctx context.Context, txID *uint16, sql *string, params *sq.PP) error {
	if txID == nil {
		log.Debug().Msg(fmt.Sprintf("[SQL-EXEC] %s", fmtSQLPrint(*sql)))
	} else {
		log.Debug().Msg(fmt.Sprintf("[SQL-TX<%d>-EXEC] %s", *txID, fmtSQLPrint(*sql)))
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
	if strings.Contains(sql, ";") {
		return "(...multi...)"
	}

	sql = strings.ReplaceAll(sql, "\r", "")
	sql = strings.ReplaceAll(sql, "\n", " ")

	return sql
}

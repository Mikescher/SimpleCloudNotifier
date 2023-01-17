package dbtools

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"regexp"
	"strings"
)

var rexWhitespaceRun = rext.W(regexp.MustCompile("\\s{2,}"))

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
	if strings.Contains(strings.TrimRight(sql, ";\r\n\t "), ";") {

		str := "(...multi...)"
		for _, v := range strings.Split(sql, ";") {

			v = strings.ReplaceAll(v, "\r", "")
			v = strings.ReplaceAll(v, "\n", " ")
			v = strings.TrimRight(v, ";")
			v = strings.TrimSpace(v)
			v = rexWhitespaceRun.ReplaceAll(v, " ", true)

			str += "\n" + "    " + v
		}
		return str

	} else {

		sql = strings.ReplaceAll(sql, "\r", "")
		sql = strings.ReplaceAll(sql, "\n", " ")
		sql = rexWhitespaceRun.ReplaceAll(sql, " ", true)

		return sql

	}

}

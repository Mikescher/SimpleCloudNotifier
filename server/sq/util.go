package sq

import (
	"strings"
)

func fmtSQLPrint(sql string) string {
	if strings.Contains(sql, ";") {
		return "(...multi...)"
	}

	sql = strings.ReplaceAll(sql, "\r", "")
	sql = strings.ReplaceAll(sql, "\n", " ")

	return sql
}

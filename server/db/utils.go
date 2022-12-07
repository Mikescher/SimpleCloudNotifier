package db

import (
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"strings"
	"time"
)

func bool2DB(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func time2DB(t time.Time) int64 {
	return t.UnixMilli()
}

func time2DBOpt(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	return langext.Ptr(t.UnixMilli())
}

func fmtSQLPrint(sql string) string {
	if strings.Contains(sql, ";") {
		return "(...multi...)"
	}

	sql = strings.ReplaceAll(sql, "\r", "")
	sql = strings.ReplaceAll(sql, "\n", " ")

	return sql
}

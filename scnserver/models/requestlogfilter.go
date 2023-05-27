package models

import (
	"crypto/sha512"
	"encoding/hex"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"strings"
)

type RequestLogFilter struct {
}

func (f RequestLogFilter) SQL() (string, string, sq.PP, error) {

	joinClause := ""

	// ...

	sqlClauses := make([]string, 0)

	params := sq.PP{}

	// ...

	sqlClause := ""
	if len(sqlClauses) > 0 {
		sqlClause = strings.Join(sqlClauses, " AND ")
	} else {
		sqlClause = "1=1"
	}

	return sqlClause, joinClause, params, nil
}

func (f RequestLogFilter) Hash() string {
	bh, err := dataext.StructHash(f, dataext.StructHashOptions{HashAlgo: sha512.New()})
	if err != nil {
		return "00000000"
	}

	str := hex.EncodeToString(bh)
	return str[0:mathext.Min(8, len(bh))]
}

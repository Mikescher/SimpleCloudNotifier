package dbtools

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"regexp"
	"strings"
	"sync"
)

//
// This is..., not good...
//
// for sq.ScanAll to work with (left-)joined tables _need_ to get column names aka "alias.column"
// But sqlite (and all other db server) only return "column" if we don't manually specify `alias.column as "alias.columnname"`
// But always specifying all columns (and their alias) would be __very__ cumbersome...
//
// The "solution" is this preprocessor, which translates queries of the form `SELECT tab1.*, tab2.* From tab1` into `SELECT tab1.col1 AS "tab1.col1", tab1.col2 AS "tab1.col2" ....`
//
// Prerequisites:
//   - all aliased tables must be written as `tablename AS alias` (the variant without the AS keyword is invalid)
//   - a star only expands to the (single) table in FROM. Use *, table2.* if there exists a second (joined) table
//   - No weird SQL syntax, this "parser" is not very robust...
//

type DBPreprocessor struct {
	db sq.DB

	lock         sync.Mutex
	cacheColumns map[string][]string
	cacheQuery   map[string]string
}

var regexAlias = regexp.MustCompile("([A-Za-z_\\-0-9]+)\\s+AS\\s+([A-Za-z_\\-0-9]+)")

func NewDBPreprocessor(db sq.DB) *DBPreprocessor {
	return &DBPreprocessor{
		db:           db,
		lock:         sync.Mutex{},
		cacheColumns: make(map[string][]string),
		cacheQuery:   make(map[string]string),
	}
}

func (pp *DBPreprocessor) PrePing(ctx context.Context) error {
	return nil
}

func (pp *DBPreprocessor) PreTxBegin(ctx context.Context, txid uint16) error {
	return nil
}

func (pp *DBPreprocessor) PreTxCommit(txid uint16) error {
	return nil
}

func (pp *DBPreprocessor) PreTxRollback(txid uint16) error {
	return nil
}

func (pp *DBPreprocessor) PreQuery(ctx context.Context, txID *uint16, sql *string, params *sq.PP) error {
	sqlOriginal := *sql

	pp.lock.Lock()
	v, ok := pp.cacheQuery[sqlOriginal]
	pp.lock.Unlock()

	if ok {
		*sql = v
		return nil
	}

	if !strings.HasPrefix(sqlOriginal, "SELECT ") {
		return nil
	}

	idxFrom := strings.Index(sqlOriginal, " FROM ")
	if idxFrom < 0 {
		return nil
	}

	fromTableName := strings.Split(strings.TrimSpace(sqlOriginal[idxFrom+len(" FROM"):]), " ")[0]

	sels := strings.TrimSpace(sqlOriginal[len("SELECT "):idxFrom])

	split := strings.Split(sels, ",")

	newsel := make([]string, 0)

	aliasMap := make(map[string]string)
	for _, v := range regexAlias.FindAllStringSubmatch(sqlOriginal, idxFrom+len(" FROM")) {
		aliasMap[strings.TrimSpace(v[2])] = strings.TrimSpace(v[1])
	}

	for _, expr := range split {

		expr = strings.TrimSpace(expr)

		if expr == "*" {

			columns, err := pp.getTableColumns(ctx, fromTableName)
			if err != nil {
				return err
			}

			for _, colname := range columns {
				newsel = append(newsel, fmt.Sprintf("%s.%s AS \"%s\"", fromTableName, colname, colname))
			}

		} else if strings.HasSuffix(expr, ".*") {

			tableName := expr[0 : len(expr)-2]

			if tableRealName, ok := aliasMap[tableName]; ok {

				columns, err := pp.getTableColumns(ctx, tableRealName)
				if err != nil {
					return err
				}

				for _, colname := range columns {
					newsel = append(newsel, fmt.Sprintf("%s.%s AS \"%s.%s\"", tableName, colname, tableName, colname))
				}

			} else if tableName == fromTableName {

				columns, err := pp.getTableColumns(ctx, tableName)
				if err != nil {
					return err
				}

				for _, colname := range columns {
					newsel = append(newsel, fmt.Sprintf("%s.%s AS \"%s\"", tableName, colname, colname))
				}

			} else {

				columns, err := pp.getTableColumns(ctx, tableName)
				if err != nil {
					return err
				}

				for _, colname := range columns {
					newsel = append(newsel, fmt.Sprintf("%s.%s AS \"%s.%s\"", tableName, colname, tableName, colname))
				}

			}

		} else {
			return nil
		}

	}

	newSQL := "SELECT " + strings.Join(newsel, ", ") + sqlOriginal[idxFrom:]

	pp.lock.Lock()
	pp.cacheQuery[sqlOriginal] = newSQL
	pp.lock.Unlock()

	log.Debug().Msgf("Preprocessed SQL statement from '%s' --to--> '%s'", sqlOriginal, newSQL)

	*sql = newSQL
	return nil
}

func (pp *DBPreprocessor) PreExec(ctx context.Context, txID *uint16, sql *string, params *sq.PP) error {
	return nil
}

func (pp *DBPreprocessor) PostPing(result error) {
	//
}

func (pp *DBPreprocessor) PostTxBegin(txid uint16, result error) {
	//
}

func (pp *DBPreprocessor) PostTxCommit(txid uint16, result error) {
	//
}

func (pp *DBPreprocessor) PostTxRollback(txid uint16, result error) {
	//
}

func (pp *DBPreprocessor) PostQuery(txID *uint16, sqlOriginal string, sqlReal string, params sq.PP) {
	//
}

func (pp *DBPreprocessor) PostExec(txID *uint16, sqlOriginal string, sqlReal string, params sq.PP) {
	//
}

func (pp *DBPreprocessor) getTableColumns(ctx context.Context, tablename string) ([]string, error) {
	pp.lock.Lock()
	v, ok := pp.cacheColumns[tablename]
	pp.lock.Unlock()

	if ok {
		return v, nil
	}

	type res struct {
		CID     int64   `db:"cid"`
		Name    string  `db:"name"`
		Type    string  `db:"type"`
		NotNull int     `db:"notnull"`
		DFLT    *string `db:"dflt_value"`
		PK      int     `db:"pk"`
	}

	rows, err := pp.db.Query(ctx, "PRAGMA table_info('"+tablename+"');", sq.PP{})
	if err != nil {
		return nil, err
	}

	resrows, err := sq.ScanAll[res](rows, true)
	if err != nil {
		return nil, err
	}

	columns := langext.ArrMap(resrows, func(v res) string { return v.Name })

	if len(columns) == 0 {
		return nil, errors.New("no columns in table '" + tablename + "' (table does not exist?)")
	}

	pp.lock.Lock()
	pp.cacheColumns[tablename] = columns
	pp.lock.Unlock()

	return columns, nil
}

package dbtools

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"regexp"
	"strings"
	"sync"
	"time"
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

	lock       sync.Mutex
	dbTables   []string
	dbColumns  map[string][]string
	cacheQuery map[string]string
}

var regexAlias = rext.W(regexp.MustCompile("([A-Za-z_\\-0-9]+)\\s+AS\\s+([A-Za-z_\\-0-9]+)"))

func NewDBPreprocessor(db sq.DB) (*DBPreprocessor, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	obj := &DBPreprocessor{
		db:         db,
		lock:       sync.Mutex{},
		cacheQuery: make(map[string]string),
	}

	err := obj.Init(ctx)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (pp *DBPreprocessor) Init(ctx context.Context) error {

	dbTables := make([]string, 0)
	dbColumns := make(map[string][]string, 0)

	type tabInfo struct {
		Name string `db:"name"`
	}
	type colInfo struct {
		Name string `db:"name"`
	}

	rows1, err := pp.db.Query(ctx, "PRAGMA table_list;", sq.PP{})
	if err != nil {
		return err
	}
	resrows1, err := sq.ScanAll[tabInfo](rows1, sq.SModeFast, sq.Unsafe, true)
	if err != nil {
		return err
	}
	for _, tab := range resrows1 {

		rows2, err := pp.db.Query(ctx, fmt.Sprintf("PRAGMA table_info(\"%s\");", tab.Name), sq.PP{})
		if err != nil {
			return err
		}
		resrows2, err := sq.ScanAll[colInfo](rows2, sq.SModeFast, sq.Unsafe, true)
		if err != nil {
			return err
		}
		columns := langext.ArrMap(resrows2, func(v colInfo) string { return v.Name })

		dbTables = append(dbTables, tab.Name)
		dbColumns[tab.Name] = columns
	}

	pp.dbTables = dbTables
	pp.dbColumns = dbColumns

	return nil
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
	for _, v := range regexAlias.MatchAll(sqlOriginal) {
		aliasMap[strings.TrimSpace(v.GroupByIndex(2).Value())] = strings.TrimSpace(v.GroupByIndex(1).Value())
	}

	for _, expr := range split {

		expr = strings.TrimSpace(expr)

		if expr == "*" {

			columns, ok := pp.dbColumns[fromTableName]
			if !ok {
				return errors.New(fmt.Sprintf("[preprocessor]: table '%s' not found", fromTableName))
			}

			for _, colname := range columns {
				newsel = append(newsel, fmt.Sprintf("%s.%s AS \"%s\"", fromTableName, colname, colname))
			}

		} else if strings.HasSuffix(expr, ".*") {

			tableName := expr[0 : len(expr)-2]

			if tableRealName, ok := aliasMap[tableName]; ok {

				columns, ok := pp.dbColumns[tableRealName]
				if !ok {
					return errors.New(fmt.Sprintf("[sql-preprocessor]: table '%s' not found", tableRealName))
				}

				for _, colname := range columns {
					newsel = append(newsel, fmt.Sprintf("%s.%s AS \"%s.%s\"", tableName, colname, tableName, colname))
				}

			} else if tableName == fromTableName {

				columns, ok := pp.dbColumns[tableName]
				if !ok {
					return errors.New(fmt.Sprintf("[sql-preprocessor]: table '%s' not found", tableName))
				}

				for _, colname := range columns {
					newsel = append(newsel, fmt.Sprintf("%s.%s AS \"%s\"", tableName, colname, colname))
				}

			} else {

				columns, ok := pp.dbColumns[tableName]
				if !ok {
					return errors.New(fmt.Sprintf("[sql-preprocessor]: table '%s' not found", tableName))
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

	log.Debug().Msgf("Preprocessed SQL statement from\n'%s'\n--to-->\n'%s'", sqlOriginal, newSQL)

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

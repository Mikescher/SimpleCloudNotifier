package models

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

func timeOptFmt(t *time.Time, fmt string) *string {
	if t == nil {
		return nil
	} else {
		return langext.Ptr(t.Format(fmt))
	}
}

func timeOptFromMilli(millis *int64) *time.Time {
	if millis == nil {
		return nil
	}
	return langext.Ptr(time.UnixMilli(*millis))
}

func scanSingle[TData any](rows *sqlx.Rows) (TData, error) {
	if rows.Next() {
		var data TData
		err := rows.StructScan(&data)
		if err != nil {
			return *new(TData), err
		}
		if rows.Next() {
			return *new(TData), errors.New("sql returned more than onw row")
		}
		return data, nil
	} else {
		return *new(TData), sql.ErrNoRows
	}
}

func scanAll[TData any](rows *sqlx.Rows) ([]TData, error) {
	res := make([]TData, 0)
	for rows.Next() {
		var data TData
		err := rows.StructScan(&data)
		if err != nil {
			return nil, err
		}
		res = append(res, data)
	}
	return res, nil
}

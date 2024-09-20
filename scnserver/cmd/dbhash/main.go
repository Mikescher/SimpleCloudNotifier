package main

import (
	"blackforestbytes.com/simplecloudnotifier/db/schema"
	"context"
	"database/sql"
	"fmt"
	"github.com/glebarez/go-sqlite"
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func main() {
	exerr.Init(exerr.ErrorPackageConfigInit{})

	ctx, cancel := context.WithTimeout(context.Background(), 1011*time.Second)
	defer cancel()

	if !langext.InArray("sqlite3", sql.Drivers()) {
		sqlite.RegisterAsSQLITE3()
	}

	for key, schemaObj := range langext.AsSortedBy(langext.MapToArr(schema.PrimarySchema), func(v langext.MapEntry[int, schema.Def]) int { return v.Key }) {
		var h0 string
		if key == 1 {
			h0 = "N/A"
		} else {
			var err error
			h0, err = sq.HashGoSqliteSchema(ctx, schemaObj.Value.SQL)
			if err != nil {
				h0 = "ERR"
			}
		}
		fmt.Printf("PrimarySchema    [%d] := %s%s\n", schemaObj.Key, h0, langext.Conditional(schemaObj.Key == schema.PrimarySchemaVersion, "     (active)", ""))
	}

	fmt.Printf("\n")

	for _, schemaObj := range langext.AsSortedBy(langext.MapToArr(schema.RequestsSchema), func(v langext.MapEntry[int, schema.Def]) int { return v.Key }) {
		h0, err := sq.HashGoSqliteSchema(ctx, schemaObj.Value.SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("RequestsSchema   [%d] := %s%s\n", schemaObj.Key, h0, langext.Conditional(schemaObj.Key == schema.RequestsSchemaVersion, "     (active)", ""))
	}

	fmt.Printf("\n")

	for _, schemaObj := range langext.AsSortedBy(langext.MapToArr(schema.LogsSchema), func(v langext.MapEntry[int, schema.Def]) int { return v.Key }) {
		h0, err := sq.HashGoSqliteSchema(ctx, schemaObj.Value.SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("LogsSchema       [%d] := %s%s\n", schemaObj.Key, h0, langext.Conditional(schemaObj.Key == schema.LogsSchemaVersion, "     (active)", ""))
	}

	fmt.Printf("\n")
}

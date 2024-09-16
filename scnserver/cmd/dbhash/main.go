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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if !langext.InArray("sqlite3", sql.Drivers()) {
		sqlite.RegisterAsSQLITE3()
	}

	fmt.Println()

	for i := 2; i <= schema.PrimarySchemaVersion; i++ {
		h0, err := sq.HashGoSqliteSchema(ctx, schema.PrimarySchema[i].SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("PrimarySchema%d   := %s\n", i, h0)
	}

	for i := 1; i <= schema.RequestsSchemaVersion; i++ {
		h0, err := sq.HashGoSqliteSchema(ctx, schema.RequestsSchema[i].SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("RequestsSchema%d  := %s\n", i, h0)
	}

	for i := 1; i <= schema.LogsSchemaVersion; i++ {
		h0, err := sq.HashGoSqliteSchema(ctx, schema.LogsSchema[i].SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("LogsSchema%d      := %s\n", i, h0)
	}

	fmt.Println()

}

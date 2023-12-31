package main

import (
	"blackforestbytes.com/simplecloudnotifier/db/schema"
	"context"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sqlite3.Version() // ensure slite3 loaded
	{
		h0, err := sq.HashSqliteSchema(ctx, schema.PrimarySchema[1].SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("PrimarySchema1  := %s\n", h0)
	}
	{
		h0, err := sq.HashSqliteSchema(ctx, schema.PrimarySchema[2].SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("PrimarySchema2  := %s\n", h0)
	}
	{
		h0, err := sq.HashSqliteSchema(ctx, schema.PrimarySchema[3].SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("PrimarySchema3  := %s\n", h0)
	}
	{
		h0, err := sq.HashSqliteSchema(ctx, schema.PrimarySchema[4].SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("PrimarySchema4  := %s\n", h0)
	}
	{
		h0, err := sq.HashSqliteSchema(ctx, schema.RequestsSchema[1].SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("RequestsSchema1 := %s\n", h0)
	}
	{
		h0, err := sq.HashSqliteSchema(ctx, schema.LogsSchema[1].SQL)
		if err != nil {
			h0 = "ERR"
		}
		fmt.Printf("LogsSchema1     := %s\n", h0)
	}
}

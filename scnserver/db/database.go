package db

import (
	"context"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
)

type DatabaseImpl interface {
	DB() sq.DB

	Migrate(ctx context.Context) error
	Ping(ctx context.Context) error
	BeginTx(ctx context.Context) (sq.Tx, error)
	Stop(ctx context.Context) error

	ReadSchema(ctx context.Context) (int, error)

	WriteMetaString(ctx context.Context, key string, value string) error
	WriteMetaInt(ctx context.Context, key string, value int64) error
	WriteMetaReal(ctx context.Context, key string, value float64) error
	WriteMetaBlob(ctx context.Context, key string, value []byte) error

	ReadMetaString(ctx context.Context, key string) (*string, error)
	ReadMetaInt(ctx context.Context, key string) (*int64, error)
	ReadMetaReal(ctx context.Context, key string) (*float64, error)
	ReadMetaBlob(ctx context.Context, key string) (*[]byte, error)

	DeleteMeta(ctx context.Context, key string) error
}

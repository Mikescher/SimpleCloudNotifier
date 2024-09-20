package db

import (
	"context"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
)

type DatabaseImpl interface {
	DB() sq.DB

	Migrate(ctx context.Context) error
	Ping(ctx context.Context) error
	Version(ctx context.Context) (string, string, error)
	BeginTx(ctx context.Context) (sq.Tx, error)
	Stop(ctx context.Context) error

	ReadSchema(ctx TxContext) (int, error)

	WriteMetaString(ctx TxContext, key string, value string) error
	WriteMetaInt(ctx TxContext, key string, value int64) error
	WriteMetaReal(ctx TxContext, key string, value float64) error
	WriteMetaBlob(ctx TxContext, key string, value []byte) error

	ReadMetaString(ctx TxContext, key string) (*string, error)
	ReadMetaInt(ctx TxContext, key string) (*int64, error)
	ReadMetaReal(ctx TxContext, key string) (*float64, error)
	ReadMetaBlob(ctx TxContext, key string) (*[]byte, error)

	DeleteMeta(ctx TxContext, key string) error
}

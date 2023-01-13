package util

import (
	"blackforestbytes.com/simplecloudnotifier/logic"
	"testing"
	"time"
)

func ConvertToCompatID(t *testing.T, ws *logic.Application, newid string) int64 {

	ctx := ws.NewSimpleTransactionContext(5 * time.Second)
	defer ctx.Cancel()

	uidold, _, err := ws.Database.Primary.ConvertToCompatID(ctx, newid)
	TestFailIfErr(t, err)

	if uidold == nil {
		TestFail(t, "faile to convert newid to oldid (compat)")
	}

	err = ctx.CommitTransaction()
	if err != nil {
		TestFail(t, "failed to commit")
		return 0
	}

	return *uidold
}

func CreateCompatID(t *testing.T, ws *logic.Application, idtype string, newid string) int64 {

	ctx := ws.NewSimpleTransactionContext(5 * time.Second)
	defer ctx.Cancel()

	uidold, err := ws.Database.Primary.CreateCompatID(ctx, idtype, newid)
	TestFailIfErr(t, err)

	err = ctx.CommitTransaction()
	if err != nil {
		TestFail(t, "failed to commit")
		return 0
	}

	return uidold
}

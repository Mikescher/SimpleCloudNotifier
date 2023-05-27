package util

import (
	"blackforestbytes.com/simplecloudnotifier/logic"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
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

func ConvertCompatID(t *testing.T, ws *logic.Application, oldid int64, idtype string) string {

	ctx := ws.NewSimpleTransactionContext(5 * time.Second)
	defer ctx.Cancel()

	idnew, err := ws.Database.Primary.ConvertCompatID(ctx, oldid, idtype)
	TestFailIfErr(t, err)

	if idnew == nil {
		TestFail(t, "faile to convert oldid to newid (compat)")
	}

	err = ctx.CommitTransaction()
	if err != nil {
		TestFail(t, "failed to commit")
		return ""
	}

	return *idnew
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

func SetCompatClient(t *testing.T, ws *logic.Application, cid string) {
	ctx := ws.NewSimpleTransactionContext(5 * time.Second)
	defer ctx.Cancel()

	_, err := ws.Database.Primary.DB().Exec(ctx, "INSERT INTO compat_clients (client_id) VALUES (:cid)", sq.PP{"cid": cid})
	TestFailIfErr(t, err)

	err = ctx.CommitTransaction()
	TestFailIfErr(t, err)
}

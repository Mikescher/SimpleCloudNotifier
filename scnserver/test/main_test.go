package test

import (
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if !exerr.Initialized() {
		exerr.Init(exerr.ErrorPackageConfigInit{ZeroLogErrTraces: langext.PFalse, ZeroLogAllTraces: langext.PFalse})
	}

	os.Exit(m.Run())
}

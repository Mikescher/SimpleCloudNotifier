package test

import (
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"testing"
)

func TestSearchMessageFTS(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	tt.InitDefaultData(t, ws)

	//TODO search for messages by FTS
}

//TODO test missing message-xx methods

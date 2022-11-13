package main

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api"
	"blackforestbytes.com/simplecloudnotifier/common"
	"blackforestbytes.com/simplecloudnotifier/common/ginext"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

var conf = scn.Conf

func main() {
	common.Init(conf)

	log.Info().Msg(fmt.Sprintf("Starting with config-namespace <%s>", conf.Namespace))

	sqlite, err := db.NewDatabase(context.Background(), conf)
	if err != nil {
		panic(err)
	}

	app := logic.NewApp(sqlite)

	ginengine := ginext.NewEngine(conf)

	router := api.NewRouter(app)

	app.Init(conf, ginengine)

	router.Init(ginengine)

	app.Run()
}

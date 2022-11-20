package main

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api"
	"blackforestbytes.com/simplecloudnotifier/common"
	"blackforestbytes.com/simplecloudnotifier/common/ginext"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/firebase"
	"blackforestbytes.com/simplecloudnotifier/jobs"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"fmt"
	"github.com/rs/zerolog/log"
)

var conf = scn.Conf

func main() {
	common.Init(conf)

	log.Info().Msg(fmt.Sprintf("Starting with config-namespace <%s>", conf.Namespace))

	sqlite, err := db.NewDatabase(conf)
	if err != nil {
		panic(err)
	}

	app := logic.NewApp(sqlite)

	if err := app.Migrate(); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate DB")
		return
	}

	ginengine := ginext.NewEngine(conf)

	router := api.NewRouter(app)

	fb, err := firebase.NewFirebase(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init firebase")
		return
	}

	jobRetry := jobs.NewDeliveryRetryJob(app)

	app.Init(conf, ginengine, fb, []logic.Job{jobRetry})

	router.Init(ginengine)

	app.Run()
}

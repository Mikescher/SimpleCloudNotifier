package main

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api"
	"blackforestbytes.com/simplecloudnotifier/api/ginext"
	"blackforestbytes.com/simplecloudnotifier/google"
	"blackforestbytes.com/simplecloudnotifier/jobs"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/push"
	"fmt"
	"github.com/rs/zerolog/log"
)

func main() {
	conf := scn.Conf

	scn.Init(conf)

	log.Info().Msg(fmt.Sprintf("Starting with config-namespace <%s>", conf.Namespace))

	sqlite, err := logic.NewDBPool(conf)
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

	var nc push.NotificationClient
	if conf.DummyFirebase {
		nc = push.NewDummy()
	} else {
		nc, err = push.NewFirebaseConn(conf)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to init firebase")
			return
		}
	}

	var apc google.AndroidPublisherClient
	if conf.DummyGoogleAPI {
		apc = google.NewDummy()
	} else {
		apc, err = google.NewAndroidPublisherAPI(conf)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to init google-api")
			return
		}
	}

	jobRetry := jobs.NewDeliveryRetryJob(app)

	jobReqCollector := jobs.NewRequestLogCollectorJob(app)

	app.Init(conf, ginengine, nc, apc, []logic.Job{jobRetry, jobReqCollector})

	err = router.Init(ginengine)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init router")
		return
	}

	app.Run()
}

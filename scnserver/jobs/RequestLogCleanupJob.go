package jobs

import (
	"blackforestbytes.com/simplecloudnotifier/logic"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/syncext"
	"time"
)

type RequestLogCleanupJob struct {
	app        *logic.Application
	name       string
	isRunning  *syncext.AtomicBool
	isStarted  bool
	sigChannel chan string
}

func NewRequestLogCleanupJob(app *logic.Application) *RequestLogCleanupJob {
	return &RequestLogCleanupJob{
		app:        app,
		name:       "RequestLogCleanupJob",
		isRunning:  syncext.NewAtomicBool(false),
		isStarted:  false,
		sigChannel: make(chan string),
	}
}

func (j *RequestLogCleanupJob) Start() error {
	if j.isRunning.Get() {
		return errors.New("job already running")
	}
	if j.isStarted {
		return errors.New("job was already started") // re-start after stop is not allowed
	}

	j.isStarted = true

	go j.mainLoop()

	return nil
}

func (j *RequestLogCleanupJob) Stop() {
	log.Info().Msg(fmt.Sprintf("Stopping Job [%s]", j.name))
	syncext.WriteNonBlocking(j.sigChannel, "stop")
	j.isRunning.Wait(false)
	log.Info().Msg(fmt.Sprintf("Stopped Job [%s]", j.name))
}

func (j *RequestLogCleanupJob) Running() bool {
	return j.isRunning.Get()
}

func (j *RequestLogCleanupJob) mainLoop() {
	j.isRunning.Set(true)

	var err error = nil

	for {
		interval := 1 * time.Hour

		signal, okay := syncext.ReadChannelWithTimeout(j.sigChannel, interval)
		if okay {
			if signal == "stop" {
				log.Info().Msg(fmt.Sprintf("Job [%s] received <stop> signal", j.name))
				break
			} else if signal == "run" {
				log.Info().Msg(fmt.Sprintf("Job [%s] received <run> signal", j.name))
				continue
			} else {
				log.Error().Msg(fmt.Sprintf("Received unknown job signal: <%s> in job [%s]", signal, j.name))
			}
		}

		log.Debug().Msg(fmt.Sprintf("Run job [%s]", j.name))

		t0 := time.Now()
		err = j.execute()
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("Failed to execute job [%s]: %s", j.name, err.Error()))
		} else {
			t1 := time.Now()
			log.Debug().Msg(fmt.Sprintf("Job [%s] finished successfully after %f minutes", j.name, (t1.Sub(t0)).Minutes()))
		}

	}

	log.Info().Msg(fmt.Sprintf("Job [%s] exiting main-loop", j.name))

	j.isRunning.Set(false)
}

func (j *RequestLogCleanupJob) execute() (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Error().Interface("recover", rec).Msg("Recovered panic in " + j.name)
			err = errors.New(fmt.Sprintf("Panic recovered: %v", rec))
		}
	}()

	ctx := j.app.NewSimpleTransactionContext(10 * time.Second)
	defer ctx.Cancel()

	deleted, err := j.app.Database.Requests.Cleanup(ctx, j.app.Config.ReqLogHistoryMaxCount, j.app.Config.ReqLogHistoryMaxDuration)
	if err != nil {
		return err
	}

	log.Warn().Msgf("Deleted %d entries from the request-log table", deleted)

	return nil
}

package jobs

import (
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/syncext"
	"time"
)

type RequestLogCollectorJob struct {
	app        *logic.Application
	name       string
	isRunning  *syncext.AtomicBool
	isStarted  bool
	sigChannel chan string
}

func NewRequestLogCollectorJob(app *logic.Application) *RequestLogCollectorJob {
	return &RequestLogCollectorJob{
		app:        app,
		name:       "RequestLogCollectorJob",
		isRunning:  syncext.NewAtomicBool(false),
		isStarted:  false,
		sigChannel: make(chan string),
	}
}

func (j *RequestLogCollectorJob) Start() error {
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

func (j *RequestLogCollectorJob) Stop() {
	log.Info().Msg(fmt.Sprintf("Stopping Job [%s]", j.name))
	syncext.WriteNonBlocking(j.sigChannel, "stop")
	j.isRunning.Wait(false)
	log.Info().Msg(fmt.Sprintf("Stopped Job [%s]", j.name))
}

func (j *RequestLogCollectorJob) Running() bool {
	return j.isRunning.Get()
}

func (j *RequestLogCollectorJob) mainLoop() {
	j.isRunning.Set(true)

mainLoop:
	for {
		select {
		case signal := <-j.sigChannel:
			if signal == "stop" {
				log.Info().Msg(fmt.Sprintf("Job [%s] received <stop> signal", j.name))
				break mainLoop
			} else if signal == "run" {
				log.Info().Msg(fmt.Sprintf("Job [%s] received <run> signal", j.name))
				continue
			} else {
				log.Error().Msg(fmt.Sprintf("Received unknown job signal: <%s> in job [%s]", signal, j.name))
			}
		case obj := <-j.app.RequestLogQueue:
			err := j.insertLog(obj)
			if err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("Failed to insert RequestLog {%s} into DB", obj.RequestID))
			} else {
				log.Debug().Msg(fmt.Sprintf("Inserted RequestLog '%s' into DB", obj.RequestID))
			}
		}
	}

	log.Info().Msg(fmt.Sprintf("Job [%s] exiting main-loop", j.name))

	j.isRunning.Set(false)
}

func (j *RequestLogCollectorJob) insertLog(rl models.RequestLog) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := j.app.Database.Requests.InsertRequestLog(ctx, rl.DB())
	if err != nil {
		return err
	}

	return nil
}

package jobs

import (
	"blackforestbytes.com/simplecloudnotifier/db/simplectx"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/syncext"
	"time"
)

type DeliveryRetryJob struct {
	app        *logic.Application
	name       string
	isRunning  *syncext.AtomicBool
	isStarted  bool
	sigChannel chan string
}

func NewDeliveryRetryJob(app *logic.Application) *DeliveryRetryJob {
	return &DeliveryRetryJob{
		app:        app,
		name:       "DeliveryRetryJob",
		isRunning:  syncext.NewAtomicBool(false),
		isStarted:  false,
		sigChannel: make(chan string, 1),
	}
}

func (j *DeliveryRetryJob) Start() error {
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

func (j *DeliveryRetryJob) Stop() {
	log.Info().Msg(fmt.Sprintf("Stopping Job [%s]", j.name))
	if !syncext.WriteNonBlocking(j.sigChannel, "stop") {
		log.Error().Msg(fmt.Sprintf("Failed to send Stop-Signal to Job [%s]", j.name))
	}
	j.isRunning.Wait(false)
	log.Info().Msg(fmt.Sprintf("Stopped Job [%s]", j.name))
}

func (j *DeliveryRetryJob) Running() bool {
	return j.isRunning.Get()
}

func (j *DeliveryRetryJob) mainLoop() {
	j.isRunning.Set(true)

	var fastRerun bool = false
	var err error = nil

	for {
		interval := 30 * time.Second
		if fastRerun {
			interval = 1 * time.Second
		}

		signal, okay := syncext.ReadChannelWithTimeout(j.sigChannel, interval)
		if okay {
			if signal == "stop" {
				log.Info().Msg(fmt.Sprintf("Job [%s] received <stop> signal", j.name))
				break
			} else if signal == "run" {
				log.Info().Msg(fmt.Sprintf("Job [%s] received <run> signal", j.name))
			} else {
				log.Error().Msg(fmt.Sprintf("Received unknown job signal: <%s> in job [%s]", signal, j.name))
			}
		}

		log.Debug().Msg(fmt.Sprintf("Run job [%s]", j.name))

		t0 := time.Now()
		fastRerun, err = j.execute()
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

func (j *DeliveryRetryJob) execute() (fastrr bool, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Error().Interface("recover", rec).Msg("Recovered panic in " + j.name)
			err = errors.New(fmt.Sprintf("Panic recovered: %v", rec))
			fastrr = false
		}
	}()

	ctx := j.app.NewSimpleTransactionContext(10 * time.Second)
	defer ctx.Cancel()

	deliveries, err := j.app.Database.Primary.ListRetrieableDeliveries(ctx, 32)
	if err != nil {
		return false, err
	}

	err = ctx.CommitTransaction()
	if err != nil {
		return false, err
	}

	if len(deliveries) == 32 {
		log.Warn().Msg("The delivery pipeline is greater than 32 (too much for a single cycle)")
	}

	for _, delivery := range deliveries {
		j.redeliver(ctx, delivery)
	}

	return len(deliveries) == 32, nil
}

func (j *DeliveryRetryJob) redeliver(ctx *simplectx.SimpleContext, delivery models.Delivery) {

	client, err := j.app.Database.Primary.GetClient(ctx, delivery.ReceiverUserID, delivery.ReceiverClientID)
	if err != nil {
		log.Err(err).Str("ReceiverUserID", delivery.ReceiverUserID.String()).Str("ReceiverClientID", delivery.ReceiverClientID.String()).Msg("Failed to get client")
		ctx.RollbackTransaction()
		return
	}

	msg, err := j.app.Database.Primary.GetMessage(ctx, delivery.MessageID, true)
	if err != nil {
		log.Err(err).Str("MessageID", delivery.MessageID.String()).Msg("Failed to get message")
		ctx.RollbackTransaction()
		return
	}

	if msg.Deleted {
		err = j.app.Database.Primary.SetDeliveryFailed(ctx, delivery)
		if err != nil {
			log.Err(err).Str("MessageID", delivery.MessageID.String()).Str("DeliveryID", delivery.DeliveryID.String()).Msg("Failed to update delivery")
			ctx.RollbackTransaction()
			return
		}
	} else {

		isCompatClient, err := j.app.Database.Primary.IsCompatClient(ctx, client.ClientID)
		if err != nil {
			log.Err(err).Str("MessageID", delivery.MessageID.String()).Str("ClientID", client.ClientID.String()).Msg("Failed to get <IsCompatClient>")
			ctx.RollbackTransaction()
			return
		}

		var titleOverride *string = nil
		var msgidOverride *string = nil
		if isCompatClient {

			messageIdComp, err := j.app.Database.Primary.ConvertToCompatIDOrCreate(ctx, "messageid", msg.MessageID.String())
			if err != nil {
				log.Err(err).Str("MessageID", delivery.MessageID.String()).Str("ClientID", client.ClientID.String()).Msg("Failed to query/create messageid")
				ctx.RollbackTransaction()
				return
			}

			titleOverride = langext.Ptr(j.app.CompatizeMessageTitle(ctx, msg))
			msgidOverride = langext.Ptr(fmt.Sprintf("%d", messageIdComp))
		}

		fcmDelivID, err := j.app.DeliverMessage(ctx, client, msg, titleOverride, msgidOverride)
		if err == nil {
			err = j.app.Database.Primary.SetDeliverySuccess(ctx, delivery, fcmDelivID)
			if err != nil {
				log.Err(err).Str("MessageID", delivery.MessageID.String()).Str("DeliveryID", delivery.DeliveryID.String()).Msg("Failed to update delivery")
				ctx.RollbackTransaction()
				return
			}
		} else if delivery.RetryCount+1 > delivery.MaxRetryCount() {
			err = j.app.Database.Primary.SetDeliveryFailed(ctx, delivery)
			if err != nil {
				log.Err(err).Str("MessageID", delivery.MessageID.String()).Str("DeliveryID", delivery.DeliveryID.String()).Msg("Failed to update delivery")
				ctx.RollbackTransaction()
				return
			}
			log.Warn().Str("MessageID", delivery.MessageID.String()).Str("DeliveryID", delivery.DeliveryID.String()).Msg("Delivery failed after <max> retries (set to FAILURE)")
		} else {
			err = j.app.Database.Primary.SetDeliveryRetry(ctx, delivery)
			if err != nil {
				log.Err(err).Str("MessageID", delivery.MessageID.String()).Str("DeliveryID", delivery.DeliveryID.String()).Msg("Failed to update delivery")
				ctx.RollbackTransaction()
				return
			}
		}

	}

	err = ctx.CommitTransaction()

}

package jobs

import (
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"github.com/rs/zerolog/log"
	"time"
)

type DeliveryRetryJob struct {
	app         *logic.Application
	running     bool
	stopChannel chan bool
}

func NewDeliveryRetryJob(app *logic.Application) *DeliveryRetryJob {
	return &DeliveryRetryJob{
		app:         app,
		running:     true,
		stopChannel: make(chan bool, 8),
	}
}

func (j *DeliveryRetryJob) Start() {
	if !j.running {
		panic("cannot re-start job")
	}

	go j.mainLoop()
}

func (j *DeliveryRetryJob) Stop() {
	j.running = false
}

func (j *DeliveryRetryJob) mainLoop() {
	fastRerun := false

	for j.running {
		if fastRerun {
			j.sleep(1 * time.Second)
		} else {
			j.sleep(30 * time.Second)
		}
		if !j.running {
			return
		}

		fastRerun = j.run()

	}
}

func (j *DeliveryRetryJob) run() bool {
	defer func() {
		if rec := recover(); rec != nil {
			log.Error().Interface("recover", rec).Msg("Recovered panic in DeliveryRetryJob")
		}
	}()

	ctx := j.app.NewSimpleTransactionContext(10 * time.Second)
	defer ctx.Cancel()

	deliveries, err := j.app.Database.ListRetrieableDeliveries(ctx, 32)
	if err != nil {
		log.Err(err).Msg("Failed to query retrieable deliveries")
		return false
	}

	err = ctx.CommitTransaction()
	if err != nil {
		log.Err(err).Msg("Failed to commit")
		return false
	}

	if len(deliveries) == 32 {
		log.Warn().Msg("The delivery pipeline is greater than 32 (too much for a single cycle)")
	}

	for _, delivery := range deliveries {
		j.redeliver(ctx, delivery)
	}

	return len(deliveries) == 32
}

func (j *DeliveryRetryJob) redeliver(ctx *logic.SimpleContext, delivery models.Delivery) {

	client, err := j.app.Database.GetClient(ctx, delivery.ReceiverUserID, delivery.ReceiverClientID)
	if err != nil {
		log.Err(err).Int64("ReceiverUserID", delivery.ReceiverUserID).Int64("ReceiverClientID", delivery.ReceiverClientID).Msg("Failed to get client")
		ctx.RollbackTransaction()
		return
	}

	msg, err := j.app.Database.GetMessage(ctx, delivery.SCNMessageID)
	if err != nil {
		log.Err(err).Int64("SCNMessageID", delivery.SCNMessageID).Msg("Failed to get message")
		ctx.RollbackTransaction()
		return
	}

	fcmDelivID, err := j.app.DeliverMessage(ctx, client, msg)
	if err == nil {
		err = j.app.Database.SetDeliverySuccess(ctx, delivery, *fcmDelivID)
		if err != nil {
			log.Err(err).Int64("SCNMessageID", delivery.SCNMessageID).Int64("DeliveryID", delivery.DeliveryID).Msg("Failed to update delivery")
			ctx.RollbackTransaction()
			return
		}
	} else if delivery.RetryCount+1 > delivery.MaxRetryCount() {
		err = j.app.Database.SetDeliveryFailed(ctx, delivery)
		if err != nil {
			log.Err(err).Int64("SCNMessageID", delivery.SCNMessageID).Int64("DeliveryID", delivery.DeliveryID).Msg("Failed to update delivery")
			ctx.RollbackTransaction()
			return
		}
		log.Warn().Int64("SCNMessageID", delivery.SCNMessageID).Int64("DeliveryID", delivery.DeliveryID).Msg("Delivery failed after <max> retries (set to FAILURE)")
	} else {
		err = j.app.Database.SetDeliveryRetry(ctx, delivery)
		if err != nil {
			log.Err(err).Int64("SCNMessageID", delivery.SCNMessageID).Int64("DeliveryID", delivery.DeliveryID).Msg("Failed to update delivery")
			ctx.RollbackTransaction()
			return
		}
	}

	err = ctx.CommitTransaction()

}

func (j *DeliveryRetryJob) sleep(d time.Duration) {
	if !j.running {
		return
	}
	afterCh := time.After(d)
	for {
		select {
		case <-j.stopChannel:
			j.stopChannel <- true
			return
		case <-afterCh:
			return
		}
	}
}

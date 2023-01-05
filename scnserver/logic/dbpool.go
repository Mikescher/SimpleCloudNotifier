package logic

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/db"
	logsdb "blackforestbytes.com/simplecloudnotifier/db/impl/logs"
	primarydb "blackforestbytes.com/simplecloudnotifier/db/impl/primary"
	requestsdb "blackforestbytes.com/simplecloudnotifier/db/impl/requests"
	"context"
)

type DBPool struct {
	Primary  *primarydb.Database
	Requests *requestsdb.Database
	Logs     *logsdb.Database
}

func NewDBPool(conf scn.Config) (*DBPool, error) {

	dbprimary, err := primarydb.NewPrimaryDatabase(conf)
	if err != nil {
		return nil, err
	}

	dbrequests, err := requestsdb.NewRequestsDatabase(conf)
	if err != nil {
		return nil, err
	}

	dblogs, err := logsdb.NewLogsDatabase(conf)
	if err != nil {
		return nil, err
	}

	return &DBPool{
		Primary:  dbprimary,
		Requests: dbrequests,
		Logs:     dblogs,
	}, nil
}

func (p DBPool) List() []db.DatabaseImpl {
	return []db.DatabaseImpl{
		p.Primary,
		p.Requests,
		p.Logs,
	}
}

func (p DBPool) Stop(ctx context.Context) error {

	var err error = nil

	for _, subdb := range p.List() {
		err2 := subdb.Stop(ctx)
		if err2 != nil && err == nil {
			err = err2
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (p DBPool) Migrate(ctx context.Context) error {
	for _, subdb := range p.List() {
		err := subdb.Migrate(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p DBPool) Ping(ctx context.Context) error {
	for _, subdb := range p.List() {
		err := subdb.Ping(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

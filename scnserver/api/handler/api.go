package handler

import (
	primarydb "blackforestbytes.com/simplecloudnotifier/db/impl/primary"
	"blackforestbytes.com/simplecloudnotifier/logic"
)

type APIHandler struct {
	app      *logic.Application
	database *primarydb.Database
}

func NewAPIHandler(app *logic.Application) APIHandler {
	return APIHandler{
		app:      app,
		database: app.Database.Primary,
	}
}

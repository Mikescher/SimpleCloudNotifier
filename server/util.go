package server

import (
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"time"
)

func QuotaDayString() string {
	return time.Now().In(timeext.TimezoneBerlin).Format("2006-01-02")
}

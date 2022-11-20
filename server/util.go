package server

import (
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"time"
)

func QuotaDayString() string {
	return time.Now().In(timeext.TimezoneBerlin).Format("2006-01-02")
}

func NextDeliveryTimestamp(now time.Time) time.Time {
	return now.Add(5 * time.Second)
}

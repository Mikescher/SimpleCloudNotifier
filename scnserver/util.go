package server

import (
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"math/rand"
	"time"
)

func QuotaDayString() string {
	return time.Now().In(timeext.TimezoneBerlin).Format("2006-01-02")
}

func NextDeliveryTimestamp(now time.Time) time.Time {
	return now.Add(5 * time.Second)
}

func RandomAuthKey() string {
	charset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	k := ""
	for i := 0; i < 64; i++ {
		k += string(charset[rand.Int()%len(charset)])
	}
	return k
}

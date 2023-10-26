package models

import (
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

func timeOptFmt(t *time.Time, fmt string) *string {
	if t == nil {
		return nil
	} else {
		return langext.Ptr(t.Format(fmt))
	}
}

func timeOptFromMilli(millis *int64) *time.Time {
	if millis == nil {
		return nil
	}
	return langext.Ptr(time.UnixMilli(*millis))
}

func timeFromMilli(millis int64) time.Time {
	return time.UnixMilli(millis)
}

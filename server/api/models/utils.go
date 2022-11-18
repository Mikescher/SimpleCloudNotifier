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

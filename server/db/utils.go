package db

import "time"

func bool2DB(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func time2DB(t time.Time) int64 {
	return t.UnixMilli()
}

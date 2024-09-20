package models

type SenderNameStatistics struct {
	SenderName     string  `json:"name"            db:"name"`
	LastTimestamp  SCNTime `json:"last_timestamp"  db:"ts_last"`
	FirstTimestamp SCNTime `json:"first_timestamp" db:"ts_first"`
	Count          int     `json:"count"           db:"count"`
}

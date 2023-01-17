package models

type CompatMessage struct {
	Title         string  `json:"title"`
	Body          *string `json:"body"`
	Priority      int     `json:"priority"`
	Timestamp     int64   `json:"timestamp"`
	UserMessageID *string `json:"usr_msg_id"`
	SCNMessageID  int64   `json:"scn_msg_id"`
	Trimmed       *bool   `json:"trimmed"`
}

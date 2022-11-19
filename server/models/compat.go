package models

type CompatMessage struct {
	Title         string `json:"title"`
	Body          string `json:"body"`
	Priority      int    `json:"priority"`
	Timestamp     int64  `json:"timestamp"`
	UserMessageID string `json:"usr_msg_id"`
	SCNMessageID  string `json:"scn_msg_id"`
}

type ShortCompatMessage struct {
	Title         string `json:"title"`
	Body          string `json:"body"`
	Trimmed       bool   `json:"trimmed"`
	Priority      int    `json:"priority"`
	Timestamp     int64  `json:"timestamp"`
	UserMessageID string `json:"usr_msg_id"`
	SCNMessageID  string `json:"scn_msg_id"`
}

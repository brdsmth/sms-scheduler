package types

import "time"

type ScheduledMessage struct {
	To       string    `json:"to"`
	Message  string    `json:"message"`
	SendTime time.Time `json:"sendTime"`
}

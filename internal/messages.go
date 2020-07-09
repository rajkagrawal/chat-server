package internal

import  "time"

type message struct {
	messageID string
	sender string
	receiver  string
	isMessageForRoom bool
	time time.Time
	isDelivered bool
}

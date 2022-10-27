package core

type MQEvent struct {
	Type    string
	Payload interface{}
}

type NewUser struct {
	ChatID int64
}

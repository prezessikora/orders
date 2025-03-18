package model

import "time"

type Order struct {
	Id      int
	UserId  int
	EventId int
	Created time.Time
	Status  string
}

// NewOrder creates Order with creation parameters
func NewOrder(eventId int, userId int) Order {
	o := Order{UserId: userId, EventId: eventId}
	o.Status = "pending"
	return o
}

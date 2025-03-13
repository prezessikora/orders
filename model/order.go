package model

import "time"

type Order struct {
	Id      int
	UserId  int
	EventId int
	Created time.Time
	Status  string
}

// Create Order to bind with request
func NewOrder(eventId int, userId int) Order {
	o := Order{UserId: userId, EventId: eventId}
	o.Created = time.Now()
	o.Status = "pending"
	return o
}

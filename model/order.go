package model

import "time"

type Order struct {
	Id      int       `json:"id"`
	UserId  int       `binding:"required" json:"user_id"`
	EventId int       `binding:"required" json:"event_id"`
	Created time.Time `json:"created"`
	Status  string    `json:"status"`
}

// Create Order to bind with request
func (o *Order) Reset() {
	o.Created = time.Now()
	o.Status = "pending"
}

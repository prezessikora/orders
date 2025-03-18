package model

import "time"

type Ticket struct {
	Id           int
	ticketNumber string
	OrderId      int
	TicketId     int
	UserId       int
	EventId      int
	eventName    string
	ticketType   string
	name         string
	validFrom    time.Time
	validThrough time.Time
	generated    time.Time
}

func NewTicket(eventId int, userId int) Ticket {
	return Ticket{EventId: eventId, generated: time.Now(), UserId: userId}
}

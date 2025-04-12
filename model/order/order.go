package order

import (
	"errors"
	"github.com/prezessikora/orders/myclient"
	"time"
)

type Order struct {
	Id      int // Id  the id of the order. Auto filed.
	UserId  int
	EventId int
	Created time.Time // Created  the timestamp when the order was created. Auto filed.
	Status  string    // Status represents the current state of the order, such as "pending", "paid", or "cancelled".
}

// NewOrder verifies if the Order can be fulfilled and when successful creates and returns the order, error otherwise
func Create(eventId int, userId int, event *myclient.Event) (Order, error) {
	newOrder := Order{UserId: userId, EventId: eventId}
	newOrder.Status = "pending"

	err := canCreateOrderForEvent(event)

	if err != nil {
		return Order{}, err
	}

	return newOrder, nil
}

// Policy expresses a rule of allowed orders to be created
type Policy interface {
	isAllowed(myclient.Event) error
}

type CapacityPolicy struct {
}

func (c CapacityPolicy) isAllowed(event myclient.Event) error {
	//TODO implement me
	panic("implement me")
}

type BookingWindowPolicy struct {
}

func (b BookingWindowPolicy) isAllowed(event myclient.Event) error {
	//TODO implement me
	panic("implement me")
}

// canCreateOrderForEvent checks if an order can be created for the given event based on its capacity and start time.
func canCreateOrderForEvent(event *myclient.Event) error {

	if !(event.Capacity > event.Bookings) {
		return errors.New("event has no capacity")
	}

	now := time.Now()
	diff := event.DateTime.Sub(now)
	if diff.Truncate(time.Hour) > time.Hour*24 {
		return errors.New("event start time is less than 24h from now")
	}
	return nil
}

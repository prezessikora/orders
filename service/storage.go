package service

import (
	"github.com/prezessikora/orders/model"
	"github.com/prezessikora/orders/model/order"
)

// Both repository interfaces declare the http service contract expectation

// OrderDataStorage is orders data store interface for various storages in-mem or sql.
type OrderDataStorage interface {
	AddOrder(order order.Order) int
	GetAll() []order.Order
	GetUserOrders(userId int) []order.Order
	GetOrderById(id int) (order.Order, error)
	CancelEventOrders(eventId int) error
}

// TicketsDataStorage OrderDataStorage is tickets data store interface for various storages in-mem or sql.
// Follows repository per bounded-context pattern and is a candidate for future DB extraction.
type TicketsDataStorage interface {
	GetAll(int) []model.Ticket
	AddTicket(model.Ticket)
	GetUserTickets(userId int) []model.Ticket
}

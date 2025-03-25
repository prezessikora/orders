package service

import "github.com/prezessikora/orders/model"

// Both repository interfaces declare the http service contract expectation

// OrderDataStorage is orders data store interface for various storages in-mem or sql.
type OrderDataStorage interface {
	AddOrder(order model.Order) int
	GetAll() []model.Order
	GetUserOrders(userId int) []model.Order
	GetOrderById(id int) (model.Order, error)
	CancelEventOrders(eventId int) error
}

// TicketsDataStorage OrderDataStorage is tickets data store interface for various storages in-mem or sql.
// Follows repository per bounded-context pattern and is a candidate for future DB extraction.
type TicketsDataStorage interface {
	GetAll(int) []model.Ticket
	AddTicket(model.Ticket)
	GetUserTickets(userId int) []model.Ticket
}

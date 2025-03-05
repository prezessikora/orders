package storage

import (
	"com.sikora/payments/orders"
	"errors"
)

// MemoryOrdersStorage is in-memory storage for orders
type MemoryOrdersStorage struct {
	orders []orders.Order
	nextId int
}

var Memory = MemoryOrdersStorage{}

func (storage *MemoryOrdersStorage) AddOrder(order orders.Order) orders.Order {
	storage.nextId += 1
	order.Id = storage.nextId
	storage.orders = append(storage.orders, order)

	return order
}

func (storage MemoryOrdersStorage) GetAll() []orders.Order {
	return storage.orders
}

func (storage MemoryOrdersStorage) GetOrderById(id int) (orders.Order, error) {
	for _, order := range storage.orders {
		if order.UserId == id {
			return order, nil
		}
	}
	return orders.Order{}, errors.New("could not find order")
}

type OrdersStorage interface {
}

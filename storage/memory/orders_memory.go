package memory

import (
	"errors"
	"github.com/prezessikora/orders/model/order"
	"time"
)

func NewDataStore() *OrdersStorage {
	mos := OrdersStorage{orders: make([]order.Order, 0, 10)}
	mos.AddOrder(order.Order{Id: 0, UserId: 1, EventId: 1, Created: time.Now(), Status: "pending"})
	return &mos
}

type OrdersStorage struct {
	orders []order.Order
	nextId int
}

func (storage *OrdersStorage) CancelEventOrders(eventId int) error {
	for _, order := range storage.orders {
		if order.EventId == eventId {
			order.Status = "canceled"
		}
	}
	return nil
}

func (storage *OrdersStorage) GetUserOrders(userId int) []order.Order {
	var result []order.Order
	for _, storedOrder := range storage.orders {
		if storedOrder.UserId == userId {
			result = append(result, storedOrder)
		}
	}
	return result
}

func (storage *OrdersStorage) AddOrder(order order.Order) int {
	storage.nextId += 1
	order.Id = storage.nextId
	storage.orders = append(storage.orders, order)
	return order.Id
}

func (storage OrdersStorage) GetAll() []order.Order {
	return storage.orders
}

func (storage OrdersStorage) GetOrderById(id int) (order.Order, error) {
	for _, order := range storage.orders {
		if order.Id == id {
			return order, nil
		}
	}
	return order.Order{}, errors.New("could not find order")
}

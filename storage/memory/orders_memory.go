package memory

import (
	"errors"
	"github.com/prezessikora/orders/model"
	"time"
)

func NewDataStore() *OrdersStorage {
	mos := OrdersStorage{orders: make([]model.Order, 0, 10)}
	mos.AddOrder(model.Order{Id: 0, UserId: 1, EventId: 1, Created: time.Now(), Status: "pending"})
	return &mos
}

type OrdersStorage struct {
	orders []model.Order
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

func (storage *OrdersStorage) GetUserOrders(userId int) []model.Order {
	var result []model.Order
	for _, storedOrder := range storage.orders {
		if storedOrder.UserId == userId {
			result = append(result, storedOrder)
		}
	}
	return result
}

func (storage *OrdersStorage) AddOrder(order model.Order) int {
	storage.nextId += 1
	order.Id = storage.nextId
	storage.orders = append(storage.orders, order)
	return order.Id
}

func (storage OrdersStorage) GetAll() []model.Order {
	return storage.orders
}

func (storage OrdersStorage) GetOrderById(id int) (model.Order, error) {
	for _, order := range storage.orders {
		if order.Id == id {
			return order, nil
		}
	}
	return model.Order{}, errors.New("could not find order")
}

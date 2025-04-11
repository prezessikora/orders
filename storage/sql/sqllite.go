package sql

import (
	"errors"
	"fmt"
	"github.com/prezessikora/orders/model/order"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// SQL model

type Order struct {
	gorm.Model
	UserId  int
	EventId int
	Status  string
	ignored string // fields that aren't exported are ignored
	Log     []Log
}

type Log struct {
	gorm.Model
	OrderID     uint
	Description string
}

type DataStorage struct {
	db *gorm.DB
}

func (s DataStorage) logOrderEvent(description string, orderId int) error {
	l := Log{OrderID: uint(orderId), Description: description}
	result := s.db.Create(&l)
	if result.Error != nil {
		log.Printf("error creating order log event: %v", result.Error)
		return result.Error
	}
	return nil

}

// CancelEventOrders changes the state of all orders given eventId that was deleted in events service
// Relevant Log entries for all cancelled events are also created
func (s DataStorage) CancelEventOrders(eventId int) error {
	var changedOrders []Order
	result := s.db.Where(&Order{EventId: eventId}).Find(&changedOrders)

	if err := result.Error; err != nil {
		log.Printf("error querying orders to be cancelled: %v\n", err)
		return err
	}

	result = s.db.Model(&Order{}).Where(&Order{EventId: eventId}).Update("status", "cancelled")

	if err := result.Error; err != nil {
		log.Printf("error cancelling orders: %v\n", err)
		return err
	}
	log.Printf("cancelled %v events", result.RowsAffected)
	for _, o := range changedOrders {
		log.Printf("creating order cancel log for order_id: %v", o.ID)
		err := s.logOrderEvent("cancelled", int(o.ID))
		if err != nil {
			log.Printf("error creating order change log: %v", err)
		}
	}
	return nil
}

func (s DataStorage) GetUserOrders(userId int) []order.Order {
	var dbOrders []Order
	result := s.db.Where(&Order{UserId: userId}).Find(&dbOrders)

	if result.Error != nil {
		log.Println(result.Error)
		return make([]order.Order, 0)
	}
	// map sql data store objects to the model
	orders := make([]order.Order, 0, len(dbOrders))
	for _, row := range dbOrders {
		orders = append(orders, order.Order{
			Id:      int(row.ID),
			UserId:  row.UserId,
			EventId: row.EventId,
			Created: row.CreatedAt,
			Status:  row.Status,
		})
	}
	return orders
}

func NewDataStore() (error, *DataStorage) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open("orders.db"), &gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println(err)
		return err, nil
	}

	err = db.AutoMigrate(&Order{}, &Log{})
	if err != nil {
		fmt.Println(err)
		return err, nil
	}

	return nil, &DataStorage{db: db}

	//
	//db.Find(&orders)
	//fmt.Printf("query all: %d\n", len(orders))
}

func (s DataStorage) AddOrder(order order.Order) int {

	dbOrder := Order{UserId: order.UserId, EventId: order.EventId, Status: order.Status,
		Log: []Log{{Description: "order created"}, {Description: "payment requested"}}}
	result := s.db.Create(&dbOrder)

	if err := result.Error; err != nil {
		log.Println(result.Error)
		// TODO return error here
		return 0
	}

	log.Printf("added order with id: %d", dbOrder.ID)
	return int(dbOrder.ID)
}

func (s DataStorage) GetAll() []order.Order {
	var dbOrders []Order
	result := s.db.Preload("Log").Find(&dbOrders)
	if result.Error != nil {
		log.Println(result.Error)
		return make([]order.Order, 0)
	}
	// map sql data store objects to the model
	orders := make([]order.Order, 0, len(dbOrders))
	for _, row := range dbOrders {
		latestStatus := row.Log[len(row.Log)-1]

		orders = append(orders, order.Order{
			Id:      int(row.ID),
			UserId:  row.UserId,
			EventId: row.EventId,
			Created: row.CreatedAt,
			Status:  latestStatus.Description,
		})
	}
	return orders
}

func (s DataStorage) GetOrderById(id int) (order.Order, error) {
	var orderById Order
	result := s.db.First(&orderById, id)
	if result.Error != nil {
		return order.Order{}, result.Error
	}
	if result.RowsAffected == 0 {
		return order.Order{}, errors.New("could not find order")
	}

	log.Printf("GetOrderById: %v\n", orderById)
	log.Println(orderById.Log)
	ro := order.Order{Id: int(orderById.ID), UserId: orderById.UserId, EventId: orderById.EventId, Created: orderById.CreatedAt, Status: orderById.Status}
	return ro, nil

}

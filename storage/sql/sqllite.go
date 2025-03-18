package sql

import (
	"errors"
	"fmt"
	"github.com/prezessikora/orders/model"
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

func (s DataStorage) GetUserOrders(userId int) []model.Order {
	var dbOrders []Order
	result := s.db.Where(&Order{UserId: userId}).Find(&dbOrders)

	if result.Error != nil {
		log.Println(result.Error)
		return make([]model.Order, 0)
	}
	// map sql data store objects to the model
	orders := make([]model.Order, 0, len(dbOrders))
	for _, row := range dbOrders {
		orders = append(orders, model.Order{
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

func (s DataStorage) AddOrder(order model.Order) int {

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

func (s DataStorage) GetAll() []model.Order {
	var dbOrders []Order
	result := s.db.Preload("Log").Find(&dbOrders)
	if result.Error != nil {
		log.Println(result.Error)
		return make([]model.Order, 0)
	}
	// map sql data store objects to the model
	orders := make([]model.Order, 0, len(dbOrders))
	for _, row := range dbOrders {
		latestStatus := row.Log[len(row.Log)-1]

		orders = append(orders, model.Order{
			Id:      int(row.ID),
			UserId:  row.UserId,
			EventId: row.EventId,
			Created: row.CreatedAt,
			Status:  latestStatus.Description,
		})
	}
	return orders
}

func (s DataStorage) GetOrderById(id int) (model.Order, error) {
	var order Order
	result := s.db.First(&order, id)
	if result.Error != nil {
		return model.Order{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Order{}, errors.New("could not find order")
	}

	log.Printf("GetOrderById: %v\n", order)
	log.Println(order.Log)
	ro := model.Order{Id: int(order.ID), UserId: order.UserId, EventId: order.EventId, Created: order.CreatedAt, Status: order.Status}
	return ro, nil

}

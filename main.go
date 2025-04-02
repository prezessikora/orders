package main

import (
	"context"
	"fmt"
	"github.com/prezessikora/orders/notifications"
	"github.com/prezessikora/orders/service"
	"github.com/prezessikora/orders/storage/memory"
	"github.com/prezessikora/orders/storage/sql"
	_ "github.com/prezessikora/orders/storage/sql"
	"log"
	_ "log"
	"os"
)
import "github.com/gin-gonic/gin"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func memoryStoreService() service.OrderDataStorage {
	log.Println("Using memory.NewDataStore")
	return memory.NewDataStore()
}

func memoryTicketsStoreService() *memory.TicketsMemoryOrdersStorage {
	return memory.NewTicketsDataStore()
}

func sqlStorageService() service.OrderDataStorage {
	log.Println("Using sql.NewDataStore")
	err, storage := sql.NewDataStore()
	if err != nil {
		log.Fatal("error creating sql data store %v", err)
		return nil
	}
	return storage
}

func main() {

	storage := createDataStore()

	// gin HTTP routes
	server := gin.Default()
	ordersService := service.NewOrdersService(storage)
	ordersService.RegisterRoutes(server)
	// tickets
	ticketsService := service.NewTicketsService(memoryTicketsStoreService(), ordersService)
	ticketsService.RegisterRoutes(server)

	// async message broker notifications
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notifications.HandleEventsDeletions(ctx, storage)

	err := server.Run(":8081")
	if err != nil {
		log.Println(err)
		return
	}

}

func createDataStore() service.OrderDataStorage {
	if len(os.Args) == 1 {
		return memoryStoreService()
	} else { // switch on second arg value
		switch os.Args[1] {
		case "sql":
			return sqlStorageService()
		case "mem":
			return memoryStoreService()
		default:
			fmt.Println("Wrong command line parameters.\nRun as: orders [mem | sql]")
			log.Fatal("wrong command line parameters")
		}
	}
	return nil
}

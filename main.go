package main

import (
	"fmt"
	"github.com/prezessikora/orders/service"
	"github.com/prezessikora/orders/storage/sql"
	"log"
)
import "github.com/gin-gonic/gin"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	fmt.Printf("Orders_Manager!")

	server := gin.Default()

	//ordersService := service.NewOrdersService(memory.NewDataStore())
	// SQL
	err, storage := sql.NewDataStore()
	if err != nil {
		log.Println(err)
		log.Fatal("Error creating sql data store")
		return
	}
	ordersService := service.NewOrdersService(storage)
	// SQL

	ordersService.RegisterRoutes(server)

	err = server.Run(":8081")
	if err != nil {
		fmt.Println(err)
		return
	}

}

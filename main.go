package main

import (
	"com.sikora/orders/service"
	"com.sikora/orders/storage/memory"
	"fmt"
)
import "github.com/gin-gonic/gin"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	fmt.Printf("Orders_Manager!")

	//db, err := gorm.Open(sqlite.Open("orders.db"), &gorm.Config{})

	server := gin.Default()

	ordersService := service.NewOrdersService(memory.NewDataStore())
	ordersService.RegisterRoutes(server)

	err := server.Run(":8081")
	if err != nil {
		fmt.Println(err)
		return
	}

}

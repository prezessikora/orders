package main

import (
	"com.sikora/payments/service"
	"fmt"
)
import "github.com/gin-gonic/gin"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	fmt.Printf("Orders_Manager!")
	server := gin.Default()
	service.RegisterRoutes(server)

	server.Run(":8081")

}

package service

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/orders", createOrder)
	server.GET("/orders/", getAllPaidOrders)
	server.GET("/orders/:id", orderStatus)
}

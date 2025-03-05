package service

import (
	"com.sikora/payments/orders"
	"com.sikora/payments/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Initiate the registration createOrder for the given event and user
func createOrder(ctx *gin.Context) {
	var order orders.Order

	err := ctx.ShouldBindBodyWithJSON(&order)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not order request data"})
		fmt.Println(err)
		return
	}
	order.Reset()
	order = storage.Memory.AddOrder(order)
	fmt.Println(order)
	ctx.JSONP(http.StatusCreated, gin.H{"order": order})
}

// returns all paid orders for which ticket can be created
func getAllPaidOrders(context *gin.Context) {
	all := storage.Memory.GetAll()
	context.JSONP(http.StatusOK, gin.H{"orders": all})
}

// check order status for user dashboard
func orderStatus(ctx *gin.Context) {
	ctx.JSONP(http.StatusCreated, gin.H{"orderId": 123, "status": "pending"})
	idParam := ctx.Param("id")
	orderId, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not convert request param"})
		return
	}

	order, err := storage.MemoryOrdersStorage{}.GetOrderById(orderId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusNotFound, gin.H{"message": "could not find order"})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"order": order})
}

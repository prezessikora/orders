package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prezessikora/events/client"
	"github.com/prezessikora/orders/model"
	"log"
	"net/http"
	"strconv"
)

// OrdersService The service and its deps
type OrdersService struct {
	storage OrderDataStorage
}

func NewOrdersService(storage OrderDataStorage) *OrdersService {
	return &OrdersService{storage: storage}
}

type OrderRequest struct {
	UserId  int `binding:"required" json:"user_id"`
	EventId int `binding:"required" json:"event_id"`
}

// Initiate the registration createOrder for the given event and user
func (service OrdersService) createOrder(ctx *gin.Context) {
	var orderRequest OrderRequest

	err := ctx.ShouldBindBodyWithJSON(&orderRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse orderRequest request data"})
		fmt.Println(err)
		return
	}
	// check with events service the order is valid
	event, err := client.NewEvents().GetEvent(orderRequest.EventId)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "order event not verified"})
		log.Print(err)
		return
	}
	log.Print(event)
	order := model.NewOrder(orderRequest.EventId, orderRequest.UserId)
	orderId := service.storage.AddOrder(order)

	ctx.JSONP(http.StatusCreated, gin.H{"order_id": orderId})
}

// returns all paid model for which ticket can be created
func (service OrdersService) getAllUserPaidOrdersHttp(ctx *gin.Context) {
	idParam := ctx.Param("userId")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not convert request param"})
		return
	}

	all := service.getAllPaidOrders(userId)
	ctx.JSONP(http.StatusOK, gin.H{"orders": all})
}

// Inter service interface so that they dont have to go through HTTP
func (service OrdersService) getAllPaidOrders(userId int) []model.Order {
	return service.storage.GetUserOrders(userId)
}

// check order status for user dashboard
func (service OrdersService) orderStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	orderId, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not convert request param"})
		return
	}

	order, err := service.storage.GetOrderById(orderId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusNotFound, gin.H{"message": "could not find order"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"order": order})
}

func (service OrdersService) RegisterRoutes(server *gin.Engine) {
	server.POST("/orders", service.createOrder)
	server.GET("/orders/user/:userId", service.getAllUserPaidOrdersHttp)
	server.GET("/orders/:id", service.orderStatus)
}

package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prezessikora/orders/model"
	"net/http"
	"strconv"
)

// The service and its deps
type OrdersService struct {
	storage OrderDataStorage
}

func NewOrdersService(storage OrderDataStorage) *OrdersService {
	return &OrdersService{storage: storage}
}

// Data store interface for various storages, interface on client side!
type OrderDataStorage interface {
	AddOrder(order model.Order) int
	GetAll() []model.Order
	GetOrderById(id int) (model.Order, error)
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
	order := model.NewOrder(orderRequest.EventId, orderRequest.UserId)
	orderId := service.storage.AddOrder(order)

	ctx.JSONP(http.StatusCreated, gin.H{"order_id": orderId})
}

// returns all paid model for which ticket can be created
func (service OrdersService) getAllPaidOrders(context *gin.Context) {
	all := service.storage.GetAll()
	context.JSONP(http.StatusOK, gin.H{"orders": all})
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
	server.GET("/orders/", service.getAllPaidOrders)
	server.GET("/orders/:id", service.orderStatus)
}

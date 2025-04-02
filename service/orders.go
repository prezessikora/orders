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
	s, ok := ctx.Get(XEventsHeaderKey)
	var coorelationId string
	if ok {
		s2, _ := s.(string)
		coorelationId = s2
	}

	log.Printf("[INFO] [Orders] [%v] create order", coorelationId)

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

// Inter service interface for tickets so that they dont have to go through HTTP
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

	correlated := server.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	correlated.Use(correlationId)

	correlated.POST("/orders", service.createOrder)
	correlated.GET("/orders/user/:userId", service.getAllUserPaidOrdersHttp)
	correlated.GET("/orders/:id", service.orderStatus)

}

const XEventsHeaderKey = "X-Events-Request-Id"

func correlationId(c *gin.Context) {

	// this header is used as coorelation id for requests

	header := c.GetHeader(XEventsHeaderKey)
	if header == "" {
		id := "XXX-100"
		// further to be set as HTTP header for in between services calls
		c.Set(XEventsHeaderKey, id)
		log.Printf("setting coorelation id [%v] : %v", XEventsHeaderKey, id)
	}
	// before request
	c.Next()
	// after request

}

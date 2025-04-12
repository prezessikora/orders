package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prezessikora/orders/common"
	"github.com/prezessikora/orders/model/order"
	"github.com/prezessikora/orders/myclient"
	"log/slog"

	"log"
	"net/http"
	"strconv"
)

type EventPublisher struct {
}

func (p EventPublisher) notifyOrderCreated(order order.Order) {
	// TODO switch notification code here
}

// OrdersService The service and its deps
type OrdersService struct {
	storage              OrderDataStorage
	domainEventPublisher EventPublisher
}

type LogLevel int

// Logger is facade for slog with extra features to enhance the log output with request id and alike
type Logger struct {
}

func (logger Logger) Info(message string) {
	slog.Info(message)
}

func (logger Logger) Warn(message string) {
	slog.Warn(message)
}

func NewOrdersService(storage OrderDataStorage) *OrdersService {
	return &OrdersService{storage: storage, domainEventPublisher: EventPublisher{}}
}

type OrderRequest struct {
	UserId  int `binding:"required" json:"user_id"`
	EventId int `binding:"required" json:"event_id"`
}

// Initiate the registration createOrder for the given event and user
func (service OrdersService) createOrder(ctx *gin.Context) {

	// TODO move this into logger as struct on the service since the UUID is only needed there at the moment

	correlationId, err := common.GetRequestUUIDFromContext(ctx)
	if err != nil {
		sendResponse(ctx, http.StatusInternalServerError, "request id missing in service layer", err)
		return
	}
	log.Printf("[INFO] [Orders] [%v] create newOrder", correlationId)

	// parse the request
	var orderRequest OrderRequest
	err = ctx.ShouldBindBodyWithJSON(&orderRequest)
	if err != nil {
		sendResponse(ctx, http.StatusBadRequest, "incorrect request", err)
		return
	}

	// fetch newOrder event
	event, err := myclient.NewEventsServiceClient().GetEvent(orderRequest.EventId, ctx)
	if err != nil {
		sendErrorResponse(ctx, "could not fetch event to create newOrder", err)
		return
	}

	// MODEL interaction
	// check with events service the newOrder is valid & there are enough places
	newOrder, err := order.Create(orderRequest.EventId, orderRequest.UserId, event)
	if err != nil {
		sendResponse(ctx, http.StatusBadRequest, "could not create newOrder", err)
		return
	}
	// persist the newOrder get obtain it's id
	orderId := service.storage.AddOrder(newOrder)

	newOrder.Id = orderId
	service.domainEventPublisher.notifyOrderCreated(newOrder)
	// return created newOrder to the client
	ctx.JSON(http.StatusCreated, gin.H{"newOrder": newOrder})
}

func sendErrorResponse(ctx *gin.Context, msg string, err error) {
	sendResponse(ctx, http.StatusInternalServerError, msg, err)
}

func sendResponse(ctx *gin.Context, code int, msg string, err error) {
	ctx.JSON(code, gin.H{"message": msg, "error": fmt.Sprint(err)})
	log.Println(err) // TODO log ERROR it
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
func (service OrdersService) getAllPaidOrders(userId int) []order.Order {
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

	orderById, err := service.storage.GetOrderById(orderId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusNotFound, gin.H{"message": "could not find orderById"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"orderById": orderById})
}

func (service OrdersService) RegisterRoutes(server *gin.Engine) {

	correlated := server.Group("/orders")
	// per group common! in this case we use the custom created
	// AuthRequired() common just in the "authorized" group.
	correlated.Use(common.CorrelationId)

	correlated.POST("/", service.createOrder)
	correlated.GET("/user/:userId", service.getAllUserPaidOrdersHttp)
	correlated.GET("/:id", service.orderStatus)

}

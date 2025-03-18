package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prezessikora/orders/model"
	"net/http"
	"slices"
	"strconv"
)

// TicketsService generates and offers tickets download for fully purchased orders
type TicketsService struct {
	storage       TicketsDataStorage
	ordersService *OrdersService
}

func NewTicketsService(storage TicketsDataStorage, ordersService *OrdersService) TicketsService {
	return TicketsService{
		storage:       storage,
		ordersService: ordersService,
	}
}

// returns all elements in set1 that are not in set2
func subtract(set1 []int, set2 []int) []int {
	result := make([]int, 0, len(set1))

	for _, e := range set1 {
		if found := slices.Index(set2, e); found > -1 {
			continue
		} else {
			result = append(result, e)
		}
	}
	return result
}

// getAllTickets returns all fully purchased orders tickets
// ticket can only be generated once so the method checks if there are new paid orders
// and generates the tickets
func (tickets TicketsService) getAllUserTickets(ctx *gin.Context) {
	idParam := ctx.Param("userId")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not convert userId request param"})
		return
	}

	paidOrders := tickets.ordersService.getAllPaidOrders(userId)
	paidSet := make([]int, 0, len(paidOrders))
	for _, order := range paidOrders {
		paidSet = append(paidSet, order.EventId)
	}

	generatedTickets := tickets.storage.GetAll(userId)
	generatedSet := make([]int, 0, len(generatedTickets))
	for _, ticket := range generatedTickets {
		generatedSet = append(generatedSet, ticket.EventId)
	}

	fmt.Println(paidSet)
	fmt.Println(generatedSet)
	newTickets := subtract(paidSet, generatedSet)
	fmt.Println(newTickets)

	// TODO reach out to events for event details
	tickets.generateNewTickets(newTickets, userId)

	currentTickets := tickets.storage.GetAll(userId)
	ctx.JSONP(http.StatusOK, gin.H{"printable_user_tickets": currentTickets})
}

func (tickets TicketsService) generateNewTickets(eventIds []int, userId int) {
	for _, eventId := range eventIds {
		tickets.storage.AddTicket(model.NewTicket(eventId, userId))
	}
}

func (tickets TicketsService) RegisterRoutes(server *gin.Engine) {
	server.GET("/tickets/:userId", tickets.getAllUserTickets)
}

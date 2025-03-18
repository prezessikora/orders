package memory

import (
	"github.com/prezessikora/orders/model"
	"log"
)

func NewTicketsDataStore() *TicketsMemoryOrdersStorage {
	log.Println("Using memory.NewTicketsDataStore")
	mos := TicketsMemoryOrdersStorage{tickets: make([]model.Ticket, 0, 10)}
	return &mos
}

type TicketsMemoryOrdersStorage struct {
	tickets []model.Ticket
	nextId  int
}

func (m TicketsMemoryOrdersStorage) GetAll(userId int) []model.Ticket {
	return m.GetUserTickets(userId)
}

func (m *TicketsMemoryOrdersStorage) AddTicket(ticket model.Ticket) {
	m.tickets = append(m.tickets, ticket)
}

func (m TicketsMemoryOrdersStorage) GetUserTickets(userId int) []model.Ticket {
	var result []model.Ticket
	for _, ticket := range m.tickets {
		if ticket.UserId == userId {
			result = append(result, ticket)
		}
	}
	return result
}

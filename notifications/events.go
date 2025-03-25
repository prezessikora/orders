package notifications

import (
	"context"
	"fmt"
	"github.com/prezessikora/orders/service"
	"log"
	"strconv"
)

func HandleEventsDeletions(ctx context.Context, storage service.OrderDataStorage) {
	go func() {
		err := Subscribe("events.event-deletions", updateOrders, ctx, storage)
		if err != nil {
			log.Println("failed to subscribe to event-deletions exchange")
		}
	}()

}

// updateOrders changes the status to all orders for which event was deleted to 'cancelled'
func updateOrders(bytes []byte, storage service.OrderDataStorage) {
	fmt.Printf("cencalling orders with event id: %v\n", string(bytes))
	eventId, _ := strconv.Atoi(string(bytes))
	err := storage.CancelEventOrders(eventId)
	if err != nil {
		fmt.Printf("error cancelling orders: %v\n", err)
		return
	}

}

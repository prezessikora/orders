package order

import (
	"github.com/prezessikora/orders/myclient"
	"testing"
	"time"
)

func Test_SuccessCreateAndVerifyOrder(t *testing.T) {
	testEvent := myclient.Event{
		ID:          1,
		UserID:      1,
		Name:        "Test Event",
		Description: "Test Event Description",
		Location:    "Gandia",
		DateTime:    time.Now().Add(time.Hour * 25),
		Capacity:    10,
		Bookings:    0,
	}
	order, err := Create(1, 1, &testEvent)
	if err != nil {
		t.Fatal(err)
	}
	if !(order.EventId == testEvent.ID) {
		t.Errorf("order has incorrect event id: %v, expected: %v", order.EventId, testEvent.ID)
	}

}

func Test_FailToCreateOrderWithNoCapacity(t *testing.T) {
	testEvent := myclient.Event{
		ID:          1,
		UserID:      1,
		Name:        "Test Event",
		Description: "Test Event Description",
		Location:    "Gandia",
		DateTime:    time.Now().Add(time.Hour * 25),
		Capacity:    10,
		Bookings:    10,
	}
	_, err := Create(1, 1, &testEvent)
	if err != nil && (err.Error()) != "event has no capacity" {
		t.Fatal("verification should fail for full event")
	}

}

func Test_FailLatePurchaseCreateAndVerifyOrder(t *testing.T) {
	testEvent := myclient.Event{
		ID:          1,
		UserID:      1,
		Name:        "Test Event",
		Description: "Test Event Description",
		Location:    "Gandia",
		DateTime:    time.Now().Add(time.Hour * 41),
		Capacity:    10,
		Bookings:    1,
	}
	_, err := Create(1, 1, &testEvent)
	if err == nil {
		t.Fatal("event late booking should not be possible")
	}

	if (err.Error()) != "event start time is less than 24h from now" {
		t.Errorf("expecgted event late booking err but got instead: %v", err)
	}
}

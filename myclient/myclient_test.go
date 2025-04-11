package myclient_test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prezessikora/orders/common"
	"github.com/prezessikora/orders/myclient"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Inter service calls are triggered by client call that should have UUID set by middleware
// Requests missing UUID should be rejected
func TestGetEventFailsWithMissingUUID(t *testing.T) {
	service := myclient.NewEventsServiceClient()
	_, err := service.GetEvent(1, &gin.Context{
		Request:  nil,
		Writer:   nil,
		Params:   nil,
		Keys:     nil,
		Errors:   nil,
		Accepted: nil,
	})
	if err == nil {
		t.Fatal("requests missing UUID should be rejected")
	} else if err.Error() != "missing UUID" {
		t.Errorf("expected missing UUID error, got instead: %v", err)
	}

}

// Inter service calls are triggered by client call that should have UUID on context set by middleware
// Requests with fake UUID should be rejected
// TODO the events should be eventually mocked
func TestGetEventFailsWithFakeRequestId(t *testing.T) {
	service := myclient.NewEventsServiceClient()
	ctx := &gin.Context{}

	ctx.Set(common.XEventsHeaderKey, "12345")
	_, err := service.GetEvent(1, ctx)
	if err == nil {
		t.Fatalf("fake requests id calls should fail but got err instead: %v", err)
	}
}

// TestGetEventSuccessWithValidUUID stubs the server and checks if a http request is sent and Event returned
func TestGetEventSuccessWithValidUUID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a response for the test

		event := myclient.Event{
			ID:          1,
			UserID:      1,
			Name:        "Test Event",
			Description: "Test Event Description",
			Location:    "Gandia",
			DateTime:    time.Now().Add(time.Hour * 25),
			Capacity:    10,
			Bookings:    0,
		}
		eventResponse, err := json.Marshal(event)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("could not marshal json reponse"))
			if err != nil {
				fmt.Printf("error writing stubbed response: %v\n", err)
				return
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(eventResponse))
		if err != nil {
			fmt.Printf("error writing stub response: %v\n", err)
			return
		}
	}))
	defer server.Close()

	service := myclient.NewEventsServiceClient()
	service.Url = server.URL + "/%d"

	ctx := &gin.Context{}
	const eventId = 1
	ctx.Set(common.XEventsHeaderKey, uuid.New())
	event, err := service.GetEvent(eventId, ctx)
	if err != nil {
		t.Fatalf("request should return no error but got err instead: %v", err)
	}
	if event.ID != eventId {
		t.Fatalf("request should return event with requested id: %v but got: %v", eventId, event.ID)
	}
}

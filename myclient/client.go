package myclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prezessikora/orders/common"
	"io"
	"log"
	"net/http"
	"time"
)

// Event represents event from events as it is retrieved from the events service
// TODO import from events service rather than duplicate code
type Event struct {
	ID          int
	UserID      int
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	Capacity    int       `binding:"required"`
	Bookings    int       `binding:"required"`
}

// verify validates the Event struct to ensure all required fields meet the necessary constraints for orders processing.
func (e Event) verify() error {
	if !(e.ID > 0) {
		return fmt.Errorf("can not process event with invalid id: %v", e.ID)
	}
	if !(e.UserID > 0) {
		return fmt.Errorf("can not process event with invalid UserID: %v", e.ID)
	}
	if !(e.Capacity > 0) {
		return errors.New("can not process event with zero capacity")
	}
	return nil
}

// Events service API allowing other services such as orders to access events
type Events struct {
	// TODO abstract local and docker-compose urls
	Url    string
	Client http.Client
}

func NewEventsServiceClient() Events {
	return Events{
		// TODO get this from os env vars set up by docker-compose or container running env
		Url:    "http://localhost:8080/events/%d",
		Client: http.Client{Timeout: time.Duration(1) * time.Second},
	}
}

const XEventsHeaderKey = "X-Events-Request-Id"

func (service Events) GetEvent(eventId int, ctx *gin.Context) (*Event, error) {
	log.Printf("Fetching event with id: [%v]", eventId)
	url := fmt.Sprintf(service.Url, eventId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error creating request %s", err)
		return nil, err
	}

	// add the uuid header to all requests
	requestId, err := common.GetRequestUUIDFromContext(ctx)
	if err != nil {
		return nil, errors.New("missing UUID")
	}
	req.Header.Add(XEventsHeaderKey, requestId)
	// add the uuid header to all requests
	req.Header.Add("Accept", `application/json`)

	resp, err := service.Client.Do(req)

	if err != nil {
		log.Printf("error on event service request %s", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error on closing request body %s", err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)

		var event Event
		err = json.Unmarshal(body, &event)
		if err != nil {
			log.Printf("error parsing json response %s", err)
			return nil, err
		}
		log.Printf("event received:%v", event) // on debug
		err = event.verify()
		if err != nil {
			return nil, err
		}
		return &event, nil

	}
	return nil, fmt.Errorf("could not fetch the event, reponse: %v", resp.StatusCode)
}

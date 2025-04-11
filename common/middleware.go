package common

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

func CorrelationId(c *gin.Context) {

	// this header is used as coorelation id for requests

	header := c.GetHeader(XEventsHeaderKey)
	if header == "" {
		id := uuid.New()
		// further to be set as HTTP header for in between services calls
		c.Set(XEventsHeaderKey, id)
		log.Printf("setting coorelation id [%v] : %v", XEventsHeaderKey, id)
	}
	// before request
	c.Next()
	// after request
}

package common

// Move this to external shared module
import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const XEventsHeaderKey = "X-Events-Request-Id"

func GetRequestUUIDFromContext(ctx *gin.Context) (string, error) {
	requestId, uuidInContext := ctx.Get(XEventsHeaderKey)
	var requestUUID string
	if uuidInContext {
		requestId, valid := requestId.(uuid.UUID)
		if !valid {
			return "", errors.New("context value is not valid UUID")
		}

		requestUUID = requestId.String()
	} else {
		return "", errors.New("UUID is missing in context")
	}
	return requestUUID, nil
}

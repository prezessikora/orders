package common_test

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prezessikora/orders/common"
	"testing"
)

func TestGetRequestIdFromContextFailsWithFakeValue(t *testing.T) {
	ctx := gin.Context{}
	ctx.Set(common.XEventsHeaderKey, "12345")
	_, err := common.GetRequestUUIDFromContext(&ctx)
	if err == nil {
		t.Fatalf("should fail to extract incorrect UUID from context: %v", err)
	}
}

func TestGetRequestIdFromContextSuccess(t *testing.T) {
	ctx := gin.Context{}
	ctx.Set(common.XEventsHeaderKey, uuid.New())
	_, err := common.GetRequestUUIDFromContext(&ctx)
	if err != nil {
		t.Fatalf("failed to extract UUID from context: %v", err)
	}
}

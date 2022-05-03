package message

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/jhandguy/devops-playground/gateway/item"
	"github.com/jhandguy/devops-playground/gateway/object"
	"github.com/jhandguy/devops-playground/gateway/opentelemetry"
	itemPb "github.com/jhandguy/devops-playground/gateway/pb/item"
	objectPb "github.com/jhandguy/devops-playground/gateway/pb/object"
)

type Message struct {
	ID      string `json:"id" uri:"id"`
	Content string `json:"content"`
}

type API struct {
	ItemAPI   *item.API
	ObjectAPI *object.API
}

type result[T any] struct {
	response T
	error    error
}

func (api *API) CheckReadiness(c *gin.Context) {
	ctx := context.Background()
	tracer := opentelemetry.GetTracer("message/CheckReadiness")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "message/CheckReadiness")
		defer span.End()
	}

	objChan := make(chan result[*grpc_health_v1.HealthCheckResponse], 1)
	itmChan := make(chan result[*grpc_health_v1.HealthCheckResponse], 1)

	go func() {
		req := grpc_health_v1.HealthCheckRequest{Service: "liveness"}
		resp, err := api.ObjectAPI.CheckHealth(ctx, &req)
		objChan <- result[*grpc_health_v1.HealthCheckResponse]{resp, err}
	}()

	go func() {
		req := grpc_health_v1.HealthCheckRequest{Service: "liveness"}
		resp, err := api.ItemAPI.CheckHealth(ctx, &req)
		itmChan <- result[*grpc_health_v1.HealthCheckResponse]{resp, err}
	}()

	objRes := <-objChan
	itmRes := <-itmChan

	if err := objRes.error; err != nil {
		zap.S().Errorw("failed to check object readiness", "error", objRes.error.Error())
		c.Status(http.StatusServiceUnavailable)
		return
	}

	if err := itmRes.error; err != nil {
		zap.S().Errorw("failed to check item readiness", "error", itmRes.error.Error())
		c.Status(http.StatusServiceUnavailable)
		return
	}

	if status := objRes.response.GetStatus(); status != grpc_health_v1.HealthCheckResponse_SERVING {
		zap.S().Errorf("failed to check object readiness: %v", grpc_health_v1.HealthCheckResponse_ServingStatus_name[int32(status)])
		c.Status(http.StatusServiceUnavailable)
		return
	}

	if status := itmRes.response.GetStatus(); status != grpc_health_v1.HealthCheckResponse_SERVING {
		zap.S().Errorf("failed to check item readiness: %v", grpc_health_v1.HealthCheckResponse_ServingStatus_name[int32(status)])
		c.Status(http.StatusServiceUnavailable)
		return
	}

	zap.S().Debugw("successfully checked readiness", "traceID", opentelemetry.GetTraceID(ctx))
	c.Status(http.StatusOK)
}

func (api *API) CheckLiveness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (api *API) CreateMessage(c *gin.Context) {
	ctx := context.Background()
	tracer := opentelemetry.GetTracer("message/CreateMessage")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "message/CreateMessage")
		defer span.End()
	}

	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		zap.S().Errorw("failed to decode message", "error", err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}

	objChan := make(chan result[*objectPb.CreateObjectResponse], 1)
	itmChan := make(chan result[*itemPb.CreateItemResponse], 1)

	go func() {
		req := objectPb.CreateObjectRequest{
			Object: &objectPb.Object{
				Id:      msg.ID,
				Content: msg.Content,
			},
		}
		resp, err := api.ObjectAPI.CreateObject(ctx, &req)
		objChan <- result[*objectPb.CreateObjectResponse]{resp, err}
	}()

	go func() {
		req := itemPb.CreateItemRequest{
			Item: &itemPb.Item{
				Id:      msg.ID,
				Content: msg.Content,
			},
		}
		resp, err := api.ItemAPI.CreateItem(ctx, &req)
		itmChan <- result[*itemPb.CreateItemResponse]{resp, err}
	}()

	objRes := <-objChan
	itmRes := <-itmChan

	if err := objRes.error; err != nil {
		zap.S().Errorw("failed to create object", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := itmRes.error; err != nil {
		zap.S().Errorw("failed to create item", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	obj := objRes.response.GetObject()
	itm := itmRes.response.GetItem()

	if obj.GetId() != itm.GetId() || itm.GetContent() != obj.GetContent() {
		zap.S().Errorf("unexpected inconsistencies: %v != %v", obj, itm)
		c.JSON(http.StatusExpectationFailed, gin.H{"error": fmt.Sprintf("unexpected inconsistencies: %v != %v", obj, itm)})
		return
	}

	zap.S().Infow("successfully created message", "msg", msg, "traceID", opentelemetry.GetTraceID(ctx))
	c.JSON(http.StatusOK, msg)
}

func (api *API) GetMessage(c *gin.Context) {
	ctx := context.Background()
	tracer := opentelemetry.GetTracer("message/GetMessage")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "message/GetMessage")
		defer span.End()
	}

	var msg Message
	if err := c.ShouldBindUri(&msg); err != nil {
		zap.S().Errorw("missing id in request", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objChan := make(chan result[*objectPb.GetObjectResponse], 1)
	itmChan := make(chan result[*itemPb.GetItemResponse], 1)

	go func() {
		req := objectPb.GetObjectRequest{
			Id: msg.ID,
		}
		resp, err := api.ObjectAPI.GetObject(ctx, &req)
		objChan <- result[*objectPb.GetObjectResponse]{resp, err}
	}()

	go func() {
		req := itemPb.GetItemRequest{
			Id: msg.ID,
		}
		resp, err := api.ItemAPI.GetItem(ctx, &req)
		itmChan <- result[*itemPb.GetItemResponse]{resp, err}
	}()

	objRes := <-objChan
	itmRes := <-itmChan

	if err := objRes.error; err != nil {
		zap.S().Errorw("failed to get object", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := itmRes.error; err != nil {
		zap.S().Errorw("failed to get item", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	obj := objRes.response.GetObject()
	itm := itmRes.response.GetItem()

	if obj.GetId() != itm.GetId() || obj.GetContent() != itm.GetContent() {
		zap.S().Errorf("unexpected inconsistencies: %v != %v", obj, itm)
		c.JSON(http.StatusExpectationFailed, gin.H{"error": fmt.Sprintf("unexpected inconsistencies: %v != %v", obj, itm)})
		return
	}

	msg.Content = objRes.response.GetObject().GetContent()

	zap.S().Infow("successfully got message", "msg", msg, "traceID", opentelemetry.GetTraceID(ctx))
	c.JSON(http.StatusOK, msg)
}

func (api *API) DeleteMessage(c *gin.Context) {
	ctx := context.Background()
	tracer := opentelemetry.GetTracer("message/DeleteMessage")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "message/DeleteMessage")
		defer span.End()
	}

	var msg Message
	if err := c.ShouldBindUri(&msg); err != nil {
		zap.S().Errorw("missing id in request", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objChan := make(chan result[*objectPb.DeleteObjectResponse], 1)
	itmChan := make(chan result[*itemPb.DeleteItemResponse], 1)

	go func() {
		req := objectPb.DeleteObjectRequest{
			Id: msg.ID,
		}
		resp, err := api.ObjectAPI.DeleteObject(ctx, &req)
		objChan <- result[*objectPb.DeleteObjectResponse]{resp, err}
	}()

	go func() {
		req := itemPb.DeleteItemRequest{
			Id: msg.ID,
		}
		resp, err := api.ItemAPI.DeleteItem(ctx, &req)
		itmChan <- result[*itemPb.DeleteItemResponse]{resp, err}
	}()

	objRes := <-objChan
	itmRes := <-itmChan

	if err := objRes.error; err != nil {
		zap.S().Errorw("failed to delete object", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := itmRes.error; err != nil {
		zap.S().Errorw("failed to delete item", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	zap.S().Infow("successfully deleted message", "msg", msg, "traceID", opentelemetry.GetTraceID(ctx))
	c.Status(http.StatusOK)
}

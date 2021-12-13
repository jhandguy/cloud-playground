package message

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	ID      string `json:"id"`
	Content string `json:"content"`
}

type API struct {
	ItemAPI   *item.API
	ObjectAPI *object.API
}

func (api *API) CheckReadiness(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	tracer := opentelemetry.GetTracer("message/CheckReadiness")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "message/CheckReadiness")
		defer span.End()
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var resps []*grpc_health_v1.HealthCheckResponse
	var errs []error
	go func() {
		defer wg.Done()
		req := grpc_health_v1.HealthCheckRequest{Service: "liveness"}
		resp, err := api.ObjectAPI.CheckHealth(ctx, &req)
		resps = append(resps, resp)
		errs = append(errs, err)
	}()

	go func() {
		defer wg.Done()
		req := grpc_health_v1.HealthCheckRequest{Service: "liveness"}
		resp, err := api.ItemAPI.CheckHealth(ctx, &req)
		resps = append(resps, resp)
		errs = append(errs, err)
	}()

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			zap.S().Errorw("failed to check readiness", "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}

	for _, resp := range resps {
		if status := resp.GetStatus(); status != grpc_health_v1.HealthCheckResponse_SERVING {
			zap.S().Errorf("failed to check readiness: %v", grpc_health_v1.HealthCheckResponse_ServingStatus_name[int32(status)])
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}

	zap.S().Debugw("successfully checked readiness", "traceID", opentelemetry.GetTraceID(ctx))
}

func (api *API) CheckLiveness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (api *API) CreateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tracer := opentelemetry.GetTracer("message/CreateMessage")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "message/CreateMessage")
		defer span.End()
	}

	var msg Message

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		zap.S().Errorw("failed to decode message", "error", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var objResp *objectPb.CreateObjectResponse
	var objErr error
	go func() {
		defer wg.Done()
		req := objectPb.CreateObjectRequest{
			Object: &objectPb.Object{
				Id:      msg.ID,
				Content: msg.Content,
			},
		}
		objResp, objErr = api.ObjectAPI.CreateObject(ctx, &req)
	}()

	var itmResp *itemPb.CreateItemResponse
	var itmErr error
	go func() {
		defer wg.Done()
		req := itemPb.CreateItemRequest{
			Item: &itemPb.Item{
				Id:      msg.ID,
				Content: msg.Content,
			},
		}
		itmResp, itmErr = api.ItemAPI.CreateItem(ctx, &req)
	}()

	wg.Wait()

	if objErr != nil {
		zap.S().Errorw("failed to create object", "error", objErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if itmErr != nil {
		zap.S().Errorw("failed to create item", "error", itmErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if objResp.GetObject().GetId() != itmResp.GetItem().GetId() ||
		objResp.GetObject().GetContent() != itmResp.GetItem().GetContent() {
		zap.S().Errorf("unexpected inconsistencies: %v != %v", objResp, itmResp)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	err = json.NewEncoder(w).Encode(msg)
	if err != nil {
		zap.S().Errorw("failed to encode message", "error", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	zap.S().Infow("successfully created message", "msg", msg, "traceID", opentelemetry.GetTraceID(ctx))
}

func (api *API) GetMessage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tracer := opentelemetry.GetTracer("message/GetMessage")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "message/GetMessage")
		defer span.End()
	}

	msg := Message{
		ID: mux.Vars(r)["id"],
	}

	if msg.ID == "" {
		zap.S().Error("missing id in request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var objResp *objectPb.GetObjectResponse
	var objErr error
	go func() {
		defer wg.Done()
		req := objectPb.GetObjectRequest{
			Id: msg.ID,
		}
		objResp, objErr = api.ObjectAPI.GetObject(ctx, &req)
	}()

	var itmResp *itemPb.GetItemResponse
	var itmErr error
	go func() {
		defer wg.Done()
		req := itemPb.GetItemRequest{
			Id: msg.ID,
		}
		itmResp, itmErr = api.ItemAPI.GetItem(ctx, &req)
	}()

	wg.Wait()

	if objErr != nil {
		zap.S().Errorw("failed to get object", "error", objErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if itmErr != nil {
		zap.S().Errorw("failed to get item", "error", itmErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if objResp.GetObject().GetId() != itmResp.GetItem().GetId() ||
		objResp.GetObject().GetContent() != itmResp.GetItem().GetContent() {
		zap.S().Errorf("unexpected inconsistencies: %v != %v", objResp, itmResp)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	msg.Content = objResp.GetObject().GetContent()

	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		zap.S().Errorw("failed to encode message", "error", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	zap.S().Infow("successfully got message", "msg", msg, "traceID", opentelemetry.GetTraceID(ctx))
}

func (api *API) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tracer := opentelemetry.GetTracer("message/DeleteMessage")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "message/DeleteMessage")
		defer span.End()
	}

	msg := Message{
		ID: mux.Vars(r)["id"],
	}

	if msg.ID == "" {
		zap.S().Error("missing id in request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var objErr error
	go func() {
		defer wg.Done()
		req := objectPb.DeleteObjectRequest{
			Id: msg.ID,
		}
		_, objErr = api.ObjectAPI.DeleteObject(ctx, &req)
	}()

	var itmErr error
	go func() {
		defer wg.Done()
		req := itemPb.DeleteItemRequest{
			Id: msg.ID,
		}
		_, itmErr = api.ItemAPI.DeleteItem(ctx, &req)
	}()

	wg.Wait()

	if objErr != nil {
		zap.S().Errorw("failed to delete object", "error", objErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if itmErr != nil {
		zap.S().Errorw("failed to delete item", "error", itmErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	zap.S().Infow("successfully deleted message", "msg", msg, "traceID", opentelemetry.GetTraceID(ctx))
}

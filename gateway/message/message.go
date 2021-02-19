package message

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/jhandguy/devops-playground/gateway/item"
	"github.com/jhandguy/devops-playground/gateway/object"
	itemPb "github.com/jhandguy/devops-playground/gateway/pb/item"
	objectPb "github.com/jhandguy/devops-playground/gateway/pb/object"
)

type Message struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type API struct {
	ItemAPI   *item.API
	ObjectAPI *object.API
}

func (api *API) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Printf("failed to decode message: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var objResp *objectPb.CreateObjectResponse
	var objErr error
	go func() {
		defer wg.Done()
		req := objectPb.CreateObjectRequest{
			Name:    msg.Name,
			Content: msg.Content,
		}
		objResp, objErr = api.ObjectAPI.CreateObject(&req)
	}()

	var itmResp *itemPb.CreateItemResponse
	var itmErr error
	go func() {
		defer wg.Done()
		req := itemPb.CreateItemRequest{
			Name:    msg.Name,
			Content: msg.Content,
		}
		itmResp, itmErr = api.ItemAPI.CreateItem(&req)
	}()

	wg.Wait()

	if objErr != nil || itmErr != nil {
		log.Printf("failed to create message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if objResp.GetObject().GetName() != itmResp.GetItem().GetName() ||
		objResp.GetObject().GetContent() != itmResp.GetItem().GetContent() {
		log.Printf("unexpected inconsistencies: %v != %v", objResp, itmResp)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	msg.ID = itmResp.GetItem().GetId()

	err = json.NewEncoder(w).Encode(msg)
	if err != nil {
		log.Printf("failed to encode message: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
}

func (api *API) GetMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Printf("failed to decode message: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var objResp *objectPb.GetObjectResponse
	var objErr error
	go func() {
		defer wg.Done()
		req := objectPb.GetObjectRequest{
			Name: msg.Name,
		}
		objResp, objErr = api.ObjectAPI.GetObject(&req)
	}()

	var itmResp *itemPb.GetItemResponse
	var itmErr error
	go func() {
		defer wg.Done()
		req := itemPb.GetItemRequest{
			Id: msg.ID,
		}
		itmResp, itmErr = api.ItemAPI.GetItem(&req)
	}()

	wg.Wait()

	if objErr != nil || itmErr != nil {
		log.Printf("failed to get message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if objResp.GetObject().GetName() != itmResp.GetItem().GetName() ||
		objResp.GetObject().GetContent() != itmResp.GetItem().GetContent() {
		log.Printf("unexpected inconsistencies: %v != %v", objResp, itmResp)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	msg.Content = objResp.GetObject().GetContent()

	err = json.NewEncoder(w).Encode(msg)
	if err != nil {
		log.Printf("failed to encode message: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
}

func (api *API) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Printf("failed to decode message: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var objErr error
	go func() {
		defer wg.Done()
		req := objectPb.DeleteObjectRequest{
			Name: msg.Name,
		}
		_, objErr = api.ObjectAPI.DeleteObject(&req)
	}()

	var itmErr error
	go func() {
		defer wg.Done()
		req := itemPb.DeleteItemRequest{
			Id: msg.ID,
		}
		_, itmErr = api.ItemAPI.DeleteItem(&req)
	}()

	wg.Wait()

	if objErr != nil || itmErr != nil {
		log.Printf("failed to delete message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

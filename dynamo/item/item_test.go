package item

import (
	"bytes"
	"dynamo/test"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestCreateItem(t *testing.T) {
	var actTable string
	var actItem Item

	fun := func(table string, item Item) error {
		actTable = table
		actItem = item
		return nil
	}

	table := "table"
	sut := CreateItem(fun, table)

	item := Item{
		Id:      "id",
		Name:    "item",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(item)

	code, _ := test.RecordRequest(sut, http.MethodPost, "", bytes.NewReader(byt))

	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, actTable, table)
	test.AssertEqual(t, actItem, item)
}

func TestGetItem(t *testing.T) {
	var actTable string
	var actId string

	item := Item{
		Id:      "id",
		Name:    "item",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(item)

	fun := func(table string, id string) (map[string]*dynamodb.AttributeValue, error) {
		actTable = table
		actId = id

		return dynamodbattribute.MarshalMap(item)
	}

	table := "table"
	sut := GetItem(fun, table)

	code, body := test.RecordRequest(sut, http.MethodGet, item.Id, bytes.NewReader(byt))

	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, actTable, table)
	test.AssertEqual(t, actId, item.Id)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))
}

func TestDeleteItem(t *testing.T) {
	var actTable string
	var actId string

	fun := func(table string, id string) error {
		actTable = table
		actId = id

		return nil
	}

	table := "table"
	id := "id"
	sut := DeleteItem(fun, table)

	code, _ := test.RecordRequest(sut, http.MethodDelete, id, nil)

	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, actTable, table)
	test.AssertEqual(t, actId, id)
}

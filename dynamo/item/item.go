package item

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func CreateItem(createItem func(string, Item) error, table string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		item := Item{
			Id: uuid.New().String(),
		}

		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		err = createItem(table, item)
		if err != nil {
			log.Printf("failed to insert item %s (%s) in table %s\nwith content: %s\nerror: %v\n", item.Name, item.Id, table, item.Content, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(item)
		if err != nil {
			log.Printf("failed to encode item: %v\nerror: %v\n", item, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("successfully inserted item %s (%s) in table %s\nwith content: %s\n", item.Name, item.Id, table, item.Content)
	}
}

func GetItem(getItem func(string, string) (map[string]*dynamodb.AttributeValue, error), table string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getIdQuery(r)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		value, err := getItem(table, id)
		if err != nil {
			log.Printf("failed to select item %s from table %s\nerror: %v\n", id, table, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if len(value) == 0 {
			log.Printf("failed to find item %s from table %s\n", id, table)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var item Item
		err = dynamodbattribute.UnmarshalMap(value, &item)
		if err != nil {
			log.Printf("failed to unmarshal item: %v\nerror: %v\n", item, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(item)
		if err != nil {
			log.Printf("failed to encode item: %v\nerror: %v\n", item, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("successfully selected item %s (%s) from table %s\nwith content: %s\n", item.Name, item.Id, table, item.Content)
	}
}

func DeleteItem(deleteItem func(string, string) error, table string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getIdQuery(r)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := deleteItem(table, id)
		if err != nil {
			log.Printf("failed to delete item %s from table %s\nerror: %v\n", id, table, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("successfully deleted item %s from table %s\n", id, table)
	}
}

func CreateItemFunc(client *dynamodb.DynamoDB) func(string, Item) error {
	return func(table string, item Item) error {
		it, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			return err
		}

		_, err = client.PutItemWithContext(context.Background(), &dynamodb.PutItemInput{
			Item:      it,
			TableName: aws.String(table),
		})

		return err
	}
}

func GetItemFunc(client *dynamodb.DynamoDB) func(string, string) (map[string]*dynamodb.AttributeValue, error) {
	return func(table string, id string) (map[string]*dynamodb.AttributeValue, error) {
		out, err := client.GetItemWithContext(context.Background(), &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"id": {
					S: aws.String(id),
				},
			},
			TableName: aws.String(table),
		})
		if err != nil {
			return nil, err
		}

		return out.Item, nil
	}
}

func DeleteItemFunc(client *dynamodb.DynamoDB) func(string, string) error {
	return func(table string, id string) error {
		_, err := client.DeleteItemWithContext(context.Background(), &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"id": {
					S: aws.String(id),
				},
			},
			TableName: aws.String(table),
		})

		return err
	}
}

func getIdQuery(r *http.Request) (string, bool) {
	query, ok := r.URL.Query()["id"]

	if !ok || len(query[0]) < 1 {
		return "", false
	}

	return query[0], true
}

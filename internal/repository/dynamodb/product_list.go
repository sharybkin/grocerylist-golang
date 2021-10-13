package dynamo_repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/pkg/db"
	log "github.com/sirupsen/logrus"
)

type ProductList struct {
	database *db.DynamoDB
}

const (
	productListsTable = "product_lists"
)

func NewProductList(db *db.DynamoDB) *ProductList {
	return &ProductList{database: db}
}

func (p *ProductList) GetAllProductLists(userId string) ([]model.ProductList, error) {

	var lists []model.ProductList

	client, err := p.database.GetClient()

	if err != nil {
		return lists, err
	}

	id := "1df81a12-5841-470e-afa3-e9cdc10d04f8"

	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(productListsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return lists, err
	}
	var productList model.ProductList

	err = attributevalue.UnmarshalMap(out.Item, &productList)
	if err != nil {
		log.Fatalln(err.Error())
		return lists, fmt.Errorf("failed to unmarshal Items, %w", err)
	}

	lists = append(lists, productList)

	return lists, nil
}

func (p *ProductList) AddProductList(list model.ProductList) (string, error) {

	list.Id = uuid.New().String()

	client, err := p.database.GetClient()

	if err != nil {
		return "", err
	}

	av, err := attributevalue.MarshalMap(list)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Record, %w", err)
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(productListsTable),
		Item:      av,
	})

	if err != nil {
		return "", fmt.Errorf("failed to put Record, %w", err)
	}


	log.WithFields(log.Fields{
		"listName": list.Name,
	}).Info("ProductList was added")

	return list.Id, nil
}

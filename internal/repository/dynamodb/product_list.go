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

//func (p *ProductList) GetAllProductLists(userId string) ([]model.ProductList, error) {
//
//	var lists []model.ProductList
//
//	client, err := p.database.GetClient()
//
//	if err != nil {
//		return lists, err
//	}
//
//	id := "1df81a12-5841-470e-afa3-e9cdc10d04f8"
//
//	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
//		TableName: aws.String(productListsTable),
//		Key: map[string]types.AttributeValue{
//			"id": &types.AttributeValueMemberS{Value: id},
//		},
//	})
//
//	if err != nil {
//		return lists, err
//	}
//	var productList model.ProductList
//
//	err = attributevalue.UnmarshalMap(out.Item, &productList)
//	if err != nil {
//		log.Fatalln(err.Error())
//		return lists, fmt.Errorf("failed to unmarshal Items, %w", err)
//	}
//
//	lists = append(lists, productList)
//
//	return lists, nil
//}

func (p *ProductList) GetProductList(listId string) (model.ProductList, error) {

	var list model.ProductList

	client, err := p.database.GetClient()

	if err != nil {
		return list, err
	}

	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(productListsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: listId},
		},
	})

	if err != nil {
		return list, err
	}

	err = attributevalue.UnmarshalMap(out.Item, &list)
	if err != nil {
		log.Fatalln(err.Error())
		return list, fmt.Errorf("failed to unmarshal Items, %w", err)
	}

	return list, nil
}

func (p *ProductList) CreateProductList(request model.ProductListRequest) (string, error) {

	list := model.ProductList{
		Id:   uuid.New().String(),
		Name: request.Name}

	err := p.createOrUpdateProductList(list)

	if err != nil {
		return "", fmt.Errorf("failed to put Record, %w", err)
	}

	log.WithFields(log.Fields{
		"listName": request.Name,
	}).Info("ProductList was added")

	return list.Id, nil
}


func (p *ProductList) UpdateProductList(list model.ProductList) error {
	err := p.createOrUpdateProductList(list)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"listName": list.Name,
	}).Info("ProductList was updated")

	return nil
}

func (p *ProductList) createOrUpdateProductList(list model.ProductList) error {
	client, err := p.database.GetClient()

	if err != nil {
		return err
	}

	av, err := attributevalue.MarshalMap(list)
	if err != nil {
		return fmt.Errorf("failed to marshal Record, %w", err)
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(productListsTable),
		Item:      av,
	})

	return nil
}

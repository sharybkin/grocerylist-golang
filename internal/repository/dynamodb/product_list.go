package dynamo_repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	ext "github.com/sharybkin/grocerylist-golang/internal/extension"
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
		Id:       uuid.New().String(),
		Name:     request.Name,
		Products: map[string]model.Product{},
	}

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
		"listName": request.Name,
	}).Debug("ProductList was added")

	return list.Id, nil
}

func (p *ProductList) UpdateProductList(listId string, request model.ProductListRequest) error {
	client, err := p.database.GetClient()

	if err != nil {
		return err
	}

	_, err = client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(productListsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: listId},
		},
		ExpressionAttributeNames: map[string]string{
			"#name": "name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":n": &types.AttributeValueMemberS{Value: request.Name},
		},
		UpdateExpression: aws.String("set #name = :n"),
		ReturnValues:     types.ReturnValueUpdatedNew,
	})

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"listId":  listId,
		"newName": request.Name,
	}).Debug("ProductList was updated")

	return nil
}

func (p *ProductList) DeleteProductList(listId string) error {
	client, err := p.database.GetClient()

	if err != nil {
		return err
	}

	_, err = client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(productListsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: listId},
		},
	})

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"listId": listId,
	}).Debug("ProductList was deleted")

	return nil
}

func (p *ProductList) GetProducts(listId string) ([]model.Product, error) {
	productList, err := p.GetProductList(listId)

	if err != nil {
		return nil, err
	}

	return ext.GetValues(productList.Products), nil
}

func (p *ProductList) AddOrUpdateProduct(listId string, product model.Product, update bool) (string, error) {
	client, err := p.database.GetClient()
	if err != nil {
		return "", err
	}

	if !update {
		product.Id = uuid.New().String()
	}

	av, err := attributevalue.MarshalMap(product)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Record, %w", err)
	}

	_, err = client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(productListsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: listId},
		},
		ExpressionAttributeNames: map[string]string{
			"#id": product.Id,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":product": &types.AttributeValueMemberM{Value: av},
		},
		UpdateExpression: aws.String("SET products.#id = :product"),
		//ConditionExpression: aws.String("attribute_not_exists(products.#id)"),
		ReturnValues: types.ReturnValueAllNew,
	})

	if err != nil {
		return "", err
	}

	return product.Id, nil
}

func (p *ProductList) DeleteProduct(listId string, productId string) error {
	client, err := p.database.GetClient()
	if err != nil {
		return err
	}

	_, err = client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(productListsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: listId},
		},
		ExpressionAttributeNames: map[string]string{
			"#id": productId,
		},
		UpdateExpression: aws.String("REMOVE products.#id"),
		ReturnValues:     types.ReturnValueNone,
	})

	return err
}

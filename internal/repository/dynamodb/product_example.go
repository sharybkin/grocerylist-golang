package dynamo_repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/pkg/db"
	log "github.com/sirupsen/logrus"
)

type ProductExample struct {
	database *db.DynamoDB
}

const (
	productExampleTable = "product_examples"
)

func NewProductExample(db *db.DynamoDB) *ProductExample {
	return &ProductExample{database: db}
}

func (p *ProductExample) GetExamples() ([]model.ProductExample, error) {
	client, err := p.database.GetClient()
	if err != nil {
		return nil, err
	}

	result, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(productExampleTable),
	})

	if err != nil {
		return nil, err
	}

	productExamples := make([]model.ProductExample, 0, len(result.Items))

	for _, i := range result.Items {
		var example model.ProductExample

		err = attributevalue.UnmarshalMap(i, &example)

		if err != nil {
			log.Fatalln(err.Error())
			return nil, fmt.Errorf("failed to unmarshal Items, %w", err)
		}

		productExamples = append(productExamples, example)
	}

	return productExamples, nil
}

func (p *ProductExample) AddOrUpdate(example model.ProductExample) error {
	client, err := p.database.GetClient()
	if err != nil {
		return err
	}

	av, err := attributevalue.MarshalMap(example)
	if err != nil {
		return fmt.Errorf("failed to marshal Record, %w", err)
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(productExampleTable),
		Item:      av,
	})

	return err
}

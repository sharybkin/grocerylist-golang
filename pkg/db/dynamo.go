package db

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/viper"
)

type DynamoDB struct {
	config aws.Config
	client *dynamodb.Client
}

func NewDynamoDB() (*DynamoDB, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	return &DynamoDB{config: cfg}, nil
}

func (d *DynamoDB) GetClient() (*dynamodb.Client, error) {

	if d.client == nil {
		dbClient, err := createClient(d.config)
		if err != nil {
			return nil, fmt.Errorf("failed to create dynamodb client, %w", err)
		}
		d.client =  dbClient
		log.Infoln("dynamodb client was created")
	}
	return d.client, nil
}

func createClient(cfg aws.Config) (*dynamodb.Client, error) {

	client := dynamodb.NewFromConfig(cfg)
	return client, nil
}

func getConfig() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		aws := viper.GetStringMapString("aws")

		o.Region = aws["region"]
		return nil
	})

	if err != nil {
		log.Fatalln("Error while loading DynamoDB config")
	}

	return cfg, err
}

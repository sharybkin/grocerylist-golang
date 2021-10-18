package dynamo_repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/pkg/db"
	log "github.com/sirupsen/logrus"
)

type UserList struct {
	database *db.DynamoDB
}

const (
	userListTable = "user_lists"
)

func NewUserList(db *db.DynamoDB) *UserList {
	return &UserList{database: db}
}

func (u *UserList) GetUserLists(userId string) ([]model.UserProductListInfo, error) {

	//TODO: delete info log
	log.WithFields(log.Fields{
		"userId": userId,
	}).Info("GetUserLists")

	listInfo := make([]model.UserProductListInfo, 0)

	client, err := u.database.GetClient()
	if err != nil {
		return listInfo, err
	}

	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(userListTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userId},
		},
	})

	if err != nil {
		return listInfo, err
	}

	var userLists model.UserLists

	err = attributevalue.UnmarshalMap(out.Item, &userLists)
	if err != nil {
		log.Fatalln(err.Error())
		return listInfo, fmt.Errorf("failed to unmarshal Items, %w", err)
	}

	if userLists.UserId == "" {
		return listInfo, nil
	}



	return userLists.ProductLists, nil
}

func (u *UserList) LinkListToUser(userId string, listInfo model.UserProductListInfo) error {

	client, err := u.database.GetClient()

	if err != nil {
		return err
	}

	productLists, err := u.GetUserLists(userId)

	if err != nil{
		return fmt.Errorf("getting user information failed, %w", err)
	}

	productLists = append(productLists, listInfo)

	userList := model.UserLists{UserId: userId, ProductLists: productLists}

	av, err := attributevalue.MarshalMap(userList)
	if err != nil {
		return fmt.Errorf("failed to marshal Record, %w", err)
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(userListTable),
		Item:      av,
	})

	if err != nil {
		return fmt.Errorf("failed to put Record, %w", err)
	}

	log.WithFields(log.Fields{
		"userId": userId,
		"listName": listInfo.Name,
	}).Info("List was linked")

	return nil
}
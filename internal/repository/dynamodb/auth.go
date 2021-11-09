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
	"github.com/sharybkin/grocerylist-golang/pkg/extension"
	log "github.com/sirupsen/logrus"
)

const (
	userTable string = "users"
)

type Auth struct {
	database *db.DynamoDB
}

func NewAuth(db *db.DynamoDB) *Auth {
	return &Auth{database: db}
}

func (s *Auth) CreateUser(user model.User) (string, error) {

	exists, err := s.checkUserExists(user.Username)
	if err != nil {
		return "", err
	}

	if exists {
		message := fmt.Sprintf("'%s' already exists", user.Username)
		log.Infoln(message)
		return "", &extension.AuthError{HttpError: extension.HttpError{Message: message}}
	}

	user.Id = uuid.New().String()

	log.WithFields(log.Fields{
		"user": user.Name,
	}).Info("User was added")

	client, err := s.database.GetClient()
	if err != nil {
		return "", err
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(userTable),
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: user.Id},
			"username": &types.AttributeValueMemberS{Value: user.Username},
			"name":     &types.AttributeValueMemberS{Value: user.Name},
			"password": &types.AttributeValueMemberS{Value: user.Password},
		},
	})

	if err != nil {
		return "", err
	}

	return user.Id, nil
}

func (s *Auth) checkUserExists(username string) (bool, error) {
	user, err := s.getUser(username)

	if err != nil || user.Id == "" {
		return false, err
	}

	return true, nil
}

func (s *Auth) getUser(username string) (model.User, error) {

	client, err := s.database.GetClient()

	if err != nil {
		return model.User{}, err
	}

	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(userTable),
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
		},
	})

	if err != nil {
		return model.User{}, err
	}

	var userModel model.User

	err = attributevalue.UnmarshalMap(out.Item, &userModel)
	if err != nil {
		log.Fatalln(err.Error())
		return model.User{}, fmt.Errorf("failed to unmarshal Items, %w", err)
	}

	return userModel, nil
}

func (s *Auth) GetUserByCredentials(username, password string) (model.User, error) {

	userModel, err := s.getUser(username)

	if err != nil {
		return model.User{}, err
	}

	if userModel.Id == "" {

		message := fmt.Sprintf("user %s not found", username)
		log.Infoln(message)
		return model.User{}, &extension.AuthError{HttpError: extension.HttpError{Message: message}}
	}

	if userModel.Password != password {

		message := fmt.Sprintf("Invalid password for user '%s'", username)
		log.Infoln(message)

		return model.User{}, &extension.AuthError{HttpError: extension.HttpError{Message: "invalid password"}}
	}

	return userModel, nil
}

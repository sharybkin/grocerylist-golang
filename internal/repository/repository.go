package repository

import (
	"github.com/sharybkin/grocerylist-golang/internal/model"
	dynamorepo "github.com/sharybkin/grocerylist-golang/internal/repository/dynamodb"

	"github.com/sharybkin/grocerylist-golang/pkg/db"

	log "github.com/sirupsen/logrus"
)

type Authorization interface {
	CreateUser(user model.User) (string, error)
	GetUserByCredentials(username, password string) (model.User, error)
}

type ProductList interface {
	GetProductList(listId string) (model.ProductList, error)
	CreateProductList(request model.ProductListRequest) (string, error)
	UpdateProductList(list model.ProductList) error
	DeleteProductList(userId string, listId string) error
	GetProducts(listId string) ([]model.Product, error)
}

type UserList interface {
	GetUserLists(userId string) ([]model.UserProductListInfo, error)
	LinkListToUser(userId string, listInfo model.UserProductListInfo) error
	UpdateUserList(userId string, listInfo model.UserProductListInfo) error
	UnlinkListFromUser(userId string, listId string) error
}

type Repository struct {
	Authorization
	ProductList
	UserList
}

func NewRepository() *Repository {

	db, err := db.NewDynamoDB()
	if err != nil {
		log.Fatalln("DynamoDB initialization failed", err.Error())
	}

	return &Repository{
		Authorization: dynamorepo.NewAuth(db),
		ProductList:   dynamorepo.NewProductList(db),
		UserList:      dynamorepo.NewUserList(db),
	}
}

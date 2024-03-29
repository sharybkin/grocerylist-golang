package service

import (
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
)

type Authorization interface {
	CreateUser(user model.User) (string, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (string, error)
}

type ProductList interface {
	CreateProductList(userId string, request model.ProductListRequest) (string, error)
	UpdateProductList(userId string, listId string, request model.ProductListRequest) error
	DeleteProductList(userId string, listId string) error
}

type UserLists interface {
	GetUserLists(userId string) ([]model.UserProductListInfo, error)
	LinkListToUser(userId string, list model.UserProductListInfo) error
	UpdateUserList(userId string, listInfo model.UserProductListInfo) error
}

type Product interface {
	GetAllProducts(userId string, listId string) ([]model.Product, error)
	AddProduct(userId string, listId string, product model.Product) (string, error)
	UpdateProduct(userId string, listId string, product model.Product) error
	DeleteProduct(userId string, listId string, productId string) error
}

type ProductExample interface {
	GetProductExamples() []string
	UpdateUsageStatistic(name string)
}

type ServicesHolder struct {
	Authorization
	ProductList
	UserLists
	Product
	ProductExample
}

func NewService(repository *repository.Repository) *ServicesHolder {

	userListService := NewUserListService(repository.UserList)
	productExampleService := NewProductExampleService(repository.ProductExample)

	return &ServicesHolder{
		Authorization:  NewAuthService(repository.Authorization),
		ProductList:    NewProductListService(repository.ProductList, userListService),
		UserLists:      userListService,
		Product:        NewProductService(repository.ProductList, userListService, productExampleService),
		ProductExample: productExampleService,
	}
}

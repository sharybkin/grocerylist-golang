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
	CreateProductList(userId string, request model.ProductListRequest) (string, int, error)
	UpdateProductList(userId string, listId string, request model.ProductListRequest) (int, error)
	DeleteProductList(userId string, listId string) (int, error)
}

type UserLists interface {
	GetUserLists(userId string) ([]model.UserProductListInfo, error)
	LinkListToUser(userId string, list model.UserProductListInfo) error
	UpdateUserList(userId string, listInfo model.UserProductListInfo) error
}

type ServicesHolder struct {
	Authorization
	ProductList
	UserLists
}

func NewService(repository *repository.Repository) *ServicesHolder {

	userListService := NewUserListService(repository.UserList)

	return &ServicesHolder{
		Authorization: NewAuthService(repository.Authorization),
		ProductList:   NewProductListService(repository.ProductList, userListService),
		UserLists: userListService,
	}
}

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
	GetAllProductLists(userId string) ([]model.ProductList, error)
}

type UserLists interface {
	GetUserLists(userId string) ([]model.UserProductListInfo, error)
	LinkListToUser(listId string, userId string) error
}

type ServicesHolder struct {
	Authorization
	ProductList
	UserLists
}

func NewService(repository *repository.Repository) *ServicesHolder {
	return &ServicesHolder{
		Authorization: NewAuthService(repository.Authorization),
		ProductList:   NewProductListService(repository.ProductList),
		UserLists: NewUserListService(repository.UserList),
	}
}

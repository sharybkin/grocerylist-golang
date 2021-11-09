package service

import (
	"github.com/sharybkin/grocerylist-golang/internal/model"
)

type ProductService struct {
}

func (p ProductService) GetAllProducts(userId string, listId string) ([]model.Product, int, error) {
	panic("implement me")
}

func NewProductService() *ProductService {
	return &ProductService{}
}

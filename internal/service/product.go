package service

import (
	"errors"
	"fmt"
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
)

type ProductService struct {
	repo            repository.ProductList
	userListService *UserListService
}

func NewProductService(repo repository.ProductList, service *UserListService) *ProductService {
	return &ProductService{repo: repo, userListService: service}
}

func (p *ProductService) GetAllProducts(userId string, listId string) ([]model.Product, error) {
	if p.repo == nil {
		return nil, errors.New("null pointer exception")
	}

	if err := p.userListService.checkUserList(userId, listId); err != nil {
		return nil, err
	}

	products, err := p.repo.GetProducts(listId)

	if err != nil {
		return nil, fmt.Errorf("failed to load products from list [%s] %w", listId, err)
	}

	return products, nil
}

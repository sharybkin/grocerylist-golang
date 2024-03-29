package service

import (
	"errors"
	"fmt"
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
	log "github.com/sirupsen/logrus"
)

type ProductService struct {
	repo                  repository.ProductList
	userListService       *UserListService
	productExampleService *ProductExampleService
}

func NewProductService(repo repository.ProductList, listService *UserListService, exampleService *ProductExampleService) *ProductService {
	return &ProductService{
		repo:                  repo,
		userListService:       listService,
		productExampleService: exampleService,
	}
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

func (p *ProductService) AddProduct(userId string, listId string, product model.Product) (string, error) {
	if p.repo == nil {
		return "", errors.New("null pointer exception")
	}

	if err := p.userListService.checkUserList(userId, listId); err != nil {
		return "", err
	}

	productId, err := p.repo.AddOrUpdateProduct(listId, product, false)

	if err != nil {
		return "", fmt.Errorf("failed to add product to list [%s] %w", listId, err)
	}

	p.productExampleService.UpdateUsageStatistic(product.Name)

	log.WithFields(log.Fields{
		"listId":  listId,
		"product": product.Name,
	}).Debug("Product was added")

	return productId, err
}

func (p *ProductService) UpdateProduct(userId string, listId string, product model.Product) error {
	if p.repo == nil {
		return errors.New("null pointer exception")
	}

	if err := p.userListService.checkUserList(userId, listId); err != nil {
		return err
	}

	_, err := p.repo.AddOrUpdateProduct(listId, product, true)

	if err != nil {
		return fmt.Errorf("failed to update product from list [%s] %w", listId, err)
	}

	return nil
}

func (p *ProductService) DeleteProduct(userId string, listId string, productId string) error {
	if p.repo == nil {
		return errors.New("null pointer exception")
	}

	if err := p.userListService.checkUserList(userId, listId); err != nil {
		return err
	}

	if err := p.repo.DeleteProduct(listId, productId); err != nil {
		return fmt.Errorf("failed to delete product from list [%s] %w", listId, err)
	}

	return nil
}

package service

import (
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
)

type ProductListService struct {
	repo repository.ProductList
}

func NewProductListService(repo repository.ProductList) *ProductListService {
	return &ProductListService{repo: repo}
}

func (p *ProductListService) GetAllProductLists(userId string) ([]model.ProductList, error) {
	return p.repo.GetAllProductLists(userId)
}

func (p *ProductListService) AddProductList(list model.ProductList) (string, error) {
	pr1 := model.Product{Name: "молоко", Count: 1.5, IsDone: false}
	pr2 := model.Product{Name: "шашлык", Count: 1.6, IsDone: false}

	list1 := model.ProductList{Products: []model.Product{pr1, pr2}, Name: "Test"}

	return p.repo.AddProductList(list1)
}

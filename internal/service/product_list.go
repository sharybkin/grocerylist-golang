package service

import (
	"fmt"
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
	"strings"
)

type ProductListService struct {
	repo repository.ProductList
	userListService *UserListService
}

func NewProductListService(repo repository.ProductList, userListService *UserListService) *ProductListService {
	return &ProductListService{repo: repo, userListService: userListService}
}

func (p *ProductListService) GetProductList(listId string) (model.ProductList, error) {
	return p.repo.GetProductList(listId)
}

func (p *ProductListService) CreateProductList(userId string,list model.ProductList) (string, int, error) {

	list.Name = strings.Trim(list.Name, " ")

	listId, err := p.repo.CreateProductList(list)
	if err != nil {
		return "", 500, err
	}

	lists, err := p.userListService.GetUserLists(userId)
	if err != nil {
		return "", 500, fmt.Errorf("cannot get product lists for user [%s], %w", userId, err)
	}

	for _, productList := range lists {
		if strings.EqualFold(productList.Name, list.Name) {
			return "", 400, fmt.Errorf("[%s] already exists", list.Name)
		}
	}

	err = p.userListService.LinkListToUser(userId, model.UserProductListInfo{Id: listId, Name: list.Name})

	if err != nil {
		return "", 500, fmt.Errorf("cannot link list [%s] to user [%s], %w", listId, userId, err)
	}

	return listId, 200, nil
}

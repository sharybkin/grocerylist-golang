package service

import (
	"errors"
	"fmt"
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
	"strings"
)

type ProductListService struct {
	repo repository.ProductList
	userListService *UserListService
}

func (p *ProductListService) DeleteProductList(userId string, listId string) (int, error) {

	code, err := p.checkUserList(userId, listId)

	if err != nil {
		return code, err
	}

	if err := p.userListService.UnlinkListFromUser(userId, listId); err != nil {
		return 500, fmt.Errorf("cannot unlink list [%s] from user [%s], %w", listId, userId, err)
	}

	if err := p.repo.DeleteProductList(userId, listId); err != nil {
		return 500, fmt.Errorf("cannot delete list [%s], %w", listId, err)
	}

	return 200, nil
}

func NewProductListService(repo repository.ProductList, userListService *UserListService) *ProductListService {
	return &ProductListService{repo: repo, userListService: userListService}
}

func (p *ProductListService) GetProductList(listId string) (model.ProductList, error) {
	return p.repo.GetProductList(listId)
}

func (p *ProductListService) CreateProductList(userId string, request model.ProductListRequest) (string, int, error) {

	request.Name = strings.Trim(request.Name, " ")

	listId, err := p.repo.CreateProductList(request)
	if err != nil {
		return "", 500, err
	}

	lists, err := p.userListService.GetUserLists(userId)
	if err != nil {
		return "", 500, fmt.Errorf("cannot get product lists for user [%s], %w", userId, err)
	}

	for _, productList := range lists {
		if strings.EqualFold(productList.Name, request.Name) {
			return "", 400, fmt.Errorf("[%s] already exists", request.Name)
		}
	}

	err = p.userListService.LinkListToUser(userId, model.UserProductListInfo{Id: listId, Name: request.Name})

	if err != nil {
		return "", 500, fmt.Errorf("cannot link list [%s] to user [%s], %w", listId, userId, err)
	}

	return listId, 200, nil
}

func (p *ProductListService) UpdateProductList(userId string, listId string, request model.ProductListRequest) (int, error) {

	code, err := p.checkUserList(userId, listId)

	if err != nil{
		return code, err
	}

	listForUpdate, err := p.repo.GetProductList(listId)
	if err != nil {
		return 500, fmt.Errorf("cannot get list for update for user [%s], %w", userId, err)
	}

	listForUpdate.Name = strings.Trim(request.Name, " ")

	if err := p.repo.UpdateProductList(listForUpdate); err != nil {
		return 500, fmt.Errorf("cannot update list info for user [%s], %w", userId, err)
	}

	listInfo := model.UserProductListInfo{Id: listId, Name: listForUpdate.Name}
	if err := p.userListService.UpdateUserList(userId, listInfo); err != nil {
		return 500, fmt.Errorf("cannot change linked list info for user [%s], %w", userId, err)
	}

	return 200, nil
}

func (p *ProductListService) checkUserList(userId string, listId string) (int, error)  {
	lists, err := p.userListService.GetUserLists(userId)
	if err != nil {
		return 500, fmt.Errorf("cannot get product lists for user [%s], %w", userId, err)
	}

	if !listContains(listId, lists) {
		return 400, errors.New("the list does not belong to the user")
	}

	return 200, nil
}

func listContains(listId string, userLists []model.UserProductListInfo) bool {

	for _, productList := range userLists {
		if productList.Id == listId {
			return true
		}
	}

	return false
}

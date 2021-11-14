package service

import (
	"fmt"
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
	"github.com/sharybkin/grocerylist-golang/pkg/extension"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ProductListService struct {
	repo            repository.ProductList
	userListService *UserListService
}

func NewProductListService(repo repository.ProductList, userListService *UserListService) *ProductListService {
	return &ProductListService{repo: repo, userListService: userListService}
}

func (p *ProductListService) GetProductList(listId string) (model.ProductList, error) {
	return p.repo.GetProductList(listId)
}

func (p *ProductListService) CreateProductList(userId string, request model.ProductListRequest) (string, error) {

	request.Name = strings.Trim(request.Name, " ")

	listId, err := p.repo.CreateProductList(request)
	if err != nil {
		return "", err
	}

	lists, err := p.userListService.GetUserLists(userId)
	if err != nil {
		return "", fmt.Errorf("cannot get product lists for user [%s], %w", userId, err)
	}

	for _, productList := range lists {
		if strings.EqualFold(productList.Name, request.Name) {
			message := fmt.Sprintf("[%s] already exists", request.Name)
			return "", &extension.BadRequestError{HttpError: extension.HttpError{Message: message}}
		}
	}

	err = p.userListService.LinkListToUser(userId, model.UserProductListInfo{Id: listId, Name: request.Name})

	if err != nil {
		return "", fmt.Errorf("cannot link list [%s] to user [%s], %w", listId, userId, err)
	}

	log.WithFields(log.Fields{
		"listId":    listId,
		"userId":    userId,
		"list name": request.Name,
	}).Debug("CreateProductList")

	return listId, nil
}

func (p *ProductListService) DeleteProductList(userId string, listId string) error {

	if err := p.checkUserList(userId, listId); err != nil {
		return err
	}

	if err := p.userListService.UnlinkListFromUser(userId, listId); err != nil {
		return fmt.Errorf("cannot unlink list [%s] from user [%s], %w", listId, userId, err)
	}

	if err := p.repo.DeleteProductList(userId, listId); err != nil {
		return fmt.Errorf("cannot delete list [%s], %w", listId, err)
	}

	log.WithFields(log.Fields{
		"listId": listId,
		"userId": userId,
	}).Debug("DeleteProductList")

	return nil
}

func (p *ProductListService) UpdateProductList(userId string, listId string, request model.ProductListRequest) error {

	if err := p.checkUserList(userId, listId); err != nil {
		return err
	}

	listForUpdate, err := p.repo.GetProductList(listId)
	if err != nil {
		return fmt.Errorf("cannot get list for update for user [%s], %w", userId, err)
	}

	listForUpdate.Name = strings.Trim(request.Name, " ")

	if err := p.repo.UpdateProductList(listForUpdate); err != nil {
		return fmt.Errorf("cannot update list info for user [%s], %w", userId, err)
	}

	listInfo := model.UserProductListInfo{Id: listId, Name: listForUpdate.Name}
	if err := p.userListService.UpdateUserList(userId, listInfo); err != nil {
		return fmt.Errorf("cannot change linked list info for user [%s], %w", userId, err)
	}

	log.WithFields(log.Fields{
		"listId":    listId,
		"userId":    userId,
		"list name": request.Name,
	}).Debug("UpdateProductList")

	return nil
}

func (p *ProductListService) checkUserList(userId string, listId string) error {
	lists, err := p.userListService.GetUserLists(userId)
	if err != nil {
		return fmt.Errorf("cannot get product lists for user [%s], %w", userId, err)
	}

	if !listContains(listId, lists) {
		return &extension.BadRequestError{HttpError: extension.HttpError{Message: "the list does not belong to the user"}}
	}

	return nil
}

func listContains(listId string, userLists []model.UserProductListInfo) bool {

	for _, productList := range userLists {
		if productList.Id == listId {
			return true
		}
	}

	return false
}

package service

import (
	"fmt"
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
	"github.com/sharybkin/grocerylist-golang/pkg/extension"
)

type UserListService struct {
	repo repository.UserList
}

func NewUserListService(repo repository.UserList) *UserListService {
	return &UserListService{repo: repo}
}

func (u *UserListService) GetUserLists(userId string) ([]model.UserProductListInfo, error) {
	return u.repo.GetUserLists(userId)
}

func (u *UserListService) LinkListToUser(userId string, listInfo model.UserProductListInfo) error {
	return u.repo.LinkListToUser(userId, listInfo)
}

func (u *UserListService) UpdateUserList(userId string, listInfo model.UserProductListInfo) error {
	return u.repo.UpdateUserList(userId, listInfo)
}

func (u *UserListService) UnlinkListFromUser(userId string, listId string) error {
	return u.repo.UnlinkListFromUser(userId, listId)
}

func (u *UserListService) checkUserList(userId string, listId string) error {
	lists, err := u.GetUserLists(userId)
	if err != nil {
		return fmt.Errorf("cannot get product lists for user [%s], %w", userId, err)
	}

	if !listContains(listId, lists) {
		return &extension.BadRequestError{HttpError: extension.HttpError{Message: "the list does not belong to the user"}}
	}

	return nil
}

package service

import (
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
)


type UserListService struct{
	repo repository.UserList
}

func NewUserListService(repo repository.UserList) *UserListService{
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

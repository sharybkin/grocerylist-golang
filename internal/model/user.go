package model

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name" `
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

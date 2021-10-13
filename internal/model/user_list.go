package model

type UserLists struct {
	UserId       string                `json:"userId" dynamodbav:"userId"`
	ProductLists []UserProductListInfo `json:"productLists" dynamodbav:"productLists"`
}

type UserProductListInfo struct {
	Id   string `json:"id" dynamodbav:"id"`
	Name string `json:"name" dynamodbav:"name"`
}

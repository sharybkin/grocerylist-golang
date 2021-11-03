package model

type ProductList struct {
	Id       string    `json:"id" dynamodbav:"id"`
	Name     string    `json:"name" dynamodbav:"name"`
	Products []Product `json:"products" dynamodbav:"products"`
}

type Product struct {
	Name   string  `json:"name" dynamodbav:"name"`
	Count  float32 `json:"count" dynamodbav:"count"`
	IsDone bool    `json:"done" dynamodbav:"done"`
}

type ProductListRequest struct {
	Name     string    `json:"name"`
}

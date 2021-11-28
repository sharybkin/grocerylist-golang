package model

type ProductExample struct {
	Name       string `json:"name" dynamodbav:"name"`
	UsageCount int    `json:"usageCount" dynamodbav:"usageCount"`
}

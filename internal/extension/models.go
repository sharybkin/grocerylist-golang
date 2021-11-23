package extension

import "github.com/sharybkin/grocerylist-golang/internal/model"

func GetValues(products map[string]model.Product) []model.Product {
	values := make([]model.Product, 0, len(products))
	for _, v := range products {
		values = append(values, v)
	}

	return values
}

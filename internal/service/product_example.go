package service

import (
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type ProductExampleService struct {
	repo     repository.ProductExample
	examples map[string]model.ProductExample
	mu       sync.Mutex
}

func NewProductExampleService(repo repository.ProductExample) *ProductExampleService {
	return &ProductExampleService{
		repo:     repo,
		examples: getStoredExamples(repo),
	}
}

func getStoredExamples(repo repository.ProductExample) map[string]model.ProductExample {
	examples, err := repo.GetExamples()

	exampleMap := make(map[string]model.ProductExample)
	if err != nil {
		log.Fatalf("failed to get product examples:  %s", err.Error())
		return exampleMap
	}

	for _, i := range examples {
		exampleMap[i.Name] = i
	}

	log.WithFields(log.Fields{
		"example count": len(exampleMap),
	}).Debug("Product examples was received from DB")

	return exampleMap
}

func (p *ProductExampleService) UpdateUsageStatistic(name string) {
	go p.updateUsageStatisticAsync(name)
}

func (p *ProductExampleService) updateUsageStatisticAsync(name string) {
	name = strings.Trim(name, " ")
	name = strings.Title(name)

	var example model.ProductExample

	p.mu.Lock()
	defer p.mu.Unlock()

	example, ok := p.examples[name]

	if !ok {
		example.Name = name
		example.UsageCount = 0
	}

	example.UsageCount++

	p.addOrUpdate(example)



	log.WithFields(log.Fields{
		"name": name,
	}).Debug("UpdateUsageStatistic")
}

func (p *ProductExampleService) addOrUpdate(example model.ProductExample) {
	if p.repo == nil {
		log.WithFields(log.Fields{
			"class":  "ProductExampleService",
			"method": "addOrUpdate",
		}).Error("null pointer exception")

		return
	}

	if err := p.repo.AddOrUpdate(example); err != nil {
		log.Errorf("Update Product examples was failed: %s", err.Error())
		return
	}

	p.examples[example.Name] = example
}

func (p *ProductExampleService) GetProductExamples() []string {
	values := make([]string, 0, len(p.examples))

	for _, v := range p.examples {
		values = append(values, v.Name)
	}
	return values
}

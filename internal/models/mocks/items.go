package mocks

import (
	"github.com/zbsss/whatsbetter/internal/models"
	"github.com/zbsss/whatsbetter/internal/rating"
)

var items = []models.Item{
	{Name: "Spicy Kebab", Rating: 4.5, NumOfReviews: 2048},
	{Name: "Kebab King", Rating: 4.9, NumOfReviews: 80},
	{Name: "Boss Kebab", Rating: 4.5, NumOfReviews: 512},
}

func processedItems() []models.Item {
	for i := range items {
		items[i].CalculatedRating = rating.Calculate(items[i].Rating, items[i].NumOfReviews)
	}
	models.SortItems(items)
	return items
}

var MockItems = processedItems()

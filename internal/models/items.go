package models

import "sort"

type Item struct {
	Name             string
	Rating           float64
	NumOfReviews     int
	CalculatedRating float64
}

func SortItems(items []Item) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].CalculatedRating > items[j].CalculatedRating
	})
}

package rating

const numNewReviews = 5

func Calculate(currentRating float64, numReviews int) float64 {
	totalScore := currentRating * float64(numReviews)
	return (totalScore + numNewReviews) / float64(numReviews+numNewReviews)
}

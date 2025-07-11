package usecase

func CalculateSum(values []float64) (sum float64) {
	for _, v := range values {
		sum += v
	}
	return
}

func CalculatedMultiply(values []float64) (multiply float64) {
	if len(values) == 0 {
		return 0
	}
	multiply = 1
	for _, v := range values {
		multiply *= v
	}
	return
}

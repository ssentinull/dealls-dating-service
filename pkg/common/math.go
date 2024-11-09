package common

func CalculatePercentage(minValue, maxValue, currentValue int64) int64 {
	progressMade := currentValue - minValue
	totalRange := maxValue - minValue
	progress := (float64(progressMade) / float64(totalRange)) * 100
	percentage := int64(progress)

	if percentage > 100 {
		return 100
	} else if percentage < 0 {
		return 0
	} else {
		return percentage
	}
}

func MaxInt(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func MinInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func MaxInt64(a int64, b int64) int64 {
	if a > b {
		return a
	}

	return b
}

func MinInt64(a int64, b int64) int64 {
	if a < b {
		return a
	}

	return b
}

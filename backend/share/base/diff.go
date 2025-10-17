package base

func Diff[T comparable](a, b []T) []T {
	// Create a map to store the elements of b
	bMap := make(map[T]bool, max(16, len(b)))
	for _, num := range b {
		bMap[num] = true
	}

	// Create a slice to store the different elements
	diff := make([]T, 0, len(a))
	for _, num := range a {
		// Check if the element is present in b
		if _, ok := bMap[num]; !ok {
			diff = append(diff, num)
		}
	}

	return diff
}
func DiffCount[T comparable](a, b []T) (diffCount int) {
	// Create a map to store the elements of b
	bMap := make(map[T]bool, max(16, len(b)))
	for _, num := range b {
		bMap[num] = true
	}
	// Create a slice to store the different elements
	for _, num := range a {
		// Check if the element is present in b
		if _, ok := bMap[num]; !ok {
			diffCount++
		}
	}

	return diffCount
}

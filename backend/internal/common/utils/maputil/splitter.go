package utils

// SplitMap splits a map into a list of maps, each with a maximum size of chunkSize.
func SplitMap[K comparable, V any](originalMap map[K]V, chunkSize int) []map[K]V {
	var result []map[K]V
	if len(originalMap) == 0 || chunkSize <= 0 {
		return result
	}

	currentChunk := make(map[K]V)
	count := 0
	for key, value := range originalMap {
		currentChunk[key] = value
		count++
		if count == chunkSize {
			result = append(result, currentChunk)
			currentChunk = make(map[K]V)
			count = 0
		}
	}

	if len(currentChunk) > 0 {
		result = append(result, currentChunk)
	}

	return result
}

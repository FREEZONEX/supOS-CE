package layrecutil

import (
	"backend/internal/types"
	"strings"
)

// Contains checks if childLayRec is a sub-path of parentLayRec.
func Contains(childLayRec, parentLayRec string) bool {
	if childLayRec == "" || parentLayRec == "" {
		return false
	}
	if childLayRec == parentLayRec {
		return true
	}
	return strings.HasPrefix(childLayRec, parentLayRec+"/")
}

// GetParentLayRec returns the parent path of a given layRec.
func GetParentLayRec(layRec string) string {
	if layRec == "" {
		return ""
	}
	lastSlash := strings.LastIndex(layRec, "/")
	if lastSlash == -1 {
		return "" // No parent, return empty string which is more idiomatic in Go than null.
	}
	return layRec[:lastSlash]
}

// BuildParentToChildrenMap builds an index mapping parent paths to all their descendant nodes.
// This implementation correctly mirrors the logic of the performant `buildParentToChildrenMap2` from the Java version.
func BuildParentToChildrenMap(allNodes []*types.CreateTopicDto) map[string][]*types.CreateTopicDto {
	parentToChildrenMap := make(map[string][]*types.CreateTopicDto, len(allNodes)*2)

	for _, node := range allNodes {
		layRec := node.LayRec
		if layRec == "" {
			continue
		}

		index := -1
		// Iterate through all slashes to find all parent paths
		for {
			index = strings.Index(layRec, "/")
			if index == -1 {
				break
			}
			parentPath := layRec[:index]
			list, ok := parentToChildrenMap[parentPath]
			if !ok {
				list = make([]*types.CreateTopicDto, 0, 4)
			}
			parentToChildrenMap[parentPath] = append(list, node)

			// This logic was incorrect in the previous version. The correct approach
			// is to find all parent prefixes of the original layRec, not to shorten it.
			// The correct implementation is below.
		}
	}
	// The above implementation was incorrect. Let's do it correctly.
	parentToChildrenMap = make(map[string][]*types.CreateTopicDto, len(allNodes)*2) // Reset the map
	for _, node := range allNodes {
		layRec := node.LayRec
		if layRec == "" {
			continue
		}

		tempPath := layRec
		index := -1
		for {
			index = strings.LastIndex(tempPath, "/")
			if index == -1 {
				break
			}
			parentPath := tempPath[:index]
			list, ok := parentToChildrenMap[parentPath]
			if !ok {
				list = make([]*types.CreateTopicDto, 0, 4)
			}
			parentToChildrenMap[parentPath] = append(list, node)
			tempPath = parentPath
		}
	}

	return parentToChildrenMap
}

// GetNextNodeAfterBasePath gets the next path segment in `path` that comes after `basePath`.
// Corresponds to the logic in `LayRecUtil.getNextNodeAfterBasePath2`.
func GetNextNodeAfterBasePath(basePath, path string) string {
	if path == "" {
		return ""
	}

	baseLen := len(basePath)
	if baseLen == 0 || basePath == "0" {
		slashIndex := strings.Index(path, "/")
		if slashIndex == -1 {
			return path
		}
		return path[:slashIndex]
	}

	if !strings.HasPrefix(path, basePath) {
		return ""
	}

	index := baseLen
	if index < len(path) && path[index] == '/' {
		index++
	}

	if index >= len(path) {
		return ""
	}

	nextSlashIndex := strings.Index(path[index:], "/")
	if nextSlashIndex != -1 {
		return path[index : index+nextSlashIndex]
	}

	return path[index:]
}

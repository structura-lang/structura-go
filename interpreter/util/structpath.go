package util

func GetPath(m map[string]any, path []string) (any, bool) {
	var current any = m

	for _, key := range path {
		obj, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		current, ok = obj[key]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

func SetPath(m map[string]any, path []string, value any) bool {
	if len(path) == 0 {
		return false
	}

	var current any = m
	for i := 0; i < len(path)-1; i++ {
		key := path[i]

		obj, ok := current.(map[string]any)
		if !ok {
			return false
		}

		next, ok := obj[key]
		if !ok {
			return false
		}
		current = next
	}

	parentMap, ok := current.(map[string]any)
	if !ok {
		return false
	}

	leafKey := path[len(path)-1]
	parentMap[leafKey] = value
	return true
}

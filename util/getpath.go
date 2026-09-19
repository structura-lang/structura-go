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

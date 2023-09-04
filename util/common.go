package util

func Contains[T string | int | bool](container []T, include []T) bool {
	includeMap := make(map[T]struct{})
	for _, a := range include {
		includeMap[a] = struct{}{}
	}
	for _, b := range container {
		if _, ok := includeMap[b]; !ok {
			return false
		}
	}
	return true
}

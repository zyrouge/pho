package utils

func SliceContains[T comparable](slice []T, value T) bool {
	for _, x := range slice {
		if x == value {
			return true
		}
	}
	return false
}

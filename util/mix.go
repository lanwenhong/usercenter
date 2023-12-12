package util

func Foundin[T comparable](vlist []T, v T) bool {
	for _, xv := range vlist {
		if xv == v {
			return true
		}
	}
	return false
}

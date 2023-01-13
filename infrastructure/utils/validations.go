package utils

type constrains interface {
	~string | ~int | ~float32
}

func ItemInSlice[T constrains](a T, list []T) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

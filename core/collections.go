package core

import "synchronex/common/hashcode"

func Contains[T any](list []T, item T) bool {
	for _, v := range list {
		if hashcode.HashCode(v) == hashcode.HashCode(item) {
			return true
		}
	}
	return false
}

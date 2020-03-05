package utils

import "reflect"

// IsNil checks nil with reflection
func IsNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

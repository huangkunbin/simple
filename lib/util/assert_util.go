package util

import (
	"fmt"
)

func Assert(check bool, message string, args ...interface{}) {
	if !check {
		panic(fmt.Sprintf(message, args...))
	}
}

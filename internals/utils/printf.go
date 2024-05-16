package util

import "fmt"

func P(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

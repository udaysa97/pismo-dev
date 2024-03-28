package utils

import (
	"fmt"
	"strconv"
)

func ToInt(value any) int {
	switch value.(type) {
	case int:
		return value.(int)
	case float64:
		return int(value.(float64))
	case string:
		val, _ := strconv.ParseInt(value.(string), 10, 64)
		return int(val)
	case int64:
		return int(value.(int64))
	}
	return 0
}

func ToString(value any) string {
	switch value.(type) {
	case int:
	case float64:
	case int64:
		return fmt.Sprint(value)
	case string:
		return value.(string)
	}
	return ""
}

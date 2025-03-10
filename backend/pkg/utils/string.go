package utils

import "strconv"

func StrToInt64(key string) int64 {
	value, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

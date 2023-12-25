package util

import "strconv"

func StringToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

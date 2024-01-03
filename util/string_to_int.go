package util

import "strconv"

func StringToInt(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

func IntToString(val int) string {
	return strconv.Itoa(val)
}

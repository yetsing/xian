package parser

import (
	"strconv"
	"strings"
)

func toint64(s string) (int64, error) {
	s = strings.ReplaceAll(s, "_", "")
	return strconv.ParseInt(s, 0, 64)
}

func tofloat64(s string) (float64, error) {
	s = strings.ReplaceAll(s, "_", "")
	return strconv.ParseFloat(s, 64)
}

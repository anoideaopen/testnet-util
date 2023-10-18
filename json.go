package utils

import (
	"strings"
)

func MakeJSONArrayFromNdJSON(body []byte) string {
	lenToCt := 2
	str := string(body)
	cn := strings.Split(str, "\n")
	str = strings.Replace(str, "\n", ",", len(cn)-lenToCt)
	return "[" + str + "]"
}

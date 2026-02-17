package main

import (
	"strconv"
	"strings"
)

func ParseHTTPStatus(status string) (code int, reason string) {
	p := strings.SplitN(status, " ", 3)
	code, _ = strconv.Atoi(p[1])
	reason = p[2]
	return code, reason
}

package main

import "fmt"

func BuildHTTPResponse(status string, headers string, body string) string {
	return fmt.Sprintf("%s\r\n%s\r\n%s", status, headers, body)
}

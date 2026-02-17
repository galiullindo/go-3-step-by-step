package main

import "fmt"

func BuildHTTPRequest(method string, url string, host string, headers string, body string) string {
	return fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\n%s\r\n%s", method, url, host, headers, body)
}

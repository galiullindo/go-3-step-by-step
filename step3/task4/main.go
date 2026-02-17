package main

import (
	"fmt"
	"strings"
)

func MakeCurlCommand(method string, url string, headers string, body string) string {
	var args []string

	// -X ...
	if method != "GET" {
		args = append(args, fmt.Sprintf("-X %s", method))
	}

	// -H '...' -H '...'
	if headers != "" {
		sliceHeaders := make([]string, 0)
		for _, header := range strings.Split(headers, "\n") {
			if header != "" {
				sliceHeaders = append(sliceHeaders, fmt.Sprintf("-H '%s'", strings.TrimSpace(header)))
			}
		}
		args = append(args, strings.Join(sliceHeaders, " "))
	}

	// --data '...'
	if body != "" {
		args = append(args, fmt.Sprintf("--data '%s'", body))
	}

	args = append(args, url)

	return fmt.Sprintf("curl %s", strings.Join(args, " "))
}

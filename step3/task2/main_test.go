package main

import "testing"

func Test(t *testing.T) {
	var tests = []struct {
		name     string
		status   string
		headers  string
		body     string
		expected string
	}{
		{
			name:     "Empty response",
			status:   "",
			headers:  "",
			body:     "",
			expected: "\r\n\r\n",
		},
		{
			name:     "Full response",
			status:   "HTTP/1.1 404 Not Found",
			headers:  "Content-Type: text/html\r\nDate: Tue, 3 Jun 2025 19:00:00 GMT\r\n",
			body:     "<h1>Not Found</h1>",
			expected: "HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\nDate: Tue, 3 Jun 2025 19:00:00 GMT\r\n\r\n<h1>Not Found</h1>",
		},
		{
			name:     "Response without headers",
			status:   "HTTP/1.1 404 Not Found",
			headers:  "",
			body:     "<h1>Not Found</h1>",
			expected: "HTTP/1.1 404 Not Found\r\n\r\n<h1>Not Found</h1>",
		},
		{
			name:     "Response without body",
			status:   "HTTP/1.1 404 Not Found",
			headers:  "Content-Type: text/html\r\nDate: Tue, 3 Jun 2025 19:00:00 GMT\r\n",
			body:     "",
			expected: "HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\nDate: Tue, 3 Jun 2025 19:00:00 GMT\r\n\r\n",
		},
		{
			name:     "Response without headers and body",
			status:   "HTTP/1.1 404 Not Found",
			headers:  "",
			body:     "",
			expected: "HTTP/1.1 404 Not Found\r\n\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := BuildHTTPResponse(test.status, test.headers, test.body)
			if got != test.expected {
				t.Errorf("unexpected request: %v expected %v\n", got, test.expected)
			}
		})
	}
}

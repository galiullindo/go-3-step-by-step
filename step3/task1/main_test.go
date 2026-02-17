package main

import "testing"

func Test(t *testing.T) {
	var tests = []struct {
		name     string
		method   string
		url      string
		host     string
		headers  string
		body     string
		expected string
	}{
		{
			name:     "Empty request",
			method:   "",
			url:      "",
			host:     "",
			headers:  "",
			body:     "",
			expected: "  HTTP/1.1\r\nHost: \r\n\r\n",
		},
		{
			name:     "Full request",
			method:   "POST",
			url:      "/api/users",
			host:     "example.com",
			headers:  "Content-Type: application/json\r\nAuthorization: Bearer abc123\r\n",
			body:     "{\"name\": \"John Doe\", \"email\": \"johndoe@example.com\", \"password\": \"123456\"}",
			expected: "POST /api/users HTTP/1.1\r\nHost: example.com\r\nContent-Type: application/json\r\nAuthorization: Bearer abc123\r\n\r\n{\"name\": \"John Doe\", \"email\": \"johndoe@example.com\", \"password\": \"123456\"}",
		},
		{
			name:     "Request without headers",
			method:   "POST",
			url:      "/api/users",
			host:     "example.com",
			body:     "{\"name\": \"John Doe\", \"email\": \"johndoe@example.com\", \"password\": \"123456\"}",
			expected: "POST /api/users HTTP/1.1\r\nHost: example.com\r\n\r\n{\"name\": \"John Doe\", \"email\": \"johndoe@example.com\", \"password\": \"123456\"}",
		},
		{
			name:     "Request without body",
			method:   "POST",
			url:      "/api/users",
			host:     "example.com",
			headers:  "Content-Type: application/json\r\nAuthorization: Bearer abc123\r\n",
			expected: "POST /api/users HTTP/1.1\r\nHost: example.com\r\nContent-Type: application/json\r\nAuthorization: Bearer abc123\r\n\r\n",
		},
		{
			name:     "Request without headers and body",
			method:   "POST",
			url:      "/api/users",
			host:     "example.com",
			expected: "POST /api/users HTTP/1.1\r\nHost: example.com\r\n\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := BuildHTTPRequest(test.method, test.url, test.host, test.headers, test.body)
			if got != test.expected {
				t.Errorf("unexpected request: %v expected %v\n", got, test.expected)
			}
		})
	}
}

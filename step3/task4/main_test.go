package main

import "testing"

func Test(t *testing.T) {
	var tests = []struct {
		name     string
		method   string
		url      string
		headers  string
		body     string
		expected string
	}{
		{
			name:     "Full POST",
			method:   "POST",
			url:      "https://example.com/api/users",
			headers:  "Content-Type: application/json\nAuthorization: Bearer abc123\n",
			body:     `{"name":"John Doe","email":"johndoe@example.com","password":"123456"}`,
			expected: "curl -X POST -H 'Content-Type: application/json' -H 'Authorization: Bearer abc123' --data '{\"name\":\"John Doe\",\"email\":\"johndoe@example.com\",\"password\":\"123456\"}' https://example.com/api/users",
		},
		{
			name:     "POST without headers",
			method:   "POST",
			url:      "https://example.com/api/users",
			body:     `{"name":"John Doe","email":"johndoe@example.com","password":"123456"}`,
			expected: "curl -X POST --data '{\"name\":\"John Doe\",\"email\":\"johndoe@example.com\",\"password\":\"123456\"}' https://example.com/api/users",
		},
		{
			name:     "POST without data",
			method:   "POST",
			url:      "https://example.com/api/users",
			headers:  "Content-Type: application/json\nAuthorization: Bearer abc123\n",
			expected: "curl -X POST -H 'Content-Type: application/json' -H 'Authorization: Bearer abc123' https://example.com/api/users",
		},
		{
			name:     "POST without headers and data",
			method:   "POST",
			url:      "https://example.com/api/users",
			expected: "curl -X POST https://example.com/api/users",
		},
		{
			name:     "Full GET",
			method:   "GET",
			url:      "https://example.com/api/users",
			headers:  "Content-Type: application/json\nAuthorization: Bearer abc123\n",
			body:     `{"name":"John Doe","email":"johndoe@example.com","password":"123456"}`,
			expected: "curl -H 'Content-Type: application/json' -H 'Authorization: Bearer abc123' --data '{\"name\":\"John Doe\",\"email\":\"johndoe@example.com\",\"password\":\"123456\"}' https://example.com/api/users",
		},
		{
			name:     "GET without headers",
			method:   "GET",
			url:      "https://example.com/api/users",
			body:     `{"name":"John Doe","email":"johndoe@example.com","password":"123456"}`,
			expected: "curl --data '{\"name\":\"John Doe\",\"email\":\"johndoe@example.com\",\"password\":\"123456\"}' https://example.com/api/users",
		},
		{
			name:     "GET without data",
			method:   "GET",
			url:      "https://example.com/api/users",
			headers:  "Content-Type: application/json\nAuthorization: Bearer abc123\n",
			expected: "curl -H 'Content-Type: application/json' -H 'Authorization: Bearer abc123' https://example.com/api/users",
		},
		{
			name:     "GET without headers and data",
			method:   "GET",
			url:      "https://example.com/api/users",
			expected: "curl https://example.com/api/users",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := MakeCurlCommand(test.method, test.url, test.headers, test.body)
			if got != test.expected {
				t.Errorf("unexpected curl: %q expected %q\n", got, test.expected)
			}
		})
	}
}

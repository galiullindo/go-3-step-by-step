package main

import "testing"

func Test(t *testing.T) {
	var tests = []struct {
		name           string
		status         string
		expectedCode   int
		expectedReason string
	}{
		{
			name:           "Status 200",
			status:         "HTTP/1.1 200 OK",
			expectedCode:   200,
			expectedReason: "OK",
		},
		{
			name:           "Status 418",
			status:         "HTTP/1.1 418 I'm a teapot",
			expectedCode:   418,
			expectedReason: "I'm a teapot",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotCode, gotReason := ParseHTTPStatus(test.status)
			if gotCode != test.expectedCode || gotReason != test.expectedReason {
				t.Errorf("unexpected code or reason: %v, %v expected %v, %v\n", gotCode, gotReason, test.expectedCode, test.expectedReason)
			}
		})
	}
}

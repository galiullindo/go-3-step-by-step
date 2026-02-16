package main

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"
)

var diff = 5 * time.Millisecond

func Test(t *testing.T) {

	var tests = []struct {
		name        string
		timeout     time.Duration
		numbers     []int
		fn          func(i int) int
		maxWorkers  int
		expected    []int
		expectedErr error
	}{
		{
			name:       "Case multiply by 2",
			timeout:    10 * time.Millisecond,
			numbers:    []int{0, 1, 2, 3, 4, 5},
			fn:         func(i int) int { return i * 2 },
			maxWorkers: 5,
			expected:   []int{0, 2, 4, 6, 8, 10},
		},
		{
			name:       "Case square",
			timeout:    10 * time.Millisecond,
			numbers:    []int{0, 1, 2, 3, 4, 5},
			fn:         func(i int) int { return i * i },
			maxWorkers: 5,
			expected:   []int{0, 1, 4, 9, 16, 25},
		},
		{
			name:        "Case timeout 0",
			timeout:     0 * time.Millisecond,
			numbers:     []int{0, 1, 2, 3, 4, 5},
			fn:          func(i int) int { return i * 2 },
			maxWorkers:  5,
			expectedErr: context.DeadlineExceeded,
		},
		{
			name:        "Case testing the max number of workers control for 5 workers",
			timeout:     10 * time.Millisecond,
			numbers:     []int{0, 1, 2, 3, 4, 5},
			fn:          func(i int) int { time.Sleep(8 * time.Millisecond); return i * i },
			maxWorkers:  5,
			expectedErr: context.DeadlineExceeded,
		},
		{
			name:       "Case testing the max number of workers control for 6 workers",
			timeout:    10 * time.Millisecond,
			numbers:    []int{0, 1, 2, 3, 4, 5},
			fn:         func(i int) int { time.Sleep(8 * time.Millisecond); return i * i },
			maxWorkers: 6,
			expected:   []int{0, 1, 4, 9, 16, 25},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			start := time.Now()
			got, err := ParallelMapCtx(ctx, test.numbers, test.fn, test.maxWorkers)
			duration := time.Since(start)

			expectedTime := test.timeout + diff
			if duration > expectedTime {
				t.Errorf("unexpected execution time %v expected %v\n", duration, expectedTime)
			}

			if !errors.Is(err, test.expectedErr) {
				t.Errorf("unexpected error %v expected %v\n", err, test.expectedErr)
			}

			if !slices.Equal(got, test.expected) {
				t.Errorf("unexpected value %v expected %v\n", got, test.expected)
			}
		})
	}
}

package main

import (
	"testing"
	"time"
)

func Test_writeUntil(t *testing.T) {
	for _, test := range []struct {
		in     time.Time
		expect time.Time
	}{
		{
			in:     time.Date(1, 1, 1, 1, 58, 59, 0, time.UTC),
			expect: time.Date(1, 1, 1, 1, 59, 0, 0, time.UTC),
		},
		{
			in:     time.Date(1, 1, 1, 1, 59, 0, 0, time.UTC),
			expect: time.Date(1, 1, 1, 1, 59, 0, 0, time.UTC),
		},
		{
			in:     time.Date(1, 1, 1, 1, 59, 1, 1, time.UTC),
			expect: time.Date(1, 1, 1, 2, 0, 0, 0, time.UTC),
		},
		{
			in:     time.Date(1, 1, 1, 1, 59, 30, 0, time.UTC),
			expect: time.Date(1, 1, 1, 2, 0, 0, 0, time.UTC),
		},
	} {
		res := writeUntil(test.in)
		if !res.Equal(test.expect) {
			t.Errorf("expected: %v, received: %v\n", test.expect, res)
		}
	}
}

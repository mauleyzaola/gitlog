package cmd

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	d1 := time.Date(2018, 11, 25, 23, 59, 59, 0, time.UTC)
	d2 := time.Date(2018, 3, 4, 23, 59, 59, 0, time.UTC)
	cases := []struct {
		input    string
		expected *time.Time
		error    bool
	}{
		{
			input:    "",
			expected: nil,
			error:    false,
		},
		{
			input:    "20181125",
			expected: &d1,
			error:    false,
		},
		{
			input:    "20180304",
			expected: &d2,
			error:    false,
		},
		{
			input:    "20180304103015",
			expected: nil,
			error:    true,
		},
		{
			input:    "x",
			expected: nil,
			error:    true,
		},
	}

	for i, c := range cases {
		actual, err := parseDate(c.input)
		if c.error {
			if err == nil {
				t.Errorf("[%d] - expected: error actual: nil", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] - expected: nil actual: %s", i, err)
				continue
			}
			if c.expected == nil || actual == nil {
				if c.expected != actual {
					t.Errorf("[%d] - expected:%v actual:%v", i, c.expected, actual)
				}
				continue
			}
			if c.expected.Unix() != actual.Unix() {
				t.Errorf("[%d] - expected:%v actual:%v", i, c.expected, actual)
			}
		}
	}
}

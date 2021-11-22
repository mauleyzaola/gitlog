package internal

import (
	"strings"
	"testing"
	"time"
)

func TestCommits_Filter(t *testing.T) {
	john := &Author{
		Name:  "John",
		Email: "john.foo@mail.com",
	}
	mary := &Author{
		Name:  "Mary",
		Email: "mary.bar@mail.com",
	}
	jan2018 := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	feb2018 := time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC)
	mar2018 := time.Date(2018, 3, 1, 0, 0, 0, 0, time.UTC)
	apr2018 := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)

	input := Commits{
		{
			Author: john,
			Date:   jan2018,
		},
		{
			Author: john,
			Date:   jan2018,
		},
		{
			Author: john,
			Date:   feb2018,
		},
		{
			Author: john,
			Date:   mar2018,
		},
		{
			Author: mary,
			Date:   jan2018,
		},
		{
			Author: mary,
			Date:   jan2018,
		},
	}

	if expected, actual := 4, len(input.Filter(strings.Fields("john.foo@mail.com"), nil, nil)); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := 6, len(input.Filter(strings.Fields("john.foo@mail.com mary.bar@mail.com"), nil, nil)); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := 0, len(input.Filter(nil, &apr2018, nil)); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := 6, len(input.Filter(nil, &jan2018, nil)); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := 1, len(input.Filter(nil, &mar2018, nil)); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := 2, len(input.Filter(nil, &feb2018, nil).Filter(nil, nil, &mar2018)); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

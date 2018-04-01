package monthx

import (
	"fmt"
	"testing"
	"time"
)

// mustParse parses value in the format YYYY-MM-DD failing the test on error.
func mustParse(t *testing.T, value string) time.Time {
	const layout = "2006-01-02"
	d, err := time.Parse(layout, value)
	if err != nil {
		t.Fatalf("time.Parse(%q, %q) unexpected error: %v", layout, value, err)
	}
	return d
}

func TestStartOfMonth(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2016-01-13"), mustParse(t, "2016-01-01")},
		{mustParse(t, "2016-01-01"), mustParse(t, "2016-01-01")},
		{mustParse(t, "2016-12-30"), mustParse(t, "2016-12-01")},
	}

	for _, c := range cases {
		got := StartOfMonth(c.in)
		if got != c.want {
			t.Errorf("StartOfMonth(%s) => %s, want %s", c.in, got, c.want)
		}
	}
}

func TestEndOfMonth(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2016-01-01"), mustParse(t, "2016-01-31")},
		{mustParse(t, "2016-01-31"), mustParse(t, "2016-01-31")},
		{mustParse(t, "2016-11-01"), mustParse(t, "2016-11-30")},
		{mustParse(t, "2016-12-31"), mustParse(t, "2016-12-31")},
		// Leap test.
		{mustParse(t, "2012-02-01"), mustParse(t, "2012-02-29")},
		{mustParse(t, "2013-02-01"), mustParse(t, "2013-02-28")},
	}

	for _, c := range cases {
		got := EndOfMonth(c.in)
		if got != c.want {
			t.Errorf("EndOfMonth(%s) => %s, want %s", c.in, got, c.want)
		}
	}
}

func TestMonthsTo(t *testing.T) {
	day := 24 * time.Hour
	cases := []struct {
		in   time.Time
		want int
	}{
		{time.Now(), 1},
		{time.Now().Add(day * 35), 1},
		{time.Now().Add(day * 65), 2},
		{time.Now().Add(day * 370), 12},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := MonthsTo(tc.in)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	m := Month(January)
	monthName := m.String()
	if monthName != "January" {
		t.Fatal("For", "1",
			"expected", "January",
			"got", monthName,
		)
	}
}

func TestLastDay(t *testing.T) {
	lastDays := [12]int{
		31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31,
	}

	// Normal years
	for i := 1; i <= 12; i++ {
		m := Month(i)
		actual := m.LastDay(2015)
		expected := lastDays[i-1]
		if actual != expected {
			t.Error(
				"For", m,
				"expected", expected,
				"got", actual)
		}
	}

	// Leap year
	actual := Month(February).LastDay(2012)
	if actual != 29 {
		t.Error(
			"For leap year",
			"expected", 29,
			"got", actual)
	}
}

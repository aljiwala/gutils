// Package monthx provides functionality, pertaining to months of the year, that
// is not available in standard Go libraries.
package monthx

import (
	"time"

	"github.com/aljiwala/gutils/timex"
)

// daysIn should return numeric value as in total days in given month and year.
func daysIn(year int, m time.Month) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// -----------------------------------------------------------------------------

// TimeMonth should returns time.Month representation of the given month.
func (m Month) TimeMonth() time.Month {
	return time.Month(m)
}

// Int should return int representation of give month.
func (m Month) Int() int {
	return int(m.TimeMonth())
}

// String returns an English language based representation.
// e.g. January, February, etc.
func (m Month) String() string {
	return m.TimeMonth().String()
}

// -----------------------------------------------------------------------------

// IsJanuary should return true if the given month is January.
func (m Month) IsJanuary() bool {
	return IsJanuary(m.TimeMonth())
}

// IsFebruary should return true if the given month is February.
func (m Month) IsFebruary() bool {
	return IsFebruary(m.TimeMonth())
}

// IsMarch should return true if the given month is March.
func (m Month) IsMarch() bool {
	return IsMarch(m.TimeMonth())
}

// IsApril should return true if the given month is April.
func (m Month) IsApril() bool {
	return IsApril(m.TimeMonth())
}

// IsMay should return true if the given month is May.
func (m Month) IsMay() bool {
	return IsMay(m.TimeMonth())
}

// IsJune should return true if the given month is June.
func (m Month) IsJune() bool {
	return IsJune(m.TimeMonth())
}

// IsJuly should return true if the given month is July.
func (m Month) IsJuly() bool {
	return IsJuly(m.TimeMonth())
}

// IsAugust should return true if the given month is August.
func (m Month) IsAugust() bool {
	return IsAugust(m.TimeMonth())
}

// IsSeptember should return true if the given month is September.
func (m Month) IsSeptember() bool {
	return IsSeptember(m.TimeMonth())
}

// IsOctober should return true if the given month is October.
func (m Month) IsOctober() bool {
	return IsOctober(m.TimeMonth())
}

// IsNovember should return true if the given month is November.
func (m Month) IsNovember() bool {
	return IsNovember(m.TimeMonth())
}

// IsDecember should return true if the given month is December.
func (m Month) IsDecember() bool {
	return IsDecember(m.TimeMonth())
}

// -----------------------------------------------------------------------------

// IsJanuary should return true if given month is January.
func IsJanuary(m time.Month) bool {
	return m.String() == January.String()
}

// IsFebruary should return true if given month is February.
func IsFebruary(m time.Month) bool {
	return m.String() == February.String()
}

// IsMarch should return true if given month is March.
func IsMarch(m time.Month) bool {
	return m.String() == March.String()
}

// IsApril should return true if given month is April.
func IsApril(m time.Month) bool {
	return m.String() == April.String()
}

// IsMay should return true if given month is May.
func IsMay(m time.Month) bool {
	return m.String() == May.String()
}

// IsJune should return true if given month is June.
func IsJune(m time.Month) bool {
	return m.String() == June.String()
}

// IsJuly should return true if given month is July.
func IsJuly(m time.Month) bool {
	return m.String() == July.String()
}

// IsAugust should return true if given month is August.
func IsAugust(m time.Month) bool {
	return m.String() == August.String()
}

// IsSeptember should return true if given month is September.
func IsSeptember(m time.Month) bool {
	return m.String() == September.String()
}

// IsOctober should return true if given month is October.
func IsOctober(m time.Month) bool {
	return m.String() == October.String()
}

// IsNovember should return true if given month is November.
func IsNovember(m time.Month) bool {
	return m.String() == November.String()
}

// IsDecember should return true if given month is December.
func IsDecember(m time.Month) bool {
	return m.String() == December.String()
}

// -----------------------------------------------------------------------------

// DaysIn should return total number of days of particular month and year.
//
// Example:
//   - timex.DaysIn(2020, monthx.June.TimeMonth()) returns 30.
func (m Month) DaysIn(year int) int {
	return daysIn(year, m.TimeMonth())
}

// LastDay returns the last numeric day of the given month and year. Also takes
// leap years into account.
//
// Example:
//   - monthx.June.LastDay(1992) should return 30.
func (m Month) LastDay(year int) int {
	return timex.LastDayOfMonth(year, time.Month(m))
}

// DaysIn should return total number of days of particular month and year.
//
// Example:
//   - timex.DaysIn(2020, monthx.June.TimeMonth()) returns 30.
func DaysIn(year int, month time.Month) int {
	return daysIn(year, month)
}

// StartOfMonth returns the first day of the month of date.
func StartOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

// EndOfMonth returns the last day of the month of date.
func EndOfMonth(date time.Time) time.Time {
	// Go to the next month, then a day of 0 removes a day leaving us at the
	// last day of dates month.
	return time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())
}

// GetStartAndEndOfMonth should return start and end time.Time value by getting
// the same type.
func GetStartAndEndOfMonth(t time.Time) (start, end time.Time) {
	year, month, _ := t.Date()
	start = time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	end = start.AddDate(0, 1, -1)
	return
}

// MonthsTo returns the number of months from the current date to the given date.
// The number of months is always rounded down, with a minimal value of 1.
//
// Example:
//   - MonthsTo(time.Now().Add(24 * time.Hour * 70)) should return 2.
//
// Dates in the past are not supported, and their behaviour is undefined!
func MonthsTo(a time.Time) int {
	var days int
	startDate := time.Now()
	lastDayOfYear := func(t time.Time) time.Time {
		return time.Date(t.Year(), 12, 31, 0, 0, 0, 0, t.Location())
	}

	firstDayOfNextYear := func(t time.Time) time.Time {
		return time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location())
	}

	cur := startDate
	for cur.Year() < a.Year() {
		// add 1 to count the last day of the year too.
		days += lastDayOfYear(cur).YearDay() - cur.YearDay() + 1
		cur = firstDayOfNextYear(cur)
	}

	days += a.YearDay() - cur.YearDay()
	if startDate.AddDate(0, 0, days).After(a) {
		days--
	}

	months := (days / 30)
	if months == 0 {
		months = 1
	}

	return months
}

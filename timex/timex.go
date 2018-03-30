package timex

import (
	"math"
	"strings"
	"time"
)

// TimeSpan holds the block of the time period.
type TimeSpan struct {
	Begin time.Time
	End   time.Time
}

// absValue returns the abs value if needed.
func absValue(needsAbs bool, value int64) int64 {
	if needsAbs && value < 0 {
		return -value
	}
	return value
}

// now returns the current local time.
func now() time.Time {
	return time.Now()
}

// nowIn should return time.Time in provided location.
func nowIn(loc *time.Location) time.Time {
	return now().In(loc)
}

// create returns a new carbon pointe. It is a helper function to create new
// dates.
func create(y int, mon time.Month, d, h, m, s, ns int, l *time.Location) time.Time {
	return time.Date(y, mon, d, h, m, s, ns, l)
}

// isPast determines if the current time is in the past, i.e. less (before)
// than Now().
func isPast(t time.Time) bool {
	return t.Before(Now())
}

// isToday should check if the given date is matched with todays date (alias of
// `isSameDay`).
func isToday(t time.Time) bool {
	return isSameDay(t)
}

// isFuture determines if the current time is in the future, i.e. greater (after)
// than Now().
func isFuture(t time.Time) bool {
	return t.After(Now())
}

// isCurrentDay determines if the current time is in the current day.
func isCurrentDay(t time.Time) bool {
	return t.Day() == Now().Day()
}

// isCurrentMonth determines if the current time is in the current month.
func isCurrentMonth(t time.Time) bool {
	return t.Month() == Now().Month()
}

// isCurrentYear determines if the current time is in the current year.
func isCurrentYear(t time.Time) bool {
	return t.Year() == Now().Year()
}

// isSameDay checks if the given date is the same day as the current day.
func isSameDay(t time.Time) bool {
	n := nowIn(t.Location())
	return t.Year() == n.Year() && t.Month() == n.Month() && t.Day() == n.Day()
}

// isSameMonth checks if month of the given date is same as month of the current
// date.
func isSameMonth(t time.Time, sameYear bool) bool {
	m := nowIn(t.Location()).Month()
	if sameYear {
		return isSameYear(t) && m == t.Month()
	}
	return m == t.Month()
}

// isSameYear checks if given date is in current year.
func isSameYear(t time.Time) bool {
	return t.Year() == nowIn(t.Location()).Year()
}

// addMilliSecond adds a millisecond to the time.
// Positive values travels forward while negative values travels into the past.
func addMilliSecond(t time.Time) time.Time {
	return addMilliSeconds(t, 1)
}

// addMilliSeconds adds milliseconds to the current time.
func addMilliSeconds(t time.Time, s time.Duration) time.Time {
	d := time.Duration(s) * time.Millisecond
	return t.Add(d)
}

// addSecond adds a second to the time.
// Positive values travels forward while negative values travels into the past.
func addSecond(t time.Time) time.Time {
	return addSeconds(t, 1)
}

// addSeconds adds seconds to the current time.
func addSeconds(t time.Time, s time.Duration) time.Time {
	d := time.Duration(s) * time.Second
	return t.Add(d)
}

// addMinute adds a minute to the time.
// Positive values travels forward while negative values travels into the past.
func addMinute(t time.Time) time.Time {
	return addMinutes(t, 1)
}

// addMinutes adds minutes to the current time.
func addMinutes(t time.Time, m time.Duration) time.Time {
	d := time.Duration(m) * time.Minute
	return t.Add(d)
}

// addHour adds an hour to the time.
// Positive values travels forward while negative values travels into the past.
func addHour(t time.Time) time.Time {
	return addHours(t, 1)
}

// addHours adds hours to the current time.
func addHours(t time.Time, h int) time.Time {
	d := time.Duration(h) * time.Hour
	return t.Add(d)
}

// addDay adds a day to the time.
// Positive values travels forward while negative values travels into the past.
func addDay(t time.Time) time.Time {
	return addDays(t, 1)
}

// addDays adds days to the current time.
func addDays(t time.Time, d int) time.Time {
	return t.AddDate(0, 0, d)
}

// addWeek adds a week to the time.
// Positive values travels forward while negative values travels into the past.
func addWeek(t time.Time) time.Time {
	return addWeeks(t, 1)
}

// addWeeks adds weeks to the current time.
func addWeeks(t time.Time, w int) time.Time {
	return t.AddDate(0, 0, daysPerWeek*w)
}

// addMonth adds a month to the time.
// Positive values travels forward while negative values travels into the past.
func addMonth(t time.Time) time.Time {
	return addMonths(t, 1)
}

// addMonths adds months to the current time.
func addMonths(t time.Time, m int) time.Time {
	return t.AddDate(0, m, 0)
}

// addQuarter adds a quarter to the time.
// Positive values travels forward while negative values travels into the past.
func addQuarter(t time.Time) time.Time {
	return addQuarters(t, 1)
}

// addQuarters adds quarters to the current time.
func addQuarters(t time.Time, q int) time.Time {
	return t.AddDate(0, monthsPerQuarter*q, 0)
}

// addYear adds a year to the time.
// Positive values travels forward while negative values travels into the past.
func addYear(t time.Time) time.Time {
	return addYears(t, 1)
}

// addYears adds years to the current time.
func addYears(t time.Time, y int) time.Time {
	return t.AddDate(y, 0, 0)
}

// addCentury adds a century to the time.
// Positive values travels forward while negative values travels into the past.
func addCentury(t time.Time) time.Time {
	return addCenturies(t, 1)
}

// addCenturies adds centuries to the current time.
func addCenturies(t time.Time, c int) time.Time {
	return t.AddDate(yearsPerCenturies*c, 0, 0)
}

// -----------------------------------------------------------------------------

// subMilliSecond removes a millisecond to the time.
// Positive values travels forward while negative values travels into the past.
func subMilliSecond(t time.Time) time.Time {
	return subMilliSeconds(t, 1)
}

// subMilliSeconds removes milliseconds to the current time.
func subMilliSeconds(t time.Time, s time.Duration) time.Time {
	return addMilliSeconds(t, -s)
}

// subSecond removes a second to the time.
// Positive values travels forward while negative values travels into the past.
func subSecond(t time.Time) time.Time {
	return subSeconds(t, -1)
}

// subSeconds removes seconds to the current time.
func subSeconds(t time.Time, s time.Duration) time.Time {
	return addSeconds(t, -s)
}

// subMinute removes a minute to the time.
// Positive values travels forward while negative values travels into the past.
func subMinute(t time.Time) time.Time {
	return subMinutes(t, 1)
}

// subMinutes removes minutes to the current time.
func subMinutes(t time.Time, m time.Duration) time.Time {
	return addMinutes(t, -m)
}

// subHour removes an hour to the time.
// Positive values travels forward while negative values travels into the past.
func subHour(t time.Time) time.Time {
	return subHours(t, 1)
}

// subHours removes hours to the current time.
func subHours(t time.Time, h int) time.Time {
	return addHours(t, -h)
}

// subDay removes a day to the time.
// Positive values travels forward while negative values travels into the past.
func subDay(t time.Time) time.Time {
	return subDays(t, 1)
}

// subDays removes days to the current time.
func subDays(t time.Time, d int) time.Time {
	return addDays(t, -d)
}

// subWeek removes a week to the time.
// Positive values travels forward while negative values travels into the past.
func subWeek(t time.Time) time.Time {
	return subWeeks(t, 1)
}

// subWeeks removes weeks to the current time.
func subWeeks(t time.Time, w int) time.Time {
	return addWeeks(t, -w)
}

// subMonth removes a month to the time.
// Positive values travels forward while negative values travels into the past.
func subMonth(t time.Time) time.Time {
	return subMonths(t, 1)
}

// addMonths removes months to the current time.
func subMonths(t time.Time, m int) time.Time {
	return addMonths(t, -m)
}

// subQuarter removes a quarter to the time.
// Positive values travels forward while negative values travels into the past.
func subQuarter(t time.Time) time.Time {
	return subQuarters(t, 1)
}

// subQuarters removes quarters to the current time.
func subQuarters(t time.Time, q int) time.Time {
	return addQuarters(t, -q)
}

// subYear removes a year to the time.
// Positive values travels forward while negative values travels into the past.
func subYear(t time.Time) time.Time {
	return subYears(t, 1)
}

// subYears removes years to the current time.
func subYears(t time.Time, y int) time.Time {
	return addYears(t, -y)
}

// subCentury removes a century to the time.
// Positive values travels forward while negative values travels into the past.
func subCentury(t time.Time) time.Time {
	return subCenturies(t, 1)
}

// subCenturies removes centuries to the current time.
func subCenturies(t time.Time, c int) time.Time {
	return addCenturies(t, -c)
}

// -----------------------------------------------------------------------------

// isYesterday determines if the current time is yesterday.
func isYesterday(t time.Time) bool {
	n := addDay(Now())
	return isSameDay(n)
}

// isTomorrow determines if the current time is tomorrow.
func isTomorrow(t time.Time) bool {
	n := addDay(Now())
	return isSameDay(n)
}

// isSunday checks if this day is a Sunday.
func isSunday(t time.Time) bool {
	return t.Weekday() == time.Sunday
}

// isMonday checks if this day is a Monday.
func isMonday(t time.Time) bool {
	return t.Weekday() == time.Monday
}

// isTuesday checks if this day is a Tuesday.
func isTuesday(t time.Time) bool {
	return t.Weekday() == time.Tuesday
}

// isWednesday checks if this day is a Wednesday.
func isWednesday(t time.Time) bool {
	return t.Weekday() == time.Wednesday
}

// isThursday checks if this day is a Thursday.
func isThursday(t time.Time) bool {
	return t.Weekday() == time.Thursday
}

// isFriday checks if this day is a Friday.
func isFriday(t time.Time) bool {
	return t.Weekday() == time.Friday
}

// isSaturday checks if this day is a Saturday.
func isSaturday(t time.Time) bool {
	return t.Weekday() == time.Saturday
}

// weekOfMonth returns the week of the month.
func weekOfMonth(t time.Time) int {
	w := math.Ceil(float64(t.Day() / daysPerWeek))
	return int(w + 1)
}

// weekOfYear returns the week of the current year (alias for time.ISOWeek).
func weekOfYear(t time.Time) (int, int) {
	return t.ISOWeek()
}

// isLongYear determines if the instance is a long year.
func isLongYear(t time.Time) bool {
	n := create(t.Year(), time.December, 31, 0, 0, 0, 0, t.Location())
	_, w := weekOfYear(n)
	return w == weeksPerLongYear
}

// isLastWeek returns true is the date is within last week.
func isLastWeek(t time.Time) bool {
	secondsInWeek := float64(secondsInWeek)
	diff := Now().Sub(t)
	if diff.Seconds() > 0 && diff.Seconds() < secondsInWeek {
		return true
	}
	return false
}

// isLastMonth returns true is the date is within last month.
func isLastMonth(t time.Time) bool {
	now := Now()
	monthDiff := now.Month() - t.Month()
	if absValue(true, int64(monthDiff)) != 1 {
		return false
	}
	if now.UnixNano() > t.UnixNano() && monthDiff == 1 {
		return true
	}
	return false
}

// previousMonthLastDay returns the last day of the previous month.
func previousMonthLastDay(t time.Time) time.Time {
	return t.AddDate(0, 0, -t.Day())
}

func diffInYears(t time.Time) {

}

// Comparision

// Eq, EqualTo, Ne, NotEqualTo, Gt, GreaterThan, Gte, GreaterThanOrEqualTo, Lt, LessThan, Lte, LessThanOrEqualTo, Between, Closest, Farthest,

// DiffInMonths
// DiffDurationInString
// DiffInWeeks
// DiffInDays
// DiffInNights
// DiffInSeconds
// DiffInMinutes
// DiffInHours
// SecondsSinceMidnight
// SecondsUntilEndOfDay
// swap
// StartOfDay returns the time at 00:00:00 of the same day
// EndOfDay returns the time at 23:59:59 of the same day
// StartOfMonth
// EndOfMonth
// StartOfQuarter
// EndOfQuarter
// StartOfYear
// EndOfYear
// StartOfDecade
// EndOfDecade
// StartOfCentury
// EndOfCentury
// StartOfWeek
// EndOfWeek
// Next
// NextWeekday
// PreviousWeekday
// NextWeekendDay
// PreviousWeekendDay
// Previous
// FirstOfMonth
// LastOfMonth
// LastDayOfMonth
// FirstDayOfMonth
// NthOfMonth
// FirstOfQuarter
// LastOfQuarter
// NthOfQuarter
// FirstOfYear
// LastOfYear
// NthOfYear
// Average

// set everything above. SetDay, SetHour, SetMinute, SetSecond

// isLeapYear determines if current current time is a leap year.
func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// daysInYear returns the number of days in the year.
func daysInYear(year int) int {
	if isLeap(year) {
		return daysInLeapYear
	}
	return daysInNormalYear
}

// endOfMonth returns the date at the end of the month and time at 23:59:59
func endOfMonth(year int, month time.Month) time.Time {
	return create(year, month+1, 0, 23, 59, 59, maxNSecs, time.Local)
}

// copy should return similar time instance as provided one.
func copy(t time.Time) time.Time {
	return create(
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(),
		t.Nanosecond(), t.Location(),
	)
}

// quarter should return the current quarter.
func quarter(month time.Month) int {
	switch {
	case month < 4:
		return 1
	case month >= 4 && month < 7:
		return 2
	case month >= 7 && month < 10:
		return 3
	}
	return 4
}

// loadLocation should return location instance based on given name.
func loadLocation(name string) (time.Location, error) {
	var (
		loc *time.Location
		err error
	)

	name = strings.TrimSpace(name)
	if name == "" {
		loc, err = time.LoadLocation(time.UTC.String())
	} else {
		loc, err = time.LoadLocation(name)
	}
	if err != nil {
		return time.Location{}, err
	}

	return *loc, nil
}

// inLocation should return given time in given location with error (if any).
func inLocation(t time.Time, l string) (time.Time, error) {
	loc, err := loadLocation(l)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(&loc), nil
}

// lastDayOfMonth should return last numeric day of the given month and year.
//
// Example:
//   - timex.LastDayOfMonth(2018, 2) should return 28.
func lastDayOfMonth(year int, month time.Month) int {
	// Special case `February` month.
	if month == time.February {
		if isLeap(year) {
			return 29
		}
		return 28
	}

	if month <= 7 {
		month++
	}

	if month&0x0001 == 0 {
		return 31
	}

	return 30
}

// -----------------------------------------------------------------------------

// Now returns the current local time.
func Now() time.Time {
	return now()
}

// EndOfMonth should return time value at end of the given month and year.
func EndOfMonth(year int, month time.Month) time.Time {
	return endOfMonth(year, month)
}

// Copy should return similar time instance as provided one.
func Copy(t time.Time) time.Time {
	return copy(t)
}

// Quarter should return the current quarter.
func Quarter(month time.Month) int {
	return quarter(month)
}

// LoadLocation should return location instance based on given name.
func LoadLocation(name string) (time.Location, error) {
	return loadLocation(name)
}

// IsPast determines if the current time is in the past, i.e. less (before)
// than Now().
func IsPast(t time.Time) bool {
	return isPast(t)
}

// IsFuture determines if the current time is in the future, i.e. greater (after)
// than Now().
func IsFuture(t time.Time) bool {
	return isFuture(t)
}

// IsLeap determines if current current time is a leap year.
func IsLeap(year int) bool {
	return isLeap(year)
}

// Create returns returns a time.Time from a specific date and time.
// If the location is invalid, it returns an error instead.
func Create(y int, mon time.Month, d, h, m, s, ns int, location string) (time.Time, error) {
	l, err := loadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	return create(y, mon, d, h, m, s, ns, &l), nil
}

// CreateFromDate returns a time.Time from a date.
// The time portion is set to time.Now().
// If the location is invalid, it returns an error instead.
func CreateFromDate(y int, mon time.Month, d int, location string) (time.Time, error) {
	now := now()
	h, m, s := now.Clock()
	return Create(y, mon, d, h, m, s, now.Nanosecond(), location)
}

// UnixMilli returns the number of milliseconds elapsed since January 1, 1970
// UTC.
func UnixMilli() int64 {
	return now().UnixNano() / time.Millisecond.Nanoseconds()
}

// NowIn should return time.Time in provided location.
func NowIn(loc *time.Location) time.Time {
	return nowIn(loc)
}

// NowInLocation returns a current time in given location.
// The location is in IANA Time Zone database, such as "America/New_York".
func NowInLocation(l string) (time.Time, error) {
	return InLocation(now(), l)
}

// InLocation should return given time in given location with error (if any).
func InLocation(t time.Time, l string) (time.Time, error) {
	return inLocation(t, l)
}

// InLocationFormat should return given time in given location with error (if
// any). With given format layout.
func InLocationFormat(t time.Time, locStr, layout string) (string, error) {
	t, err := inLocation(t, locStr)
	if err != nil {
		return "", err
	}
	return t.Format(layout), nil
}

// DaysBetween returns the number of whole days between the start date and the
// end date.
func DaysBetween(fromDate, toDate time.Time) int {
	return int(toDate.Sub(fromDate) / (24 * time.Hour))
}

// LastDayOfMonth should return last numeric day of the given month and year.
//
// Example:
//   - timex.LastDayOfMonth(2018, 2) should return 28.
func LastDayOfMonth(year int, month time.Month) int {
	return lastDayOfMonth(year, month)
}

// HoursDiff should return hours' difference between given two dates.
func HoursDiff(t, u time.Time) float64 {
	return t.Sub(u).Hours()
}

// MinutesDiff should return minutes' difference between given two dates.
func MinutesDiff(t, u time.Time) float64 {
	return t.Sub(u).Minutes()
}

// SecondsDiff should return seconds' difference between given two dates.
func SecondsDiff(t, u time.Time) float64 {
	return t.Sub(u).Seconds()
}

// NanosecondsDiff should return nanoseconds' difference between given two dates.
func NanosecondsDiff(t, u time.Time) int64 {
	return t.Sub(u).Nanoseconds()
}

// GetOneDayBeginOfTime returns the begin of the time t.
func GetOneDayBeginOfTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetOneDayEndOfTime returns the end of the time t.
func GetOneDayEndOfTime(t time.Time) time.Time {
	return GetOneDayBeginOfTime(t).Add(24 * time.Hour).Add(-1 * time.Nanosecond)
}

// TimeBeginningOfWeek return the begin of the week of time t.
// sundayFirst is used to set week day. As in some countries uses `Monday` as
// the first day of the week.
func TimeBeginningOfWeek(t time.Time, sundayFirst bool) time.Time {
	weekday := int(t.Weekday())
	if !sundayFirst {
		if weekday == 0 {
			weekday = 7
		}
		weekday = weekday - 1
	}

	d := time.Duration(-weekday) * 24 * time.Hour
	t = t.Add(d)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// TimeEndOfWeek return the end of the week of time t.
// sundayFirst is used to set week day. As in some countries uses `Monday` as
// the first day of the week.
func TimeEndOfWeek(t time.Time, sundayFirst bool) time.Time {
	return TimeBeginningOfWeek(t, sundayFirst).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// TimeBeginningOfMonth return the begin of the month of time t.
func TimeBeginningOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// TimeEndOfMonth return the end of the month of time t.
func TimeEndOfMonth(t time.Time) time.Time {
	return TimeBeginningOfMonth(t).AddDate(0, 1, -1)
}

// TimeSubDaysOfTwoDays should return the days bewteen time.Time d1 and
// time.Time d2.
func TimeSubDaysOfTwoDays(d1 time.Time, d2 time.Time) int64 {
	ds1 := GetOneDayBeginOfTime(d1)
	ds2 := GetOneDayBeginOfTime(d2)
	return int64(ds1.Sub(ds2).Hours() / 24)
}

// Format ----------------------------------------------------------------------

// ANSICFormat should return given time in ANSIC format.
//
// Example: "Mon Jan _2 15:04:05 2006"
func ANSICFormat(t time.Time) string {
	return t.Format(time.ANSIC)
}

// UnixDateFormat should return given time in UnixDate format.
//
// Example: "Mon Jan _2 15:04:05 MST 2006"
func UnixDateFormat(t time.Time) string {
	return t.Format(time.UnixDate)
}

// RubyDateFormat should return given time in RubyDate format.
//
// Example: "Mon Jan 02 15:04:05 -0700 2006"
func RubyDateFormat(t time.Time) string {
	return t.Format(time.RubyDate)
}

// RFC822Format should return given time in RFC3339 format.
//
// Example: "02 Jan 06 15:04 MST"
func RFC822Format(t time.Time) string {
	return t.Format(time.RFC3339)
}

// RFC822ZFormat should return given time in RFC822Z (RFC822 with numeric
// zone) format.
//
// Example: "02 Jan 06 15:04 -0700"
func RFC822ZFormat(t time.Time) string {
	return t.Format(time.RFC822Z)
}

// RFC850Format should return given time in RFC850 format.
//
// Example: "Monday, 02-Jan-06 15:04:05 MST"
func RFC850Format(t time.Time) string {
	return t.Format(time.RFC850)
}

// RFC1123Format should return given time in RFC1123 format.
//
// Example: "Mon, 02 Jan 2006 15:04:05 MST"
func RFC1123Format(t time.Time) string {
	return t.Format(time.RFC1123)
}

// RFC1123ZFormat should return given time in RFC1123Z (RFC1123 with numeric
// zone) format.
//
// Example: "Mon, 02 Jan 2006 15:04:05 -0700"
func RFC1123ZFormat(t time.Time) string {
	return t.Format(time.RFC1123Z)
}

// RFC3339Format should return given time in RFC3339 format.
//
// Example: "2006-01-02T15:04:05Z07:00"
func RFC3339Format(t time.Time) string {
	return t.Format(time.RFC3339)
}

// RFC3339NanoFormat should return given time in RFC3339Nano format.
//
// Example: "2006-01-02T15:04:05.999999999Z07:00"
func RFC3339NanoFormat(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

// KitchenFormat should return given time in Kitchen format.
//
// Example: "3:04PM"
func KitchenFormat(t time.Time) string {
	return t.Format(time.Kitchen)
}

// Handy time stamps.

// StampFormat should return given time in Stamp format.
//
// Example: "Jan _2 15:04:05"
func StampFormat(t time.Time) string {
	return t.Format(time.Stamp)
}

// StampMilliFormat should return given time in StampMilli format.
//
// Example: "Jan _2 15:04:05.000"
func StampMilliFormat(t time.Time) string {
	return t.Format(time.StampMilli)
}

// StampMicroFormat should return given time in StampMicro format.
//
// Example: "Jan _2 15:04:05.000000"
func StampMicroFormat(t time.Time) string {
	return t.Format(time.StampMicro)
}

// StampNanoFormat should return given time in StampNano format.
//
// Example: "Jan _2 15:04:05.000000000"
func StampNanoFormat(t time.Time) string {
	return t.Format(time.StampNano)
}

// -----------------------------------------------------------------------------

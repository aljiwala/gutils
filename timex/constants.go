package timex

// Represents the number of elements in a given period.
const (
	maxNSecs          = 999999999
	secondsPerMinute  = 60
	minutesPerHour    = 60
	secondsInHour     = secondsPerMinute * minutesPerHour
	hoursPerDay       = 24
	daysPerWeek       = 7
	monthsPerQuarter  = 3
	monthsPerYear     = 12
	yearsPerCenturies = 100
	yearsPerDecade    = 10
	weeksPerLongYear  = 53
	daysInLeapYear    = 366
	daysInNormalYear  = 365
	secondsInWeek     = secondsPerMinute * minutesPerHour * hoursPerDay * daysPerWeek
	secondsInMonth    = 2678400

	SecondsPerMinute  = secondsPerMinute
	MinutesPerHour    = minutesPerHour
	HoursPerDay       = hoursPerDay
	DaysPerWeek       = daysPerWeek
	MonthsPerQuarter  = monthsPerQuarter
	MonthsPerYear     = monthsPerYear
	YearsPerCenturies = yearsPerCenturies
	YearsPerDecade    = yearsPerDecade
	WeeksPerLongYear  = weeksPerLongYear
	DaysInLeapYear    = daysInLeapYear
	DaysInNormalYear  = daysInNormalYear
	SecondsInWeek     = secondsInWeek
	SecondsInMonth    = secondsInMonth
)

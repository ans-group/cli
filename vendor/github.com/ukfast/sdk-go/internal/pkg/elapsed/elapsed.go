package elapsed

import "time"

func divmod(numerator, denominator time.Duration) (quotient, remainder time.Duration) {
	return numerator / denominator, numerator % denominator
}

// ParseDuration takes provided duration d and returns the duration in years, months, days, hours, minutes, seconds and nanoseconds.
// It uses 30 days as the length of a month in days
func ParseDuration(d time.Duration) (years, months, days, hours, minutes, seconds, nanoseconds int) {
	yearsVal, d := divmod(d, time.Hour*24*365)
	monthsVal, d := divmod(d, time.Hour*24*30)
	daysVal, d := divmod(d, time.Hour*24)
	hoursVal, d := divmod(d, time.Hour)
	minutesVal, d := divmod(d, time.Minute)
	secondsVal, d := divmod(d, time.Second)
	nanosecondsVal, d := divmod(d, time.Nanosecond)

	if nanosecondsVal < 0 {
		nanosecondsVal += 1e9
		secondsVal--
	}
	if secondsVal < 0 {
		secondsVal += 60
		minutesVal--
	}
	if minutesVal < 0 {
		minutesVal += 60
		hoursVal--
	}
	if hoursVal < 0 {
		hoursVal += 24
		daysVal--
	}
	if daysVal < 0 {
		daysVal += 30
		monthsVal--
	}
	if monthsVal < 0 {
		monthsVal += 12
		yearsVal--
	}

	return int(yearsVal), int(monthsVal), int(daysVal), int(hoursVal), int(minutesVal), int(secondsVal), int(nanosecondsVal)
}

// NewDuration takes input values and returns the corresponding duration
func NewDuration(years, months, days, hours, minutes, seconds, nanoseconds int) time.Duration {
	day := time.Hour * 24
	d := time.Duration(years) * day * 365
	d = d + time.Duration(months)*day*30
	d = d + time.Duration(days)*day
	d = d + time.Duration(hours)*time.Hour
	d = d + time.Duration(minutes)*time.Minute
	d = d + time.Duration(seconds)*time.Second
	return d + time.Duration(nanoseconds)*time.Nanosecond
}

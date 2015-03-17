package pragmash

import (
	"errors"
	"math"
	"time"
)

// StdTime implements time-related commands.
type StdTime struct{}

// Sleep stops execution for a given number of seconds which may be fractional.
func (_ StdTime) Sleep(f float64) {
	nanos := math.Floor(f * 1000000000)
	time.Sleep(time.Nanosecond * time.Duration(nanos))
}

// Time returns the UNIX epoch time in seconds which may be fractional.
func (_ StdTime) Time() float64 {
	return float64(time.Now().UnixNano()) / 1000000000
}

// TimeDay returns the day of the month of a timestamp (starting at 1).
func (_ StdTime) TimeDay(t float64, location ...string) (int, error) {
	x, err := timestampToTime(t, location)
	if err != nil {
		return 0, err
	}
	return x.Day(), nil
}

// TimeHour returns the hour (between 0 and 60) of a timestamp.
func (_ StdTime) TimeHour(t float64, location ...string) (int, error) {
	x, err := timestampToTime(t, location)
	if err != nil {
		return 0, err
	}
	return x.Hour(), nil
}

// TimeMinute returns the minute (between 0 and 60) of a timestamp.
func (_ StdTime) TimeMinute(t float64, location ...string) (int, error) {
	x, err := timestampToTime(t, location)
	if err != nil {
		return 0, err
	}
	return x.Minute(), nil
}

// TimeMonth returns the month (1 to 12 inclusive) of a timestamp.
func (_ StdTime) TimeMonth(t float64, location ...string) (int, error) {
	x, err := timestampToTime(t, location)
	if err != nil {
		return 0, err
	}
	return int(x.Month()), nil
}

// TimeSecond returns the second (between 0 and 60) of a timestamp.
func (_ StdTime) TimeSecond(t float64, location ...string) (int, error) {
	x, err := timestampToTime(t, location)
	if err != nil {
		return 0, err
	}
	return x.Second(), nil
}

// TimeYear returns the year of a timestamp.
func (_ StdTime) TimeYear(t float64, location ...string) (int, error) {
	x, err := timestampToTime(t, location)
	if err != nil {
		return 0, err
	}
	return x.Year(), nil
}

// timestampToTime returns a time.Time for a timestamp.
func timestampToTime(t float64, location []string) (time.Time, error) {
	if len(location) > 1 {
		return time.Time{}, errors.New("expected 1 or 2 arguments")
	}
	
	// Extract nanoseconds and seconds from fractional timestamp.
	var nanos int64
	var seconds int64
	if t > 0 {
		nanos = int64(1000000000 * (t - math.Floor(t)))
		seconds = int64(t)
	} else {
		nanos = -int64(1000000000 * (-t - math.Floor(-t)))
		seconds = -int64(-t)
	}
	
	res := time.Unix(seconds, nanos)
	
	if len(location) == 0 {
		return res, nil
	}

	// Use the specified location
	loc, err := time.LoadLocation(location[0])
	if err != nil {
		return time.Time{}, err
	}
	return res.In(loc), nil
}

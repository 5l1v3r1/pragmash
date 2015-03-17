package pragmash

import (
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
func (_ StdTime) TimeDay(t float64) int {
	return timestampToTime(t).Day()
}

// TimeHour returns the hour (between 0 and 60) of a timestamp.
func (_ StdTime) TimeHour(t float64) int {
	return timestampToTime(t).Hour()
}

// TimeMinute returns the minute (between 0 and 60) of a timestamp.
func (_ StdTime) TimeMinute(t float64) int {
	return timestampToTime(t).Minute()
}

// TimeMonth returns the month (1 to 12 inclusive) of a timestamp.
func (_ StdTime) TimeMonth(t float64) int {
	return int(timestampToTime(t).Month())
}

// TimeSecond returns the second (between 0 and 60) of a timestamp.
func (_ StdTime) TimeSecond(t float64) int {
	return timestampToTime(t).Second()
}

// TimeYear returns the year of a timestamp.
func (_ StdTime) TimeYear(t float64) int {
	return timestampToTime(t).Year()
}

// timestampToTime returns a time.Time for a timestamp.
func timestampToTime(t float64) time.Time {
	var nanos int64
	var seconds int64
	if t > 0 {
		nanos = int64(1000000000 * (t - math.Floor(t)))
		seconds = int64(t)
	} else {
		nanos = -int64(1000000000 * (-t - math.Floor(-t)))
		seconds = -int64(-t)
	}
	return time.Unix(seconds, nanos)
}

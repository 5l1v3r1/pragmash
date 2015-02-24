package pragmash

import (
	"math"
	"time"
)

// StdTime implements time-related commands.
type StdTime struct{}

// Sleep stops execution for a given number of seconds which may be fractional.
func (_ StdTime) Sleep(n Number) {
	nanos := math.Floor(n.Float() * 1000000000)
	time.Sleep(time.Nanosecond * time.Duration(nanos))
}

// Time returns the UNIX epoch time in seconds which may be fractional.
func (_ StdTime) Time() Value {
	number := float64(time.Now().UnixNano()) / 1000000000
	return NewNumberFloat(number)
}

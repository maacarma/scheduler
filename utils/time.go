package utils

import "time"

type Unix int64

// CurrentUTCUnix returns the current UTC unix time.
func CurrentUTCUnix() Unix {
	return Unix(time.Now().UTC().Unix())
}

// Diff returns the difference between a and b in time duration
// if (reverse = false) then (a - b)
// if (reverse = true) the (b - a)
func (a Unix) Sub(b Unix, reverse bool) time.Duration {
	if reverse {
		return time.Duration(b-a) * time.Second
	}

	return time.Duration(a-b) * time.Second
}

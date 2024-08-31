package utils

import "time"

type Unix int64

// CurrentUTCUnix returns the current UTC unix time.
func CurrentUTCUnix() Unix {
	return Unix(time.Now().UTC().Unix())
}

// Diff returns the difference between u and b in time duration
// if (reverse = false) then (u - b)
// if (reverse = true) the (b - u)
func (u Unix) Sub(b Unix, reverse bool) time.Duration {
	if reverse {
		return time.Duration(b-u) * time.Second
	}

	return time.Duration(u-b) * time.Second
}

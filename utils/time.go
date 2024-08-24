package utils

import "time"

type Unix int64

// CurrentUTCUnix returns the current UTC unix time.
func CurrentUTCUnix() Unix {
	return Unix(time.Now().UTC().Unix())
}

// Diff returns the difference between the given time and the current UTC unix time.
// If elapsed is true,  (current time - given time).
// If elapsed is false, (given time - current time).
func (u Unix) Diff(elapsed bool) time.Duration {
	curTime := CurrentUTCUnix()
	if elapsed {
		return time.Duration(curTime-u) * time.Second
	}

	return time.Duration(u-curTime) * time.Second
}

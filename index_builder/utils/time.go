package utils

import "time"

func MakeTimeRange(start, end time.Time, duration time.Duration) (times []time.Time) {
	times = make([]time.Time, 0)

	for t := start; t.Before(end); t = t.Add(duration) {
		times = append(times, t)
	}

	return
}

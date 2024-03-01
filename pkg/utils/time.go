package utils

import "time"

func SecondToNano(t int64) int64 {
	return t * 1e9
}

func Second2Time(t int64) time.Time {
	return time.Unix(t, 0)
}

func Time2Second(t time.Time) int64 {
	return t.Unix()
}

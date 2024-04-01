package utils

import "time"

const (
	Layout = time.RFC3339
)

func SecondToNano(t int64) int64 {
	return t * 1e9
}

func NanoToSecond(t int64) int64 {
	res := t / 1e9
	if res <= 0 {
		res = 1
	}
	return res
}

func Second2Time(t int64) time.Time {
	return time.Unix(t, 0)
}

func Str2Time(str string) time.Time {
	t, err := time.Parse(Layout, str)
	if err != nil {
		res, _ := time.Parse(Layout, "0001-01-01 00:00:00 +0000 UTC")
		return res
	}
	return t
}

func Time2Str(t time.Time) string {
	return t.Format(Layout)
}

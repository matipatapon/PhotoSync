package helper

import "time"

type ITimeHelper interface {
	TimeNow() int64
	TimeIn(seconds int64) int64
}

type TimeHelper struct {
}

func (TimeHelper) TimeNow() int64 {
	return time.Now().Unix()
}

func (th *TimeHelper) TimeIn(seconds int64) int64 {
	return th.TimeNow() + seconds
}

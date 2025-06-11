package mock

import (
	"fmt"
	"testing"
)

type TimeHelperMock struct {
	expectedTimeNowTime  []int64
	expectedTimeInSecond []int64
	expectedTimeInTime   []int64
	t                    *testing.T
}

func NewTimeHelperMock(t *testing.T) TimeHelperMock {
	return TimeHelperMock{[]int64{}, []int64{}, []int64{}, t}
}

func (thm *TimeHelperMock) TimeNow() int64 {
	if len(thm.expectedTimeNowTime) <= 0 {
		fmt.Print("Unexpected TimeNow()!")
		thm.t.FailNow()
	}
	time := thm.expectedTimeNowTime[len(thm.expectedTimeNowTime)-1]
	thm.expectedTimeNowTime = thm.expectedTimeNowTime[:len(thm.expectedTimeNowTime)-1]
	return time
}

func (thm *TimeHelperMock) TimeIn(seconds int64) int64 {
	if len(thm.expectedTimeInSecond) <= 0 {
		fmt.Print("Unexpected TimeIn()!")
		thm.t.FailNow()
	}
	expectedSeconds := thm.expectedTimeInSecond[len(thm.expectedTimeInSecond)-1]
	thm.expectedTimeInSecond = thm.expectedTimeInSecond[:len(thm.expectedTimeInSecond)-1]
	if expectedSeconds != seconds {
		fmt.Print("Unexpected parameter!")
		thm.t.FailNow()
	}
	time := thm.expectedTimeInTime[len(thm.expectedTimeInTime)-1]
	thm.expectedTimeInTime = thm.expectedTimeInTime[:len(thm.expectedTimeInTime)-1]
	return time
}

func (thm *TimeHelperMock) ExpectTimeNow(time int64) {
	thm.expectedTimeNowTime = append(thm.expectedTimeNowTime, time)
}

func (thm *TimeHelperMock) ExpectTimeIn(seconds int64, time int64) {
	thm.expectedTimeInSecond = append(thm.expectedTimeInSecond, seconds)
	thm.expectedTimeInTime = append(thm.expectedTimeInTime, time)
}

func (thm *TimeHelperMock) AssertAllExpectionsSatisfied() {
	if len(thm.expectedTimeInSecond) != 0 || len(thm.expectedTimeNowTime) != 0 {
		fmt.Print("Not all expections satisfied!")
		thm.t.FailNow()
	}
}

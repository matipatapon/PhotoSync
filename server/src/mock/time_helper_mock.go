package mock

import (
	"fmt"
	"photosync/src/helper"
	"testing"
)

type TimeHelperMock struct {
	expectedTimeNowTime  helper.List[int64]
	expectedTimeInSecond helper.List[int64]
	expectedTimeInTime   helper.List[int64]
	t                    *testing.T
}

func NewTimeHelperMock(t *testing.T) TimeHelperMock {
	return TimeHelperMock{t: t}
}

func (thm *TimeHelperMock) TimeNow() int64 {
	if thm.expectedTimeNowTime.Length() <= 0 {
		fmt.Print("Unexpected TimeNow()!")
		thm.t.FailNow()
	}
	return thm.expectedTimeNowTime.PopFirst()
}

func (thm *TimeHelperMock) TimeIn(seconds int64) int64 {
	if thm.expectedTimeInSecond.Length() <= 0 {
		fmt.Print("Unexpected TimeIn()!")
		thm.t.FailNow()
	}

	expectedSeconds := thm.expectedTimeInSecond.PopFirst()
	if expectedSeconds != seconds {
		fmt.Print("Unexpected parameter!")
		thm.t.FailNow()
	}

	return thm.expectedTimeInTime.PopFirst()
}

func (thm *TimeHelperMock) ExpectTimeNow(time int64) {
	thm.expectedTimeNowTime.Append(time)
}

func (thm *TimeHelperMock) ExpectTimeIn(seconds int64, time int64) {
	thm.expectedTimeInSecond.Append(seconds)
	thm.expectedTimeInTime.Append(time)
}

func (thm *TimeHelperMock) AssertAllExpectionsSatisfied() {
	if thm.expectedTimeInSecond.Length() != 0 || thm.expectedTimeNowTime.Length() != 0 {
		fmt.Print("Not all expections satisfied!")
		thm.t.FailNow()
	}
}

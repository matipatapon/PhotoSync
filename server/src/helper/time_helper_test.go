package helper_test

import (
	"photosync/src/helper"
	"testing"
	"time"
)

func TestTimeHelperShouldReturnCurrentTime(t *testing.T) {
	sut := helper.TimeHelper{}

	time1 := sut.TimeNow()
	time2 := time.Now().Unix()
	diff := time1 - time2
	if diff < 0 {
		diff = -diff
	}

	if diff > 5 {
		t.FailNow()
	}
}

func TestTimeHelperShouldReturnFutureTime(t *testing.T) {
	sut := helper.TimeHelper{}

	time1 := sut.TimeIn(500)
	time2 := time.Now().Unix() + 500
	diff := time1 - time2
	if diff < 0 {
		diff = -diff
	}

	if diff > 5 {
		t.FailNow()
	}
}

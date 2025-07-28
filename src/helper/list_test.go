package helper_test

import (
	"photosync/src/helper"
	"testing"
)

func TestListShouldAppendAndPop(t *testing.T) {
	list := helper.List[int]{}
	if list.Length() != 0 {
		t.FailNow()
	}

	list.Append(10)
	if list.Length() != 1 {
		t.FailNow()
	}

	list.Append(12)
	if list.Length() != 2 {
		t.FailNow()
	}

	if list.PopFirst() != 10 || list.Length() != 1 {
		t.FailNow()
	}

	if list.PopFirst() != 12 || list.Length() != 0 {
		t.FailNow()
	}
}

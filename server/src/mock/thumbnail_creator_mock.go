package mock

import (
	"fmt"
	"photosync/src/helper"
	"photosync/src/metadata"
	"reflect"
	"testing"
)

type ThumbnailCreatorMock struct {
	expectedFiles     helper.List[[]byte]
	expectedMIMEType  helper.List[metadata.MIMEType]
	expectedThumbnail helper.List[[]byte]
	expectedErrors    helper.List[error]
	t                 *testing.T
}

func NewThumbnailCreatorMock(t *testing.T) ThumbnailCreatorMock {
	return ThumbnailCreatorMock{t: t}
}

func (tcm *ThumbnailCreatorMock) ExpectCreate(
	file []byte,
	mimeType metadata.MIMEType,
	thumbnail []byte,
	err error,
) {
	tcm.expectedFiles.Append(file)
	tcm.expectedMIMEType.Append(mimeType)
	tcm.expectedThumbnail.Append(thumbnail)
	tcm.expectedErrors.Append(err)
}

func (tcm *ThumbnailCreatorMock) Create(file []byte, mimeType metadata.MIMEType) ([]byte, error) {
	if tcm.expectedFiles.Length() == 0 {
		fmt.Println("Unexpected Create")
		tcm.t.FailNow()
	}

	expectedFile := tcm.expectedFiles.PopFirst()
	if !reflect.DeepEqual(expectedFile, file) {
		fmt.Println("Unexpected file")
		tcm.t.FailNow()
	}

	expectedMIMEType := tcm.expectedMIMEType.PopFirst()
	if expectedMIMEType != mimeType {
		fmt.Println("Unexpected MIMEType")
		tcm.t.FailNow()
	}

	return tcm.expectedThumbnail.PopFirst(), tcm.expectedErrors.PopFirst()
}

func (tcm *ThumbnailCreatorMock) AssertAllExpectionsSatisfied() {
	if tcm.expectedFiles.Length() != 0 {
		fmt.Println("Not all expections satisfied!")
		tcm.t.FailNow()
	}
}

package mock

import (
	"fmt"
	"photosync/src/helper"
	"reflect"
	"testing"
)

type RawMetadataExtractorMock struct {
	expectedFiles    helper.List[[]byte]
	expectedMetadata helper.List[map[string]any]
	expectedErrors   helper.List[error]
	t                *testing.T
}

func NewRawMetadataExtractorMock(t *testing.T) RawMetadataExtractorMock {
	return RawMetadataExtractorMock{t: t}
}

func (rmem *RawMetadataExtractorMock) ExpectExtract(file []byte, metadata map[string]any, err error) {
	rmem.expectedFiles.Append(file)
	rmem.expectedMetadata.Append(metadata)
	rmem.expectedErrors.Append(err)
}

func (rmem *RawMetadataExtractorMock) Extract(file []byte) (map[string]any, error) {
	if rmem.expectedFiles.Length() == 0 {
		fmt.Print("Unexpected extract!")
		rmem.t.FailNow()
	}

	expectedFile := rmem.expectedFiles.PopFirst()
	if !reflect.DeepEqual(file, expectedFile) {
		fmt.Print("Unexpected file!")
		rmem.t.FailNow()
	}

	return rmem.expectedMetadata.PopFirst(), rmem.expectedErrors.PopFirst()
}

func (rmem *RawMetadataExtractorMock) AssertAllExpectionsSatisfied() {
	if rmem.expectedFiles.Length() != 0 {
		fmt.Print("Not all expections satisfied!")
		rmem.t.FailNow()
	}
}

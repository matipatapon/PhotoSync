package mock

import (
	"fmt"
	"reflect"
	"testing"
)

type RawMetadataExtractorMock struct {
	expectedFiles    [][]byte
	expectedMetadata []map[string]any
	expectedErrors   []error
	t                *testing.T
}

func NewRawMetadataExtractorMock(t *testing.T) RawMetadataExtractorMock {
	return RawMetadataExtractorMock{[][]byte{}, []map[string]any{}, []error{}, t}
}

func (rmem *RawMetadataExtractorMock) ExpectExtract(file []byte, metadata map[string]any, err error) {
	rmem.expectedFiles = append(rmem.expectedFiles, file)
	rmem.expectedMetadata = append(rmem.expectedMetadata, metadata)
	rmem.expectedErrors = append(rmem.expectedErrors, err)
}

func (rmem *RawMetadataExtractorMock) Extract(file []byte) (map[string]any, error) {
	if len(rmem.expectedFiles) == 0 {
		fmt.Print("Unexpected extract!")
		rmem.t.FailNow()
	}

	expectedFile := rmem.expectedFiles[len(rmem.expectedFiles)-1]
	rmem.expectedFiles = rmem.expectedFiles[:len(rmem.expectedFiles)-1]
	if !reflect.DeepEqual(file, expectedFile) {
		fmt.Print("Unexpected file!")
		rmem.t.FailNow()
	}

	expectedMetadata := rmem.expectedMetadata[len(rmem.expectedMetadata)-1]
	rmem.expectedMetadata = rmem.expectedMetadata[:len(rmem.expectedMetadata)-1]

	expectedError := rmem.expectedErrors[len(rmem.expectedErrors)-1]
	rmem.expectedErrors = rmem.expectedErrors[:len(rmem.expectedErrors)-1]

	return expectedMetadata, expectedError
}

func (rmem *RawMetadataExtractorMock) AssertAllExpectionsSatisfied() {
	if len(rmem.expectedFiles) != 0 {
		fmt.Print("Not all expections satisfied!")
		rmem.t.FailNow()
	}
}

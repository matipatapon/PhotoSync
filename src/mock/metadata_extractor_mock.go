package mock

import (
	"fmt"
	"photosync/src/helper"
	"photosync/src/metadata"
	"reflect"
	"testing"
)

type MetadataExtractorMock struct {
	expectedFiles    helper.List[[]byte]
	expectedMetadata helper.List[metadata.Metadata]
	t                *testing.T
}

func NewMetadataExtractorMock(t *testing.T) MetadataExtractorMock {
	return MetadataExtractorMock{t: t}
}

func (mem *MetadataExtractorMock) ExpectExtract(file []byte, metadata metadata.Metadata) {
	mem.expectedFiles.Append(file)
	mem.expectedMetadata.Append(metadata)
}

func (mem *MetadataExtractorMock) Extract(file []byte) metadata.Metadata {
	if mem.expectedFiles.Length() == 0 {
		fmt.Print("Unexpected extract")
		mem.t.FailNow()
	}

	expectedFile := mem.expectedFiles.PopFirst()
	if !reflect.DeepEqual(expectedFile, file) {
		fmt.Print("Unexpected file")
		mem.t.FailNow()
	}

	return mem.expectedMetadata.PopFirst()
}

func (mem *MetadataExtractorMock) AssertAllExpectionsSatisfied() {
	if mem.expectedFiles.Length() != 0 {
		fmt.Print("Not all expections satisfied!")
		mem.t.FailNow()
	}
}

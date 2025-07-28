package metadata_test

import (
	"errors"
	"fmt"
	"photosync/src/metadata"
	"photosync/src/mock"
	"reflect"
	"testing"
)

var file []byte = []byte("FILE")

func TestShouldReturnEmptyMetadataWhenFileDoesNotHaveTags(t *testing.T) {
	meta := map[string]any{}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	if result.CreationDate != nil || result.Location != nil || result.MIMEType != metadata.UNKNOWN {
		fmt.Print("Expected everything to be nil")
		t.FailNow()
	}
}

func TestShouldExtractMetadataFromCompositeTags(t *testing.T) {
	meta := map[string]any{
		"Composite:DateTimeOriginal": "2022.03.01 21:37:00",
		"Composite:GPSPosition":      "51 6 32.29 N, 17 1 59.30 E",
		"File:MIMEType":              "image/jpeg",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if result.CreationDate == nil || !reflect.DeepEqual(expectedMetadata, *result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if result.Location == nil || !reflect.DeepEqual(expectedLocation, *result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldExtractMetadataFromExifTags(t *testing.T) {
	meta := map[string]any{
		"EXIF:DateTimeOriginal": "2022.03.01 21:37:00",
		"Composite:GPSPosition": "51 6 32.29 N, 17 1 59.30 E",
		"File:MIMEType":         "image/jpeg",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if result.CreationDate == nil || !reflect.DeepEqual(expectedMetadata, *result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if result.Location == nil || !reflect.DeepEqual(expectedLocation, *result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldExtractMetadataFromXmpTags(t *testing.T) {
	meta := map[string]any{
		"XMP:CreateDate":        "2022.03.01 21:37:00",
		"Composite:GPSPosition": "51 6 32.29 N, 17 1 59.30 E",
		"File:MIMEType":         "image/jpeg",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if result.CreationDate == nil || !reflect.DeepEqual(expectedMetadata, *result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if result.Location == nil || !reflect.DeepEqual(expectedLocation, *result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldExtractMetadataFromQuickTimeTags(t *testing.T) {
	meta := map[string]any{
		"QuickTime:CreateDate":  "2022.03.01 21:37:00",
		"Composite:GPSPosition": "51 6 32.29 N, 17 1 59.30 E",
		"File:MIMEType":         "image/jpeg",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if result.CreationDate == nil || !reflect.DeepEqual(expectedMetadata, *result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if result.Location == nil || !reflect.DeepEqual(expectedLocation, *result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldReturnEmptyDateWhenThereAreNoCreationDateTags(t *testing.T) {
	meta := map[string]any{
		"Composite:GPSPosition": "51 6 32.29 N, 17 1 59.30 E",
		"File:MIMEType":         "image/jpeg",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	if result.CreationDate != nil {
		fmt.Print("CreationDate expected to be nil")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if result.Location == nil || !reflect.DeepEqual(expectedLocation, *result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldReturnNilForGPSWhenThereIsNoGPSTag(t *testing.T) {
	meta := map[string]any{
		"QuickTime:CreateDate": "2022.03.01 21:37:00",
		"File:MIMEType":        "image/jpeg",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if result.CreationDate == nil || !reflect.DeepEqual(expectedMetadata, *result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	if result.Location != nil {
		fmt.Print("Location expected to be nil")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestMIMeTypeShouldBeUnknownWhenThereIsNoMIMeTypeTag(t *testing.T) {
	meta := map[string]any{
		"QuickTime:CreateDate":  "2022.03.01 21:37:00",
		"Composite:GPSPosition": "51 6 32.29 N, 17 1 59.30 E",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if result.CreationDate == nil || !reflect.DeepEqual(expectedMetadata, *result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if result.Location == nil || !reflect.DeepEqual(expectedLocation, *result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.UNKNOWN {
		fmt.Print("MIMEType should be unknown")
		t.FailNow()
	}
}

func TestReturnedMIMeTypeShouldBeUnknownWhenMIMeTypeFromTagIsNotRecognized(t *testing.T) {
	meta := map[string]any{
		"QuickTime:CreateDate":  "2022.03.01 21:37:00",
		"Composite:GPSPosition": "51 6 32.29 N, 17 1 59.30 E",
		"File:MIMEType":         "image/unknown",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if result.CreationDate == nil || !reflect.DeepEqual(expectedMetadata, *result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if result.Location == nil || !reflect.DeepEqual(expectedLocation, *result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.UNKNOWN {
		fmt.Print("MIMEType should be unknown")
		t.FailNow()
	}
}

func TestShouldReturnNilForCreationDateWhenDateInTagIsInvalid(t *testing.T) {
	meta := map[string]any{
		"QuickTime:CreateDate":  "wrong date",
		"Composite:GPSPosition": "51 6 32.29 N, 17 1 59.30 E",
		"File:MIMEType":         "image/jpeg",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	if result.CreationDate != nil {
		fmt.Print("Date should be nil")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if result.Location == nil || !reflect.DeepEqual(expectedLocation, *result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldSkipInvalidCreationDateTag(t *testing.T) {
	meta := map[string]any{
		"Composite:DateTimeOriginal": "wrong date",
		"QuickTime:CreateDate":       "2022.03.01 21:37:00",
		"Composite:GPSPosition":      "51 6 32.29 N, 17 1 59.30 E",
		"File:MIMEType":              "image/jpeg",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if result.CreationDate == nil || !reflect.DeepEqual(expectedMetadata, *result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if result.Location == nil || !reflect.DeepEqual(expectedLocation, *result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldReturnNilLocationWhenGPSPositionTagIsInvalid(t *testing.T) {
	meta := map[string]any{
		"Composite:DateTimeOriginal": "2022.03.01 21:37:00",
		"Composite:GPSPosition":      "invalid tag",
		"File:MIMEType":              "image/jpeg",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if result.CreationDate == nil || !reflect.DeepEqual(expectedMetadata, *result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	if result.Location != nil {
		fmt.Print("Location expected to be nil")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldReturnEmptyMetadataWhenRawExtractionFailed(t *testing.T) {
	meta := map[string]any{
		"QuickTime:CreateDate":  "2022.03.01 21:37:00",
		"Composite:GPSPosition": "51 6 32.29 N, 17 1 59.30 E",
		"File:MIMEType":         "image/jpeg",
	}
	err := errors.New("flatlined")
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, err)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	if result.CreationDate != nil {
		fmt.Print("Date expected to be nil")
		t.FailNow()
	}

	if result.Location != nil {
		fmt.Print("Location expected to be nil")
		t.FailNow()
	}

	if result.MIMEType != metadata.UNKNOWN {
		fmt.Print("MIMEType should be unknown")
		t.FailNow()
	}
}

package metadata_test

import (
	"fmt"
	"photosync/src/metadata"
	"photosync/src/mock"
	"reflect"
	"testing"
)

var file []byte = []byte("FILE")

func TestShouldReturnErrorWhenFileDoesNotHaveMetadata(t *testing.T) {
	meta := map[string]any{}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)

	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	if sut.Extract(file) != nil {
		fmt.Print("Expected nil")
		t.FailNow()
	}
}

func TestShouldExtractCompositeMetadata(t *testing.T) {
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

	if result == nil {
		fmt.Print("Nil received")
		t.FailNow()
	}

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if !reflect.DeepEqual(expectedMetadata, result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if !reflect.DeepEqual(expectedLocation, result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldExtractExifMetadata(t *testing.T) {
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

	if result == nil {
		fmt.Print("Nil received")
		t.FailNow()
	}

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if !reflect.DeepEqual(expectedMetadata, result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if !reflect.DeepEqual(expectedLocation, result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldExtractXmpMetadata(t *testing.T) {
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

	if result == nil {
		fmt.Print("Nil received")
		t.FailNow()
	}

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if !reflect.DeepEqual(expectedMetadata, result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if !reflect.DeepEqual(expectedLocation, result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestShouldExtractQuickTimeMetadata(t *testing.T) {
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

	if result == nil {
		fmt.Print("Nil received")
		t.FailNow()
	}

	expectedMetadata, _ := metadata.NewDate("2022.03.01 21:37:00")
	if !reflect.DeepEqual(expectedMetadata, result.CreationDate) {
		fmt.Print("Date mismatch")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if !reflect.DeepEqual(expectedLocation, result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

func TestDateShouldBeNilWhenThereAreNoCreationDateTags(t *testing.T) {
	meta := map[string]any{
		"Composite:GPSPosition": "51 6 32.29 N, 17 1 59.30 E",
		"File:MIMEType":         "image/jpeg",
	}
	rawMetadataExtractorMock := mock.NewRawMetadataExtractorMock(t)
	defer rawMetadataExtractorMock.AssertAllExpectionsSatisfied()
	rawMetadataExtractorMock.ExpectExtract(file, meta, nil)
	sut := metadata.NewMetadataExtractor(&rawMetadataExtractorMock)

	result := sut.Extract(file)

	if result == nil {
		fmt.Print("Nil received")
		t.FailNow()
	}

	if result.CreationDate != nil {
		fmt.Print("CreationDate expected to be nil")
		t.FailNow()
	}

	expectedLocation, _ := metadata.NewGPS("51 6 32.29 N, 17 1 59.30 E")
	if !reflect.DeepEqual(expectedLocation, result.Location) {
		fmt.Print("Location mismatch")
		t.FailNow()
	}

	if result.MIMEType != metadata.JPG {
		fmt.Print("MIMEType mismatch")
		t.FailNow()
	}
}

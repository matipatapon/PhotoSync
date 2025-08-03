package metadata_test

import (
	"fmt"
	"os"
	"photosync/src/metadata"
	"testing"
)

func expectThatContainsSpecificTag(m map[string]any, tag string, expectedValue string, t *testing.T) {
	value, ok := m[tag]
	if !ok {
		fmt.Printf("Missing '%s' tag", tag)
		t.FailNow()
	}

	valueStr, ok := value.(string)
	if !ok {
		fmt.Printf("'%s' tag has non-string type", tag)
		t.FailNow()
	}

	if valueStr != expectedValue {
		fmt.Printf("Expected'%s', Got '%s'", expectedValue, valueStr)
		t.FailNow()
	}
}

func runTest(filename string, t *testing.T) map[string]any {
	sut := metadata.NewRawMetadataExtractor("../../exiftool/exiftool")
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Print("Failed to access test file")
		t.FailNow()
	}
	result, err := sut.Extract(bytes)
	if err != nil {
		fmt.Printf("Something went wrong: '%s'", err.Error())
		t.FailNow()
	}
	return result
}

func TestShouldReturnExifMetadataFromJpgFile(t *testing.T) {
	result := runTest("../../test/images/exif.jpg", t)
	expectThatContainsSpecificTag(result, "File:MIMEType", "image/jpeg", t)
	expectThatContainsSpecificTag(result, "EXIF:DateTimeOriginal", "2023.06.07 12:30:45", t)
}

func TestShouldReturnIptcMetadataFromJpgFile(t *testing.T) {
	result := runTest("../../test/images/iptc.jpg", t)
	expectThatContainsSpecificTag(result, "File:MIMEType", "image/jpeg", t)
	expectThatContainsSpecificTag(result, "Composite:DateTimeOriginal", "2022.03.01 21:37:00", t)
}

func TestShouldReturnXmpMetadataFromJpgFile(t *testing.T) {
	result := runTest("../../test/images/xmp.jpg", t)
	expectThatContainsSpecificTag(result, "File:MIMEType", "image/jpeg", t)
	expectThatContainsSpecificTag(result, "XMP:CreateDate", "2007.12.20 00:24:13", t)
}

func TestShouldReturnQuickTimeMetadataFromMp4File(t *testing.T) {
	result := runTest("../../test/images/quick_time.mp4", t)
	expectThatContainsSpecificTag(result, "File:MIMEType", "video/mp4", t)
	expectThatContainsSpecificTag(result, "QuickTime:CreateDate", "2021.06.30 12:55:13", t)
}

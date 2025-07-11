package metadata_test

import (
	"fmt"
	"os"
	"photosync/src/metadata"
	"testing"
)

func TestShouldReturnMetadataFromJpgFile(t *testing.T) {
	sut := metadata.NewRawMetadataExtractor("../../exiftool/exiftool")
	bytes, err := os.ReadFile("../../test/images/1.jpg")
	if err != nil {
		fmt.Print("Failed to access test file")
		t.FailNow()
	}

	result, err := sut.Extract(bytes)
	if err != nil {
		fmt.Printf("Something went wrong: '%s'", err.Error())
		t.FailNow()
	}

	orginalDate, ok := result["Composite:SubSecDateTimeOriginal"]
	if !ok {
		fmt.Print("Missing Composite:SubSecDateTimeOriginal tag")
		t.FailNow()
	}

	orginalDateStr, ok := orginalDate.(string)
	if !ok {
		fmt.Print("Composite:SubSecDateTimeOriginal has invalid type")
		t.FailNow()
	}

	expectedDate := "2024.11.29 18:45:02"
	if orginalDateStr != expectedDate {
		fmt.Printf("'%s' != '%s'", orginalDateStr, expectedDate)
		t.FailNow()
	}

	gpsPosition, ok := result["Composite:GPSPosition"]
	if !ok {
		fmt.Print("Missing GPSPosition tag")
		t.FailNow()
	}

	gpsPositionStr, ok := gpsPosition.(string)
	if !ok {
		fmt.Print("Composite:GPSPosition has invalid type")
		t.FailNow()
	}

	expectedGps := "51 6 32.29 N, 17 1 59.30 E"
	if gpsPositionStr != expectedGps {
		fmt.Printf("'%s' != '%s'", gpsPositionStr, expectedGps)
		t.FailNow()
	}
}

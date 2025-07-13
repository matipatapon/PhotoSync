package metadata_test

import (
	"fmt"
	"photosync/src/metadata"
	"reflect"
	"testing"
)

func TestShouldConvertStringToGPS(t *testing.T) {
	stringToExpectedGPS := map[string]metadata.GPS{
		"50 30 40.32 N, 130 50 22.32 E": metadata.GPS{
			Latitude: metadata.Coordinate{
				Degree:     50,
				Minute:     30,
				Second:     40.32,
				Hemisphere: metadata.North},
			Longitude: metadata.Coordinate{
				Degree:     130,
				Minute:     50,
				Second:     22.32,
				Hemisphere: metadata.East}},
	}

	for str, expectedGPS := range stringToExpectedGPS {
		result, err := metadata.NewGPS(str)
		if err != nil || !reflect.DeepEqual(result, expectedGPS) {
			fmt.Println("Unexpected error or result")
			t.FailNow()
		}
	}
}

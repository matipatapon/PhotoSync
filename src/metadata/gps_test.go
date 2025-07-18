package metadata_test

import (
	"fmt"
	"photosync/src/metadata"
	"reflect"
	"testing"
)

func TestShouldConvertStringToGPS(t *testing.T) {
	stringToExpectedGPS := map[string]metadata.GPS{
		"50 30 40.32 N, 130 50 22.32 E": { // Return north and east
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
		"12 34 56.78 S, 9 11 12.13 E": { // Return south and east
			Latitude: metadata.Coordinate{
				Degree:     12,
				Minute:     34,
				Second:     56.78,
				Hemisphere: metadata.South},
			Longitude: metadata.Coordinate{
				Degree:     9,
				Minute:     11,
				Second:     12.13,
				Hemisphere: metadata.East}},
		"12 34 56.78 S, 9 11 12.13 W": { // Return south and west
			Latitude: metadata.Coordinate{
				Degree:     12,
				Minute:     34,
				Second:     56.78,
				Hemisphere: metadata.South},
			Longitude: metadata.Coordinate{
				Degree:     9,
				Minute:     11,
				Second:     12.13,
				Hemisphere: metadata.West}},
		"12 34 56 S, 9 11 12 W": { // Return rounded seconds
			Latitude: metadata.Coordinate{
				Degree:     12,
				Minute:     34,
				Second:     56,
				Hemisphere: metadata.South},
			Longitude: metadata.Coordinate{
				Degree:     9,
				Minute:     11,
				Second:     12,
				Hemisphere: metadata.West}},
		"90 0 0 S, 180 0 0 W": { // Return max degrees
			Latitude: metadata.Coordinate{
				Degree:     90,
				Minute:     0,
				Second:     0,
				Hemisphere: metadata.South},
			Longitude: metadata.Coordinate{
				Degree:     180,
				Minute:     0,
				Second:     0,
				Hemisphere: metadata.West}},
		"89 60 0 S, 179 60 0 W": { // Return max minutes
			Latitude: metadata.Coordinate{
				Degree:     89,
				Minute:     60,
				Second:     0,
				Hemisphere: metadata.South},
			Longitude: metadata.Coordinate{
				Degree:     179,
				Minute:     60,
				Second:     0,
				Hemisphere: metadata.West}},
		"89 59 60 S, 179 59 60 W": { // Return max seconds
			Latitude: metadata.Coordinate{
				Degree:     89,
				Minute:     59,
				Second:     60,
				Hemisphere: metadata.South},
			Longitude: metadata.Coordinate{
				Degree:     179,
				Minute:     59,
				Second:     60,
				Hemisphere: metadata.West}},
	}

	for str, expectedGPS := range stringToExpectedGPS {
		result, err := metadata.NewGPS(str)
		if err != nil || !reflect.DeepEqual(result, expectedGPS) {
			fmt.Println("Unexpected error or result")
			t.FailNow()
		}
	}
}

func TestShouldReturnErrorWhenStringIsInvalid(t *testing.T) {
	invalidStrings := []string{
		"50 30 40.32, 130 50 22.32 E",   // Missing first hemisphere
		"50 30 40.32 N, 130 50 22.32",   // Missing second hemisphere
		"50 30 40.32 N, 130 50 22.32 I", // Unknown hemisphere
		"U 30 40.32 N, 130 50 22.32 E",  // Degrees are not int
		"50 W 40.32 N, 130 50 22.32 E",  // Minutes are not int
		"50 30 40.32 N, 130 50 U E",     // Seconds are not float
		"50 30 40.32 N",                 // One coordinate is missing
		"12 34 56.78 S,",                // empty second coordinate part
		"12 34 56.78 S. 9 11 12.13 W",   // invalid coordinate separator
		"12-34-56.78-S, 9-11-12.13-W",   // invalid coordinate components separator
		"50 30 40.32 W, 130 50 22.32 E", // not north or south hemisphere in first part
		"50 30 40.32 E, 130 50 22.32 E", // not north or south hemisphere in first part
		"50 30 40.32 N, 130 50 22.32 S", // not east or west hemisphere in second part
		"50 30 40.32 N, 130 50 22.32 N", // not east or west hemisphere in second part
		"91 0 0 N, 130 50 22.32 E",      // latitude degree is above 90
		"90 1 0 N, 130 50 22.32 E",      // latitude degree is equal to 90 but minutes are not 0
		"90 0 1 N, 130 50 22.32 E",      // latitude degree is equal to 90 but seconds are not 0
		"55 0 0 N, 181 50 22.32 E",      // longitude degree is above 180
		"64 0 0 N, 180 1 0 E",           // longitude degree is equal to 180 but minutes are not 0
		"21 0 0 N, 180 0 1 E",           // longitude degree is equal to 180 but seconds are not 0
		"21 0 0 N, 179 61 0 E",          // minute is above 60 degrees
		"21 0 0 N, 179 60 1 E",          // minute is equal 60 degrees but seconds are non zero
		"21 0 0 N, 179 59 61 E",         // second is above 60 degrees
		"prefix 21 0 0 N, 179 59 0 E",   // has prefix
		"21 0 0 N, 179 59 0 E postfix",  // has postfix
	}

	for _, str := range invalidStrings {
		_, err := metadata.NewGPS(str)
		if err == nil {
			fmt.Printf("Expected error for '%s'", str)
			t.FailNow()
		}
	}
}

package metadata_test

import (
	"fmt"
	"photosync/src/metadata"
	"testing"
)

func TestShouldCreateDateFromString(t *testing.T) {
	date, err := metadata.NewDate("2025.07.06 17:39:55")
	if err != nil {
		fmt.Printf("Unexpected error: '%s'", err.Error())
		t.FailNow()
	}
	if date.Year != 2025 || date.Month != 7 || date.Day != 6 ||
		date.Hour != 17 || date.Minute != 39 || date.Second != 55 {
		fmt.Printf("Invalid date returned")
		t.FailNow()
	}
}

func TestShouldReturnErrorForInvalidString(t *testing.T) {
	invalidStrings := []string{
		"2025.07:06 17:39:55",          // Wrong separator between months and days
		"2025:07.06 17:39:55",          // Wrong separator between years and months
		"2025.07.06.05 17:39:55",       // Too many data in date part
		"A512.07.06 17:39:55",          // Year is not a int
		"2025.MT.06 17:39:55",          // Month is not a int
		"2025.07.DA 17:39:55",          // Day is not a int
		"17:39:55",                     // Missing date part
		"2025.07.06 17:39:55 17:39:55", // Too many parts
		"2025.07.06 17.39:55",          // Wrong separator between hour and minute
		"2025.07.06 17:39.55",          // Wrong separator between minute and second
		"2025.07.06 17:39:55:45",       // Too many data in time part
		"2025.07.06 B2:39:55",          // Hour is not a int
		"2025.07.06 17:9S:55",          // Minute is not a int
		"2025.07.06 17:39:A2",          // Second is not a int
	}

	for _, invalidString := range invalidStrings {
		_, err := metadata.NewDate(invalidString)
		if err == nil {
			fmt.Printf("Expected error for: '%s'", invalidString)
			t.FailNow()
		}
	}
}

func TestShouldCreateStringFromDate(t *testing.T) {
	dates := map[string]metadata.Date{
		"2014.10.11 12:30:44": {Year: 2014, Month: 10, Day: 11, Hour: 12, Minute: 30, Second: 44}, // Normal date
		"2014.01.02 03:04:05": {Year: 2014, Month: 1, Day: 2, Hour: 3, Minute: 4, Second: 5},      // Should add leading zeroes
	}

	for expectedString, date := range dates {
		result := date.ToString()
		if result != expectedString {
			fmt.Printf("Expected '%s', Got '%s'", expectedString, result)
			t.FailNow()
		}
	}
}

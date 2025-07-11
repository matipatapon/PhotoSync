package metadata

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Date struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

func toStringAndAddLeadingZero(number int) string {
	if number >= 10 {
		return fmt.Sprintf("%d", number)
	} else {
		return fmt.Sprintf("0%d", number)
	}
}

func NewDate(str string) (Date, error) {
	isFormatOk, _ := regexp.Match("^[0-9]{1,4}\\.[0-9]{2}\\.[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}$", []byte(str))
	if !isFormatOk {
		return Date{}, fmt.Errorf("'%s' has invalid format", str)
	}
	parts := strings.Split(str, " ")
	yearMonthDay := strings.Split(parts[0], ".")
	hourMinuteSecond := strings.Split(parts[1], ":")
	year, _ := strconv.Atoi(yearMonthDay[0])
	month, _ := strconv.Atoi(yearMonthDay[1])
	day, _ := strconv.Atoi(yearMonthDay[2])
	hour, _ := strconv.Atoi(hourMinuteSecond[0])
	minute, _ := strconv.Atoi(hourMinuteSecond[1])
	second, _ := strconv.Atoi(hourMinuteSecond[2])

	return Date{Year: year, Month: month, Day: day, Hour: hour, Minute: minute, Second: second}, nil
}

func (date *Date) ToString() string {
	return fmt.Sprintf("%d.%s.%s %s:%s:%s",
		date.Year,
		toStringAndAddLeadingZero(date.Month),
		toStringAndAddLeadingZero(date.Day),
		toStringAndAddLeadingZero(date.Hour),
		toStringAndAddLeadingZero(date.Minute),
		toStringAndAddLeadingZero(date.Second),
	)
}

package metadata

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Hemisphere int

const (
	East  Hemisphere = iota
	West  Hemisphere = iota
	North Hemisphere = iota
	South Hemisphere = iota
)

type Coordinate struct {
	Degree     int
	Minute     int
	Second     float32
	Hemisphere Hemisphere
}

type GPS struct {
	Latitude  Coordinate
	Longitude Coordinate
}

func NewGPS(str string) (GPS, error) {
	isFormatOk, _ := regexp.Match(gpsRegex, []byte(str))
	if !isFormatOk {
		return GPS{}, fmt.Errorf("'%s' has invalid format", str)
	}
	coordinates := strings.Split(str, ",")
	latitude, err := stringToCoordinate(coordinates[0])
	if err != nil {
		return GPS{}, err
	}
	longitude, err := stringToCoordinate(coordinates[1][1:])
	if err != nil {
		return GPS{}, err
	}
	return GPS{Latitude: latitude, Longitude: longitude}, nil
}

var gpsRegex string = "^[0-9]{1,3} [0-9]{1,2} [0-9]{1,2}(\\.[0-9]{1,2})? [N,S], [0-9]{1,3} [0-9]{1,2} [0-9]{1,2}(\\.[0-9]{1,2})? [E,W]$"

func stringToCoordinate(str string) (Coordinate, error) {
	parts := strings.Split(str, " ")
	degree, _ := strconv.ParseInt(parts[0], 10, 32)
	minute, _ := strconv.ParseInt(parts[1], 10, 32)
	second, _ := strconv.ParseFloat(parts[2], 32)
	hemisphere := stringToHemishpere(parts[3])

	if !isCoordinateValid(degree, minute, second, hemisphere) {
		return Coordinate{}, fmt.Errorf("'%s' coordniate is invalid", str)
	}

	return Coordinate{
		Degree:     int(degree),
		Minute:     int(minute),
		Second:     float32(second),
		Hemisphere: hemisphere,
	}, nil
}

func isCoordinateValid(degree int64, minute int64, second float64, hemisphere Hemisphere) bool {
	maxDegree := 180
	if hemisphere == South || hemisphere == North {
		maxDegree = 90
	}
	if int(degree) > maxDegree || (int(degree) == maxDegree && (minute != 0 || second != 0)) {
		return false
	}
	if minute > 60 || (minute == 60 && second != 0) {
		return false
	}
	if second > 60 {
		return false
	}
	return true
}

func stringToHemishpere(str string) Hemisphere {
	switch str {
	case "N":
		return North
	case "S":
		return South
	case "W":
		return West
	default:
		return East
	}
}

package metadata

import (
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
	coordinates := strings.Split(str, ",")
	latitude, _ := stringToCoordinate(coordinates[0])
	longitude, _ := stringToCoordinate(coordinates[1][1:])
	return GPS{Latitude: latitude, Longitude: longitude}, nil
}

func stringToCoordinate(str string) (Coordinate, error) {
	parts := strings.Split(str, " ")
	degree, _ := strconv.ParseInt(parts[0], 10, 32)
	minute, _ := strconv.ParseInt(parts[1], 10, 32)
	second, _ := strconv.ParseFloat(parts[2], 32)
	hemisphere := East
	if parts[3] == "N" {
		hemisphere = North
	}
	return Coordinate{
		Degree:     int(degree),
		Minute:     int(minute),
		Second:     float32(second),
		Hemisphere: hemisphere,
	}, nil
}

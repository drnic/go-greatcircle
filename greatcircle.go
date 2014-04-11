package gogreatcircle

import (
	"math"
)

type vector [3]float64

func cross(v1, v2 vector) vector {
	return vector{
		v1[1]*v2[2] - v1[2]*v2[1],
		v1[2]*v2[0] - v1[0]*v2[2],
		v1[0]*v2[1] - v1[1]*v2[0],
	}
}

func CoordinateToDecimalDegree(degrees, minutes, seconds float64) float64 {
	return degrees + (minutes / 60) + (seconds / 3600)
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func RadiansToDegrees(distanceRadians float64) float64 {
	return distanceRadians * (180 / math.Pi)
}

func NMToRadians(nauticalMiles float64) float64 {
	return (math.Pi / (180 * 60)) * nauticalMiles
}

func RadiansToNM(radians float64) float64 {
	return ((180 * 60) / math.Pi) * radians
}

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	return ((math.Acos(math.Sin(lat1)*math.Sin(lat2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon1-lon2))) * 180 * 60) / math.Pi
}

func InitialBearing(lat1, lon1, lat2, lon2 float64) float64 {
	dLon := (lon2 - lon1)
	y := math.Sin(dLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLon)
	// bearing in radians
	bearing := math.Atan2(y, x)
	// convert the bearing back to degrees to get the compas bearing
	return math.Mod(RadiansToDegrees(bearing)+360, 360)
}

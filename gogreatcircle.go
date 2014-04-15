package gogreatcircle

import (
	"errors"
	"math"
)

type vector [3]float64

type Coordinate struct {
	latitude  float64
	longitude float64
}

type Radial struct {
	Coordinate
	bearing float64
}

func CoordinateToDecimalDegree(degrees, minutes, seconds float64) float64 {
	// converting coordiantes from degrees, minutes, seconds into decimal degrees
	return degrees + (minutes / 60) + (seconds / 3600)
}

func DegreesToRadians(degrees float64) float64 {
	// converting decimal degrees to radians
	return degrees * (math.Pi / 180)
}

func RadiansToDegrees(radians float64) float64 {
	// converting radians to decimal degrees
	return radians * (180 / math.Pi)
}

func NMToRadians(nauticalMiles float64) float64 {
	// converting nautical miles to radians
	return (math.Pi / (180 * 60)) * nauticalMiles
}

func RadiansToNM(radians float64) float64 {
	// converting radians to nautical miles
	return ((180 * 60) / math.Pi) * radians
}

func Distance(point1, point2 *Coordinate) float64 {
	// distance between 2 coordiantes, returned in nautical miles
	return (math.Acos(math.Sin(point1.latitude)*math.Sin(point2.latitude)+math.Cos(point1.latitude)*math.Cos(point2.latitude)*math.Cos(point1.longitude-point2.longitude)) * 180 * 60) / math.Pi

}

func InitialBearing(point1, point2 *Coordinate) float64 {
	// calculate the initial true course from point1 to point2
	dLon := (point2.longitude - point1.longitude)
	y := math.Sin(dLon) * math.Cos(point2.latitude)
	x := math.Cos(point1.latitude)*math.Sin(point2.latitude) - math.Sin(point1.latitude)*math.Cos(point2.latitude)*math.Cos(dLon)
	// bearing calculated in radians
	bearing := math.Atan2(y, x)
	return bearing
}

func IntersectionRadials(radial1, radial2 *Radial) (coordinate Coordinate, err error) {
	// adapted from http://williams.best.vwh.net/avform.htm#Intersection
	dLat := radial2.Coordinate.latitude - radial1.Coordinate.latitude
	dLon := radial2.Coordinate.longitude - radial1.Coordinate.longitude

	dist12 := 2 * math.Asin(math.Sqrt(math.Sin(dLat/2)*math.Sin(dLat/2)+math.Cos(radial1.Coordinate.latitude)*math.Cos(radial2.Coordinate.latitude)*math.Sin(dLon/2)*math.Sin(dLon/2)))
	if dist12 == 0 {
		return Coordinate{0, 0}, errors.New("dist 0")
	}

	// initial/final bearings between points
	brngA := math.Acos((math.Sin(radial2.Coordinate.latitude) - math.Sin(radial1.Coordinate.latitude)*math.Cos(dist12)) / (math.Sin(dist12) * math.Cos(radial1.Coordinate.latitude)))
	brngB := math.Acos((math.Sin(radial1.Coordinate.latitude) - math.Sin(radial2.Coordinate.latitude)*math.Cos(dist12)) / (math.Sin(dist12) * math.Cos(radial2.Coordinate.latitude)))

	var brng12 float64
	var brng21 float64
	if math.Sin(radial2.Coordinate.longitude-radial1.Coordinate.longitude) > 0 {
		brng12 = brngA
		brng21 = 2*math.Pi - brngB
	} else {
		brng12 = 2*math.Pi - brngA
		brng21 = brngB
	}

	alpha1 := math.Mod((radial1.bearing-brng12+math.Pi), (2*math.Pi)) - math.Pi // angle 2-1-3
	alpha2 := math.Mod((brng21-radial2.bearing+math.Pi), (2*math.Pi)) - math.Pi // angle 1-2-3

	if math.Sin(alpha1) == 0 && math.Sin(alpha2) == 0 {
		return Coordinate{0, 0}, errors.New("infinite intersections")
	}
	if math.Sin(alpha1)*math.Sin(alpha2) < 0 {
		return Coordinate{0, 0}, errors.New("ambiguous intersection")
	}

	alpha3 := math.Acos(-math.Cos(alpha1)*math.Cos(alpha2) + math.Sin(alpha1)*math.Sin(alpha2)*math.Cos(dist12))
	dist13 := math.Atan2(math.Sin(dist12)*math.Sin(alpha1)*math.Sin(alpha2), math.Cos(alpha2)+math.Cos(alpha1)*math.Cos(alpha3))
	// latitude of the intersection point
	coordinate.latitude = math.Asin(math.Sin(radial1.Coordinate.latitude)*math.Cos(dist13) + math.Cos(radial1.Coordinate.latitude)*math.Sin(dist13)*math.Cos(radial1.bearing))
	dLon13 := math.Atan2(math.Sin(radial1.bearing)*math.Sin(dist13)*math.Cos(radial1.Coordinate.latitude), math.Cos(dist13)-math.Sin(radial1.Coordinate.latitude)*math.Sin(coordinate.latitude))
	// longitude of intersection point
	coordinate.longitude = radial1.Coordinate.longitude + dLon13
	coordinate.longitude = math.Mod((coordinate.longitude+3*math.Pi), (2*math.Pi)) - math.Pi // normalise to -180..+180º

	return coordinate, nil

}

func CrossTrackError(point1, point2, point3 *Coordinate) float64 {
	// distance between point A and point D
	dist_AD := NMToRadians(Distance(point1, point3))
	// course of point A to point D
	crs_AD := math.Acos((math.Sin(point3.latitude) - math.Sin(point1.latitude)*math.Cos(dist_AD)) / (math.Sin(dist_AD) * math.Cos(point1.latitude)))
	initialBearing := InitialBearing(point1, point2)
	// crosstrack error
	xtd := math.Asin(math.Sin(dist_AD) * math.Sin(crs_AD-initialBearing))
	return xtd
}

func AlongTrackDistance(point1, point2, point3 *Coordinate) float64 {
	// distance between point A and point D
	dist_AD := NMToRadians(Distance(point1, point3))
	// along track distance
	xtd := CrossTrackError(point1, point2, point3)
	atd := math.Asin(math.Sqrt(math.Pow((math.Sin(dist_AD)), 2)-math.Pow((math.Sin(xtd)), 2)) / math.Cos(xtd))
	// http://williams.best.vwh.net/avform.htm#XTE - "Note that we can also use the above formulae to find the point of closest approach to the point D on the great circle through A and B"
	return atd
}

func ClosestPoint(point1, point2, point3 *Coordinate) (coordinate Coordinate) {
	// coordinates on the route from point1 to point2 of a given point3
	// calculated using the formula from http://williams.best.vwh.net/avform.htm#Example - enroute waypoint
	bearing := InitialBearing(point1, point2)
	distance := AlongTrackDistance(point1, point2, point3)
	coordinate.latitude = math.Asin(math.Sin(point1.latitude)*math.Cos(distance) + math.Cos(point1.latitude)*math.Sin(distance)*math.Cos(bearing))
	coordinate.longitude = -(math.Mod(math.Abs(point1.longitude)-math.Asin(math.Sin(bearing)*math.Sin(distance)/math.Cos(coordinate.latitude))+math.Pi, 2*math.Pi) - math.Pi)
	return coordinate
}

func PointInReach(point1, point2, point3 *Coordinate, distance float64) (response bool) {
	// first we use the ClosestPoint function to get the first point and then compute
	//the distance to the given point3 and compare it against the given distance
	closestpoint := ClosestPoint(point1, point2, point3)
	distanceBetweenPoints := Distance(&closestpoint, point3)
	if distanceBetweenPoints <= distance {
		return true
	} else {
		return false
	}
}

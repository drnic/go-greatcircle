package gogreatcircle

import (
	"errors"
	"math"
	"sort"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type Radial struct {
	Coordinate
	Bearing float64
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

func Distance(point1, point2 Coordinate) float64 {
	// distance between 2 coordiantes, returned in nautical miles
	return (math.Acos(math.Sin(point1.Latitude)*math.Sin(point2.Latitude)+math.Cos(point1.Latitude)*math.Cos(point2.Latitude)*math.Cos(point1.Longitude-point2.Longitude)) * 180 * 60) / math.Pi

}

func InitialBearing(point1, point2 Coordinate) float64 {
	// calculate the initial true course from point1 to point2
	dLon := (point2.Longitude - point1.Longitude)
	y := math.Sin(dLon) * math.Cos(point2.Latitude)
	x := math.Cos(point1.Latitude)*math.Sin(point2.Latitude) - math.Sin(point1.Latitude)*math.Cos(point2.Latitude)*math.Cos(dLon)
	// bearing calculated in radians
	Bearing := math.Atan2(y, x)
	return Bearing
}

func IntersectionRadials(radial1, radial2 Radial) (coordinate Coordinate, err error) {
	// adapted from http://williams.best.vwh.net/avform.htm#Intersection
	dLat := radial2.Coordinate.Latitude - radial1.Coordinate.Latitude
	dLon := radial2.Coordinate.Longitude - radial1.Coordinate.Longitude

	dist12 := 2 * math.Asin(math.Sqrt(math.Sin(dLat/2)*math.Sin(dLat/2)+math.Cos(radial1.Coordinate.Latitude)*math.Cos(radial2.Coordinate.Latitude)*math.Sin(dLon/2)*math.Sin(dLon/2)))
	if dist12 == 0 {
		return Coordinate{0, 0}, errors.New("dist 0")
	}

	// initial/final bearings between points
	brngA := math.Acos((math.Sin(radial2.Coordinate.Latitude) - math.Sin(radial1.Coordinate.Latitude)*math.Cos(dist12)) / (math.Sin(dist12) * math.Cos(radial1.Coordinate.Latitude)))
	brngB := math.Acos((math.Sin(radial1.Coordinate.Latitude) - math.Sin(radial2.Coordinate.Latitude)*math.Cos(dist12)) / (math.Sin(dist12) * math.Cos(radial2.Coordinate.Latitude)))

	var brng12 float64
	var brng21 float64
	if math.Sin(radial2.Coordinate.Longitude-radial1.Coordinate.Longitude) > 0 {
		brng12 = brngA
		brng21 = 2*math.Pi - brngB
	} else {
		brng12 = 2*math.Pi - brngA
		brng21 = brngB
	}

	alpha1 := math.Mod((radial1.Bearing-brng12+math.Pi), (2*math.Pi)) - math.Pi // angle 2-1-3
	alpha2 := math.Mod((brng21-radial2.Bearing+math.Pi), (2*math.Pi)) - math.Pi // angle 1-2-3

	if math.Sin(alpha1) == 0 && math.Sin(alpha2) == 0 {
		return Coordinate{0, 0}, errors.New("infinite intersections")
	}
	if math.Sin(alpha1)*math.Sin(alpha2) < 0 {
		return Coordinate{0, 0}, errors.New("ambiguous intersection")
	}

	alpha3 := math.Acos(-math.Cos(alpha1)*math.Cos(alpha2) + math.Sin(alpha1)*math.Sin(alpha2)*math.Cos(dist12))
	dist13 := math.Atan2(math.Sin(dist12)*math.Sin(alpha1)*math.Sin(alpha2), math.Cos(alpha2)+math.Cos(alpha1)*math.Cos(alpha3))
	// latitude of the intersection point
	coordinate.Latitude = math.Asin(math.Sin(radial1.Coordinate.Latitude)*math.Cos(dist13) + math.Cos(radial1.Coordinate.Latitude)*math.Sin(dist13)*math.Cos(radial1.Bearing))
	dLon13 := math.Atan2(math.Sin(radial1.Bearing)*math.Sin(dist13)*math.Cos(radial1.Coordinate.Latitude), math.Cos(dist13)-math.Sin(radial1.Coordinate.Latitude)*math.Sin(coordinate.Latitude))
	// longitude of intersection point
	coordinate.Longitude = radial1.Coordinate.Longitude + dLon13
	coordinate.Longitude = math.Mod((coordinate.Longitude+3*math.Pi), (2*math.Pi)) - math.Pi // normalise to -180..+180º

	return coordinate, nil

}

func CrossTrackError(point1, point2, point3 Coordinate) float64 {
	// distance between point A and point D
	dist_AD := NMToRadians(Distance(point1, point3))
	// course of point A to point D
	crs_AD := math.Acos((math.Sin(point3.Latitude) - math.Sin(point1.Latitude)*math.Cos(dist_AD)) / (math.Sin(dist_AD) * math.Cos(point1.Latitude)))
	initialBearing := InitialBearing(point1, point2)
	// crosstrack error
	xtd := math.Asin(math.Sin(dist_AD) * math.Sin(crs_AD-initialBearing))
	return xtd
}

func AlongTrackDistance(point1, point2, point3 Coordinate) float64 {
	// distance between point A and point D
	dist_AD := NMToRadians(Distance(point1, point3))
	// along track distance
	xtd := CrossTrackError(point1, point2, point3)
	atd := math.Asin(math.Sqrt(math.Pow((math.Sin(dist_AD)), 2)-math.Pow((math.Sin(xtd)), 2)) / math.Cos(xtd))
	// http://williams.best.vwh.net/avform.htm#XTE - "Note that we can also use the above formulae to find the point of closest approach to the point D on the great circle through A and B"
	return atd
}

func ClosestPoint(point1, point2, point3 Coordinate) (coordinate Coordinate) {
	// coordinates on the route from point1 to point2 of a given point3
	// calculated using the formula from http://williams.best.vwh.net/avform.htm#Example - enroute waypoint
	Bearing := InitialBearing(point1, point2)
	distance := AlongTrackDistance(point1, point2, point3)
	coordinate.Latitude = math.Asin(math.Sin(point1.Latitude)*math.Cos(distance) + math.Cos(point1.Latitude)*math.Sin(distance)*math.Cos(Bearing))
	earthRadius := NMToRadians(3440.07)
	coordinate.Longitude = point1.Longitude + math.Atan2(math.Sin(Bearing)*math.Sin(distance/earthRadius)*math.Cos(point1.Latitude), math.Cos(distance/earthRadius)-math.Sin(point1.Latitude)*math.Sin(point2.Latitude))
	return coordinate
}

// helper function for Pointinreach and Pointsinreach

func pointOfReachDistance(point1, point2, point3 Coordinate) float64 {
	// first we use the ClosestPoint function to get the first point (the closest to the provided point3)
	// and then compute the distance to the given point3 and compare it against the given distance
	closestpoint := ClosestPoint(point1, point2, point3)
	distanceBetweenPoints := Distance(closestpoint, point3)

	return distanceBetweenPoints
}

func PointInReach(point1, point2, point3 Coordinate, distance float64) (response bool) {
	// using the helper function above we find the point3 in reach and get the distance
	// of to point3. Comparing the expected distance with the provided distance
	//the function returns true if point3 is in range, else false
	distanceBetweenPoints := pointOfReachDistance(point1, point2, point3)
	if distanceBetweenPoints <= distance {
		return true
	} else {
		return false
	}
}

func PointsInReach(point1, point2 Coordinate, distance float64, points []Coordinate) []Coordinate {
	// providing an array of points, the helper function is used to get the distance
	// to those points and then compared with the distance provided.
	// If the points are within distance, they are returned sorted
	pointsInReach := make(map[float64]Coordinate)
	var sortedPointsInReach []Coordinate

	for _, point := range points {
		distanceBetweenPoints := pointOfReachDistance(point1, point2, point)
		if distanceBetweenPoints <= distance {
			pointsInReach[distanceBetweenPoints] = point
		}

	}
	keys := make([]float64, 0, len(pointsInReach))
	for k := range pointsInReach {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	for _, k := range keys {
		sortedPointsInReach = append(sortedPointsInReach, pointsInReach[k])
	}

	return sortedPointsInReach

}

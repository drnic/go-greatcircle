package gogreatcircle

/*
Library for various Earth coordinate & Great Circle calcuations.

North latitudes and West longitudes are treated as positive, and South latitudes and East longitudes negative
*/

import (
	"errors"
	"math"
	"sort"
)

/*
Coordinate is the position on earth in latitude/longitude.

Negative longitude represents a longitude line in the western hemisphere.
*/
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

/*
Radial describes a great circle from a specific starting coordinate
upon an initial bearing.

Across the entire great circle, the bearing will be different at different
coordinates. */
type Radial struct {
	Coordinate
	Bearing float64
}

/*
DegreeUnitsToDecimalDegree converts from (degrees, minutes, seconds) into decimal degrees
*/
func DegreeUnitsToDecimalDegree(degrees, minutes, seconds float64) float64 {
	return degrees + (minutes / 60) + (seconds / 3600)
}

/*
DegreesToRadians converts a decimal degree into radians.
*/
func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

/*
RadiansToDegrees converts a radian into decimal degrees.
*/
func RadiansToDegrees(radians float64) float64 {
	return radians * (180 / math.Pi)
}

/*
NMToRadians converts nautical miles into radians.
*/
func NMToRadians(nauticalMiles float64) float64 {
	return (math.Pi / (180 * 60)) * nauticalMiles
}

/*
RadiansToNM converts radians into nautical miles.
*/
func RadiansToNM(radians float64) float64 {
	return ((180 * 60) / math.Pi) * radians
}

/*
Distance calculates the shortest distance between two Coordinates.

Result is in nautical miles.

The shortest distance between two coordinates is the arc across the
great circle that includes the two points.
*/
func Distance(point1, point2 Coordinate) float64 {
	return (math.Acos(math.Sin(point1.Latitude)*math.Sin(point2.Latitude)+math.Cos(point1.Latitude)*math.Cos(point2.Latitude)*math.Cos(point1.Longitude-point2.Longitude)) * 180 * 60) / math.Pi

}

/*
InitialBearing provides the initial true course from a point to commence
a journey along a great circle to another point on that great
circle.

Result is in radians.

The bearing being used whilst travelling along a great circle
will change. This function returns the bearing at point1.
*/
func InitialBearing(point1, point2 Coordinate) float64 {
	dLon := (point2.Longitude - point1.Longitude)
	y := math.Sin(dLon) * math.Cos(point2.Latitude)
	x := math.Cos(point1.Latitude)*math.Sin(point2.Latitude) - math.Sin(point1.Latitude)*math.Cos(point2.Latitude)*math.Cos(dLon)
	// bearing calculated in radians
	Bearing := math.Atan2(y, x)
	return Bearing
}

/*
IntersectionRadials determines the Coordinate that two Radials
would interset.

Adapted from http://williams.best.vwh.net/avform.htm#Intersection
*/
func IntersectionRadials(radial1, radial2 Radial) (coordinate Coordinate, err error) {
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

/*
CrossTrackError (XTD) determines the distance off course.

Positive XTD means right of course, negative means left.

See http://williams.best.vwh.net/avform.htm#XTE for more information.
*/
func CrossTrackError(routeStartCoord, routeEndCoord, actualCoord Coordinate) float64 {
	// distance between point A and point D
	distAD := NMToRadians(Distance(routeStartCoord, actualCoord))
	// course of point A to point D
	crsAD := math.Acos((math.Sin(actualCoord.Latitude) -
		math.Sin(routeStartCoord.Latitude)*math.Cos(distAD)) / (math.Sin(distAD) * math.Cos(routeStartCoord.Latitude)))
	initialBearing := InitialBearing(routeStartCoord, routeEndCoord)
	// crosstrack error
	xtd := math.Asin(math.Sin(distAD) * math.Sin(crsAD-initialBearing))
	return xtd
}

/*
AlongTrackDistance is the distance from routeStartCoord along the course towards routeEndCoord to the point abeam actualCoord.

See http://williams.best.vwh.net/avform.htm#XTE for more information.
"Note that we can also use the above formulae to find the point of closest approach to the point D on the great circle through A and B"
*/
func AlongTrackDistance(routeStartCoord, routeEndCoord, actualCoord Coordinate) float64 {
	// distance between point A and point D
	distAD := NMToRadians(Distance(routeStartCoord, actualCoord))
	// along track distance
	xtd := CrossTrackError(routeStartCoord, routeEndCoord, actualCoord)
	atd := math.Asin(math.Sqrt(math.Pow((math.Sin(distAD)), 2)-math.Pow((math.Sin(xtd)), 2)) / math.Cos(xtd))
	return atd
}

/*
ClosestPoint determines the coordinate for the closest point along a course/radial from the actualCoord.

Calculated using the formula from http://williams.best.vwh.net/avform.htm#Example - enroute waypoint.
*/
func ClosestPoint(routeStartCoord, routeEndCoord, actualCoord Coordinate) (coordinate Coordinate) {
	Bearing := InitialBearing(routeStartCoord, routeEndCoord)
	distance := AlongTrackDistance(routeStartCoord, routeEndCoord, actualCoord)
	coordinate.Latitude = math.Asin(math.Sin(routeStartCoord.Latitude)*math.Cos(distance) +
		math.Cos(routeStartCoord.Latitude)*math.Sin(distance)*math.Cos(Bearing))
	earthRadius := NMToRadians(3440.07)
	coordinate.Longitude = routeStartCoord.Longitude +
		math.Atan2(math.Sin(Bearing)*math.Sin(distance/earthRadius)*math.Cos(routeStartCoord.Latitude),
			math.Cos(distance/earthRadius)-math.Sin(routeStartCoord.Latitude)*math.Sin(routeEndCoord.Latitude))
	return coordinate
}

/*
pointOfReachDistance is a helper function for PointInReach and PointsInReach
first we use the ClosestPoint function to get the first point (the closest to the provided point3)
and then compute the distance to the given point3 and compare it against the given distance
*/
func pointOfReachDistance(point1, point2, point3 Coordinate) float64 {
	closestpoint := ClosestPoint(point1, point2, point3)
	distanceBetweenPoints := Distance(closestpoint, point3)

	return distanceBetweenPoints
}

/*
PointInReach determines if actualCoord is within testDistance of the route from

*/
func PointInReach(point1, point2, point3 Coordinate, distance float64) (response bool) {
	// using the helper function above we find the point3 in reach and get the distance
	// of to point3. Comparing the expected distance with the provided distance
	//the function returns true if point3 is in range, else false
	distanceBetweenPoints := pointOfReachDistance(point1, point2, point3)
	return distanceBetweenPoints <= distance
}

/*
PointsInReach filters a list of Coordinates to return only those Coordinates that are within testDistance
of the (routeStartCoord, routeEndCoord) route.
*/
func PointsInReach(routeStartCoord, routeEndCoord Coordinate, distance float64, coords []Coordinate) []Coordinate {
	// providing an array of points, the helper function is used to get the distance
	// to those points and then compared with the distance provided.
	// If the points are within distance, they are returned sorted
	pointsInReach := make(map[float64]Coordinate)
	var sortedPointsInReach []Coordinate

	for _, coord := range coords {
		distanceBetweenPoints := pointOfReachDistance(routeStartCoord, routeEndCoord, coord)
		if distanceBetweenPoints <= distance {
			pointsInReach[distanceBetweenPoints] = coord
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

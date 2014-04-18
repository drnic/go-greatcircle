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
NamedCoordinate includes a name or label describing the Coordinate
*/
type NamedCoordinate struct {
	Coord Coordinate
	Name  string
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
	return (math.Acos(math.Sin(point1.Latitude)*math.Sin(point2.Latitude)+
		math.Cos(point1.Latitude)*math.Cos(point2.Latitude)*math.Cos(point1.Longitude-point2.Longitude)) * 180 * 60) / math.Pi
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
	distance := NMToRadians(Distance(point1, point2))
	tcl := math.Acos((math.Sin(point2.Latitude) - math.Sin(point1.Latitude)*math.Cos(distance)) / (math.Sin(distance) * math.Cos(point1.Latitude)))
	return tcl
}

/*
IntersectionRadials determines the Coordinate that two Radials
would interset.

Adapted from http://williams.best.vwh.net/avform.htm#Intersection
*/
func IntersectionRadials(radial1, radial2 Radial) (coordinate Coordinate, err error) {
	dst12 := 2 * math.Asin(math.Sqrt(math.Pow((math.Sin((radial1.Latitude-radial2.Latitude)/2)), 2)+
		math.Cos(radial1.Latitude)*math.Cos(radial2.Latitude)*math.Pow(math.Sin((radial1.Longitude-radial2.Longitude)/2), 2)))
	// // initial/final bearings between points
	crs13 := InitialBearing(radial1.Coordinate, radial2.Coordinate)
	crs23 := InitialBearing(radial2.Coordinate, radial1.Coordinate)

	var crs12 float64
	var crs21 float64
	if math.Sin(radial2.Longitude-radial1.Longitude) < 0 {
		crs12 = math.Acos((math.Sin(radial2.Latitude) - math.Sin(radial1.Latitude)*math.Cos(dst12)) / (math.Sin(dst12) * math.Cos(radial1.Latitude)))
		crs21 = 2.*math.Pi - math.Acos((math.Sin(radial1.Latitude)-math.Sin(radial2.Latitude)*math.Cos(dst12))/(math.Sin(dst12)*math.Cos(radial2.Latitude)))
	} else {
		crs12 = 2.*math.Pi - math.Acos((math.Sin(radial2.Latitude)-math.Sin(radial1.Latitude)*math.Cos(dst12))/(math.Sin(dst12)*math.Cos(radial1.Latitude)))
		crs21 = math.Acos((math.Sin(radial1.Latitude) - math.Sin(radial2.Latitude)*math.Cos(dst12)) / (math.Sin(dst12) * math.Cos(radial2.Latitude)))
	}

	ang1 := math.Mod(crs13-crs12+math.Pi, 2.*math.Pi) - math.Pi
	ang2 := math.Mod(crs21-crs23+math.Pi, 2.*math.Pi) - math.Pi

	if math.Sin(ang1) == 0 && math.Sin(ang2) == 0 {
		return Coordinate{0, 0}, errors.New("infinity of intersections")
	} else if math.Sin(ang1)*math.Sin(ang2) < 0 {
		return Coordinate{0, 0}, errors.New("intersection ambiguous")
	} else {
		ang1 := math.Abs(ang1)
		ang2 := math.Abs(ang2)
		ang3 := math.Acos(-math.Cos(ang1)*math.Cos(ang2) + math.Sin(ang1)*math.Sin(ang2)*math.Cos(dst12))
		dst13 := math.Atan2(math.Sin(dst12)*math.Sin(ang1)*math.Sin(ang2), math.Cos(ang2)+math.Cos(ang1)*math.Cos(ang3))
		lat3 := math.Asin(math.Sin(radial1.Latitude)*math.Cos(dst13) + math.Cos(radial1.Latitude)*math.Sin(dst13)*math.Cos(crs13))
		dlon := math.Atan2(math.Sin(crs13)*math.Sin(dst13)*math.Cos(radial1.Latitude), math.Cos(dst13)-math.Sin(radial1.Latitude)*math.Sin(lat3))
		lon3 := math.Mod(radial1.Longitude-dlon+math.Pi, 2*math.Pi) - math.Pi
		return Coordinate{lat3, lon3}, nil
	}
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
	crsAD := InitialBearing(routeStartCoord, actualCoord)
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
	// earthRadius := NMToRadians(3440.07)
	// coordinate.Longitude = routeStartCoord.Longitude +
	// math.Atan2(math.Sin(Bearing)*math.Sin(distance/earthRadius)*math.Cos(routeStartCoord.Latitude),
	// math.Cos(distance/earthRadius)-math.Sin(routeStartCoord.Latitude)*math.Sin(routeEndCoord.Latitude))
	coordinate.Longitude = math.Mod(routeStartCoord.Longitude-math.Asin(math.Sin(Bearing)*math.Sin(distance)/math.Cos(coordinate.Latitude))+math.Pi, 2*math.Pi) - math.Pi
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

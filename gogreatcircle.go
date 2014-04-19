package greatcircle

/*
Library for various Earth coordinate & Great Circle calcuations.

North latitudes and West longitudes are treated as positive, and South latitudes and East longitudes negative
*/

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

/*
DegreeUnits is a value of latitude or longitude in degrees.
*/
type DegreeUnits struct {
	degree, minute, second float64
}

/*
Coordinate is the position on earth in latitude/longitude.

Negative longitude represents a longitude line in the western hemisphere.
*/
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

func NewCoordinate(latitude string, longitude string) (Coordinate, error) {
	latitudeDegrees, err := DegreeStrToDecimalDegree(latitude)
	if err != nil {
		return Coordinate{}, err
	}
	longitudeDegrees, err := DegreeStrToDecimalDegree(longitude)
	if err != nil {
		return Coordinate{}, err
	}
	return Coordinate{DegreesToRadians(latitudeDegrees), DegreesToRadians(longitudeDegrees)}, nil
}

/*
Equal compares this Coordinate to another and determines if their
Latitude & Longitudes are equivalent (to 3 decimal places)
*/
func (coord Coordinate) Equal(another Coordinate) bool {
	return math.Abs(coord.Latitude-another.Latitude) <= 0.001 &&
		math.Abs(coord.Longitude-another.Longitude) <= 0.001
}

/*
ToSkyVector returns a string that is a valid coordinate for
skyvector.com.

Response format: latitude:longitude

SkyVector uses -ve for west & +ve for east.
*/
func (coord Coordinate) ToSkyVector() (out string) {
	out = strconv.FormatFloat(RadiansToDegrees(coord.Latitude), 'f', 2, 64)
	out = out + ":"
	out = out + strconv.FormatFloat(-1*RadiansToDegrees(coord.Longitude), 'f', 2, 64)
	return
}

/*
NamedCoordinate includes a name or label describing the Coordinate
*/
type NamedCoordinate struct {
	Coord Coordinate
	Name  string
}

func NewNamedCoordinate(name string, latitude string, longitude string) (NamedCoordinate, error) {
	coord, err := NewCoordinate(latitude, longitude)
	if err != nil {
		return NamedCoordinate{}, err
	}
	return NamedCoordinate{coord, name}, nil
}

/*
Equal compares this NamedCoordinate to another and determines if their
Latitude & Longitudes are equivalent (to 3 decimal places); and
their Name is equivalent
*/
func (thisCoord NamedCoordinate) Equal(another NamedCoordinate) bool {
	return thisCoord.Equal(another) && thisCoord.Name == another.Name
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

func DegreeStrToDecimalDegree(degrees string) (float64, error) {
	units := strings.Split(degrees, ":")
	if len(units) > 3 {
		return 0, errors.New(degrees + " should have 3 or fewer portions")
	}
	unitMultiplier := 1.0
	decimalDegree := 0.0
	for _, unit := range units {
		unitValue, err := strconv.ParseFloat(unit, 64)
		if err != nil {
			return 0, err
		}
		decimalDegree = decimalDegree + unitValue/unitMultiplier
		unitMultiplier = unitMultiplier * 60
	}
	return decimalDegree, nil
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
	var tc float64
	var argacos float64
	d := math.Acos(math.Sin(point1.Latitude)*math.Sin(point2.Latitude) + math.Cos(point1.Latitude)*math.Cos(point2.Latitude)*math.Cos(point1.Longitude-point2.Longitude))
	if (d == 0.) || (point1.Latitude == -(math.Pi/180)*90.) {
		tc = 2 * math.Pi
	} else if point1.Latitude == (math.Pi/180)*90. {
		tc = math.Pi
	} else {
		argacos = (math.Sin(point2.Latitude) - math.Sin(point1.Latitude)*math.Cos(d)) / (math.Sin(d) * math.Cos(point1.Latitude))
		if math.Sin(point2.Longitude-point1.Longitude) < 0 {
			tc = math.Acos(argacos)
		} else {
			tc = 2*math.Pi - math.Acos(argacos)
		}
	}
	return tc
}

/*
IntersectionRadials determines the Coordinate that two Radials
would interset.

Adapted from http://williams.best.vwh.net/avform.htm#Intersection
*/
func IntersectionRadials(radial1, radial2 Radial) (coordinate Coordinate, err error) {
	dst12 := 2 * math.Asin(math.Sqrt(math.Pow((math.Sin((radial1.Latitude-radial2.Latitude)/2)), 2)+
		math.Cos(radial1.Latitude)*math.Cos(radial2.Latitude)*math.Pow(math.Sin((radial1.Longitude-radial2.Longitude)/2), 2)))
	// initial/final bearings between points
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
func ClosestPoint(routeStartCoord, routeEndCoord, actualCoord Coordinate) Coordinate {
	var coordinate Coordinate
	bearing := InitialBearing(routeStartCoord, routeEndCoord)
	distance := AlongTrackDistance(routeStartCoord, routeEndCoord, actualCoord)
	coordinate.Latitude = math.Asin(math.Sin(routeStartCoord.Latitude)*math.Cos(distance) +
		math.Cos(routeStartCoord.Latitude)*math.Sin(distance)*math.Cos(bearing))
	coordinate.Longitude = math.Mod(routeStartCoord.Longitude-math.Asin(math.Sin(bearing)*math.Sin(distance)/math.Cos(coordinate.Latitude))+math.Pi, 2*math.Pi) - math.Pi
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
	/*
		using the helper function above we find the point3 in reach and get the distance
		of to point3. Comparing the expected distance with the provided distance
		the function returns true if point3 is in range, else false
	*/
	distanceBetweenPoints := pointOfReachDistance(point1, point2, point3)
	return distanceBetweenPoints <= distance
}

/*
PointsInReach filters a list of Coordinates to return only those Coordinates that are within testDistance
of the (routeStartCoord, routeEndCoord) route.
*/
func PointsInReach(routeStartCoord, routeEndCoord Coordinate, distance float64, coords []Coordinate) []Coordinate {
	/*
		providing an array of points, the helper function is used to get the distance
		 to those points and then compared with the distance provided.
		 If the points are within distance, they are returned sorted
	*/
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

type MultiPointRoute []NamedCoordinate

func NewMultiPointRoute(coords []NamedCoordinate) (route MultiPointRoute) {
	route = MultiPointRoute(coords)
	return
}

/*
ToSkyVector returns a string that is a valid route for
skyvector.com.
*/
func (route MultiPointRoute) ToSkyVector() (out string) {
	out = ""
	for _, coord := range route {
		if coord.Name != "" {
			out = out + coord.Name
		} else {
			out = out + coord.Coord.ToSkyVector()
		}
		out = out + " "
	}
	return
}

/* MultiPointRoutePOIS takes a 2 lists of coordinates and a distance. The first list of coordinates
will be used to form the multi point route and the second list will be the point of interest list which will be within
the provided distance.
It returns a struct for each match with the point of interest coordinates, the neareast point on the route to the poi and the distance between the nearest poit and the poi.
*/

type MultiPoint struct {
	Poi      Coordinate
	Neareast Coordinate
	Distance float64
}

func MultiPointRoutePOIS(routePoints, pois []Coordinate, distance float64) []MultiPoint {
	var multiPoint []MultiPoint
	for i, point := range routePoints {
		if len(routePoints) > i+1 {
			poistemp := PointsInReach(point, routePoints[i+1], distance, pois)
			fmt.Println(point, routePoints[i+1], distance, pois)
			for _, poi := range poistemp {
				mps := MultiPoint{}
				mps.Poi = poi
				mps.Neareast = ClosestPoint(point, routePoints[i+1], poi)
				mps.Distance = Distance(mps.Neareast, poi)

				// append to the struct
				multiPoint = append(multiPoint, mps)
			}
		}
	}
	// poisInReach can contain duplicate pois, so let's remove the duplicates
	finalPoisInReach := []MultiPoint{}
	m := map[Coordinate]bool{}
	for _, v := range multiPoint {
		if _, seen := m[v.Poi]; !seen {
			finalPoisInReach = append(finalPoisInReach, v)
			m[v.Poi] = true
		}
	}
	return finalPoisInReach
}

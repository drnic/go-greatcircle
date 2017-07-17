package greatcircle

import (
	"fmt"
	"math"
	"testing"
)

var latKSFO, longKSFO = DegreeUnitsToDecimalDegree(37, 37, 00), DegreeUnitsToDecimalDegree(122, 22, 00)
var latKSJC, longKSJC = DegreeUnitsToDecimalDegree(37, 22, 00), DegreeUnitsToDecimalDegree(121, 55, 00)
var latKLAX, longKLAX = DegreeUnitsToDecimalDegree(33, 57, 00), DegreeUnitsToDecimalDegree(118, 24, 00)
var latKJFK, longKJFK = DegreeUnitsToDecimalDegree(40, 38, 00), DegreeUnitsToDecimalDegree(73, 47, 00)
var latKMOD, longKMOD = DegreeUnitsToDecimalDegree(37, 37, 33), DegreeUnitsToDecimalDegree(120, 57, 16)
var latKMAE, longKMAE = DegreeUnitsToDecimalDegree(36, 59, 00), DegreeUnitsToDecimalDegree(120, 7, 00)

var coordKSFO = NamedCoordinate{Coordinate{DegreesToRadians(latKSFO), DegreesToRadians(longKSFO)}, "KSFO"}
var coordKSJC = NamedCoordinate{Coordinate{DegreesToRadians(latKSJC), DegreesToRadians(longKSJC)}, "KSJC"}
var coordKLAX = NamedCoordinate{Coordinate{DegreesToRadians(latKLAX), DegreesToRadians(longKLAX)}, "KLAX"}
var coordKJFK = NamedCoordinate{Coordinate{DegreesToRadians(latKJFK), DegreesToRadians(longKJFK)}, "KJFK"}
var coordKMOD = NamedCoordinate{Coordinate{DegreesToRadians(latKMOD), DegreesToRadians(longKMOD)}, "KMOD"}
var coordKMAE = NamedCoordinate{Coordinate{DegreesToRadians(latKMAE), DegreesToRadians(longKMAE)}, "KMAE"}

var coordsByName = map[string]NamedCoordinate{
	"KSFO": coordKSFO,
	"KSJC": coordKSJC,
	"KLAX": coordKLAX,
	"KJFK": coordKJFK,
	"KMOD": coordKMOD,
	"KMAE": coordKMAE,
}

func TestShowKnownCoordinates(t *testing.T) {
	for name, coord := range coordsByName {
		fmt.Printf("%s: %v\n", name, coord)
	}
}

var degreesRadians = []struct {
	decimaldegrees float64
	radians        float64
}{
	{latKSFO, 0.6565346869585337},
	{longKSFO, 2.135701228023728},
	{latKSJC, 0.6521713638285478},
}

var nauticalMilesRadiansStruct = []struct {
	nauticalMiles float64
	radians       float64
}{
	{99.8665, 0.02905},
}

var distanceStruct = []struct {
	point1Name       string
	point2Name       string
	expectedDistance float64
}{
	{"KMOD", "KMAE", 55.5},
	{"KMAE", "KMOD", 55.5},
	{"KLAX", "KJFK", 2143.7},
	{"KJFK", "KLAX", 2143.7},
}

var initialBearing = []struct {
	point1Name      string
	point2Name      string
	expectedBearing float64
}{
	{"KMOD", "KMAE", DegreesToRadians(133)},
	{"KMAE", "KMOD", DegreesToRadians(314)},
	{"KLAX", "KJFK", DegreesToRadians(66)},
	{"KJFK", "KLAX", DegreesToRadians(274)},
}

var intersectionRadials = []struct {
	radial1 Radial
	radial2 Radial
	point3  Coordinate
	err     string
}{
	{Radial{Coordinate{0.6573, 2.1316}, 1.2392}, Radial{Coordinate{0.6568, 2.1109}, 5.4280}, Coordinate{0.6611492323068847, 2.117252771823951}, ""},
}

var crosstrack = []struct {
	point1   Coordinate
	point2   Coordinate
	point3   Coordinate
	distance float64
}{
	{coordKLAX.Coord, coordKJFK.Coord, Coordinate{0.6021386, 2.033309}, 0.0021674699088520496},
}

var alongtrack = []struct {
	point1   Coordinate
	point2   Coordinate
	point3   Coordinate
	distance float64
}{
// {coordKLAX.Coord, coordKJFK.Coord, Coordinate{0.6021386, 2.033309}, 0.028969025967186944},
}

var pointInReach = []struct {
	point1      Coordinate
	point2      Coordinate
	point3      Coordinate
	distance    float64
	isItInReach bool
}{
	{Coordinate{0.6629, 2.1301}, Coordinate{0.6717, 2.1132}, Coordinate{0.6692, 2.1193}, 30, true},
	{Coordinate{0.6629, 2.1301}, Coordinate{0.6717, 2.1132}, Coordinate{0.6774, 2.1269}, 18, false},
	{Coordinate{0.9427, 0.4892}, Coordinate{0.9593, 0.8124}, Coordinate{0.9595, 0.6364}, 1, false},
}

var pointsInReach = []struct {
	point1       Coordinate
	point2       Coordinate
	points       []Coordinate
	distance     float64
	pointswithin []Coordinate
}{
	{Coordinate{0.6629, 2.1301}, Coordinate{0.6717, 2.1132}, []Coordinate{Coordinate{0.6692, 2.1193}, Coordinate{0.6673, 2.1239}, Coordinate{0.6747, 2.1279}}, 30, []Coordinate{Coordinate{0.6692, 2.1193}, Coordinate{0.6673, 2.1239}}},
	{Coordinate{0.9427, 0.4892}, Coordinate{0.9593, 0.8124}, []Coordinate{Coordinate{0.9595, 0.6364}, Coordinate{0.9654, 0.6665}, Coordinate{1.0075, 0.6750}}, 28, []Coordinate{Coordinate{0.9595, 0.6364}, Coordinate{0.9654, 0.6665}}},
}

var multiPoint = []struct {
	routePoints      []Coordinate
	pois             []Coordinate
	distance         float64
	finalPoisInReach []MultiPoint
}{
	{[]Coordinate{
		Coordinate{latKSFO, longKSFO},
		Coordinate{latKSJC, longKSJC},
		Coordinate{latKLAX, longKLAX},
		Coordinate{latKMAE, longKMAE}},
		[]Coordinate{
			Coordinate{DegreeUnitsToDecimalDegree(36, 46.74, 0), DegreeUnitsToDecimalDegree(120, 8.23, 0)},
			Coordinate{DegreeUnitsToDecimalDegree(37, 20.66, 0), DegreeUnitsToDecimalDegree(121, 36.23, 0)},
			Coordinate{DegreeUnitsToDecimalDegree(34, 59.94, 0), DegreeUnitsToDecimalDegree(120, 19.43, 0)}}, 100,
		[]MultiPoint{
			MultiPoint{Coordinate{DegreeUnitsToDecimalDegree(36, 46.74, 0), DegreeUnitsToDecimalDegree(120, 8.23, 0)}, Coordinate{0, 0}, 7},
			MultiPoint{Coordinate{DegreeUnitsToDecimalDegree(37, 20.66, 0), DegreeUnitsToDecimalDegree(121, 36.23, 0)}, Coordinate{0, 0}, 5},
			MultiPoint{Coordinate{DegreeUnitsToDecimalDegree(34, 59.94, 0), DegreeUnitsToDecimalDegree(120, 19.43, 0)}, Coordinate{0, 0}, 8}},
	},
}

func TestDegreeStringToDegreeUnits(t *testing.T) {
	degree, err := DegreeStrToDecimalDegree("37:37:00")
	if err != nil {
		t.Fatalf("Error parsing 37:37:00; error %v", err)
	}
	if degree != DegreeUnitsToDecimalDegree(37, 37, 0) {
		t.Fatalf("Failed to parse 37:37:00; result %v", DegreeUnitsToDecimalDegree(37, 37, 0))
	}

	degree, err = DegreeStrToDecimalDegree("120:57:16")
	if err != nil {
		t.Fatalf("Error parsing 120:57:16; error %v", err)
	}
	if degree != DegreeUnitsToDecimalDegree(120, 57, 16) {
		t.Fatalf("Failed to parse 120:57:16; result %v", DegreeUnitsToDecimalDegree(120, 57, 16))
	}
}

func TestDegreesToRadians(t *testing.T) {
	for _, v := range degreesRadians {
		result := DegreesToRadians(v.decimaldegrees)
		if result != v.radians {
			t.Fatalf("Expected: %v, received %v", v.radians, result)
		}
	}
}

func TestRadiansToDegrees(t *testing.T) {
	for _, v := range degreesRadians {
		result := RadiansToDegrees(v.radians)
		if math.Abs(result-v.decimaldegrees) > 0.0000001 {
			t.Fatalf("Expected: %v, received %v", v.decimaldegrees, result)
		}
	}
}

func TestNMToRadians(t *testing.T) {
	for _, v := range nauticalMilesRadiansStruct {
		result := NMToRadians(v.nauticalMiles)
		if math.Abs(result-v.radians) > 0.0001 {
			t.Fatalf("Expected: %v, received %v", v.radians, result)
		}
	}
}
func TestRadiansToNM(t *testing.T) {
	for _, v := range nauticalMilesRadiansStruct {
		result := RadiansToNM(v.radians)
		if math.Abs(result-v.nauticalMiles) > 0.0001 {
			t.Fatalf("Expected: %v, received %v", v.nauticalMiles, result)
		}
	}
}

func TestNamedCoordinateEqual(t *testing.T) {
	if !(coordKSFO == coordKSFO) {
		t.Fatalf("Expected %v to equal itself", coordKSFO)
	}
	if coordKSFO == coordKLAX {
		t.Fatalf("Expected %v to not equal %v", coordKSFO, coordKLAX)
	}
}

func TestDistance(t *testing.T) {
	for _, v := range distanceStruct {
		point1, point2 := coordsByName[v.point1Name], coordsByName[v.point2Name]
		result := Distance(point1.Coord, point2.Coord)
		if math.Abs(result-v.expectedDistance) > 0.1 {
			t.Fatalf("Distance between %s %s expected: %v, received %v", v.point1Name, v.point2Name, v.expectedDistance, result)
		}
	}
}

func TestInitialBearing(t *testing.T) {
	for _, v := range initialBearing {
		point1, point2 := coordsByName[v.point1Name], coordsByName[v.point2Name]
		result := InitialBearing(point1.Coord, point2.Coord)
		if math.Abs(result-v.expectedBearing) > 0.5 {
			t.Fatalf("Initial bearing of %s %s expected: %v, received %v", v.point1Name, v.point2Name,
				RadiansToDegrees(v.expectedBearing), RadiansToDegrees(result))
		}
	}
}

func TestIntersection(t *testing.T) {
	for _, v := range intersectionRadials {
		resCoordinate, reserr := IntersectionRadials(v.radial1, v.radial2)
		if resCoordinate != v.point3 && reserr == nil {
			t.Fatalf("Expected: latitude: %v longitude: %v err: %v, received latitude: %v longitude: %v err: %v ", v.point3.Latitude, v.point3.Longitude, v.err, resCoordinate.Latitude, resCoordinate.Longitude, reserr)
		}
	}
}

func TestCrossTrackError(t *testing.T) {
	for _, v := range crosstrack {
		result := CrossTrackError(v.point1, v.point2, v.point3)
		if math.Abs(result-v.distance) > 0.0001 {
			t.Fatalf("Expected: %v, received %v", v.distance, result)
		}
	}
}

func TestAlongTrackDistance(t *testing.T) {
	for _, v := range alongtrack {
		result := AlongTrackDistance(v.point1, v.point2, v.point3)
		if result != v.distance {
			t.Fatalf("Expected: %v, received %v", v.distance, result)
		}
	}
}

func TestClosest(t *testing.T) {
	var closestPoint = []struct {
		point1      Coordinate
		point2      Coordinate
		point3      Coordinate
		coordinates Coordinate
	}{
		{coordKLAX.Coord, coordKJFK.Coord, Coordinate{0.6021386, 2.033309}, Coordinate{0.6041329655944052, 2.034339700924182}},
		{Coordinate{0.6629, 2.1301}, Coordinate{0.6717, 2.1132}, Coordinate{0.6692, 2.1193}, Coordinate{0.6687501299912878, 2.1189029245160818}},
		{Coordinate{0.9427, 0.4892}, Coordinate{0.9593, 0.8124}, Coordinate{0.9595, 0.6364}, Coordinate{0.9565336530696015, 0.6373752108069288}},
	}

	for _, v := range closestPoint {
		result := ClosestPoint(v.point1, v.point2, v.point3)
		if !result.Equal(v.coordinates) {
			t.Fatalf("Expected: %v, received %v", v.coordinates, result)
		}
	}
}

func TestPointInReach(t *testing.T) {
	for _, v := range pointInReach {
		result := PointInReach(v.point1, v.point2, v.point3, v.distance)
		if result != v.isItInReach {
			t.Fatalf("Expected: %v, received %v", v.isItInReach, result)
		}
	}
}

// helper function to find difference in 2 arrays
func difference(slice1 []Coordinate, slice2 []Coordinate) []Coordinate {
	var diff []Coordinate

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func TestPointsInReach(t *testing.T) {
	for _, v := range pointsInReach {
		results := PointsInReach(v.point1, v.point2, v.distance, v.points)
		diff := difference(results, v.pointswithin)
		if len(diff) > 0 {
			t.Fatalf("Expected no extra coordinates, received %v", diff)
		}
		// simpler method to test
		// if len(results) != len(v.pointswithin) {
		// 	t.Fatalf("Expected %v points , received %v", len(results), len(v.pointswithin))
		// }

	}
}

func TestMultiPointRoutePOIS(t *testing.T) {
	for _, v := range multiPoint {
		results := MultiPointRoutePOIS(v.routePoints, v.pois, v.distance)
		fmt.Println("here")
		fmt.Println(results)
		for _, result := range results {
			fmt.Println(result)
		}

	}

}

package gogreatcircle

import (
	"testing"
)

var degreesRadians = []struct {
	decimaldegrees float64
	radians        float64
}{
	{CoordinateToDecimalDegree(37, 37, 00), 0.6565346869585337},
	{CoordinateToDecimalDegree(37, 22, 00), 0.6521713638285478},
	{CoordinateToDecimalDegree(48, 26.57, 00), 0.8454869406615264},
	{CoordinateToDecimalDegree(37, 42.66, 00), 0.6581811142195816},
}

var nauticalMilesRadiansStruct = []struct {
	nauticalMiles float64
	radians       float64
}{
	{5000, 1.454441043328608},
}

var distanceStruct = []struct {
	point1   *Coordinate
	point2   *Coordinate
	distance float64
}{
	{&Coordinate{0.592539, -2.066470}, &Coordinate{0.709186, -1.287762}, 2143.727060139769},
	{&Coordinate{0.65392, -2.13134}, &Coordinate{0.65653, -2.11098}, 56.218067123787776},
}

var initialBearing = []struct {
	point1  *Coordinate
	point2  *Coordinate
	bearing float64
}{
	{&Coordinate{0.592539, -2.066470}, &Coordinate{0.709186, -1.287762}, 1.15003394270832},
	{&Coordinate{0.65392, -2.13134}, &Coordinate{0.65653, -2.11098}, 1.404312223088645},
	{&Coordinate{0.657782598, -2.126090282}, &Coordinate{0.657302632, -2.131588069}, -1.678971437808961},
	{&Coordinate{0.657302632, -2.131588069}, &Coordinate{0.657782598, -2.126090282}, 1.459261107627339},
}

var intersectionRadials = []struct {
	radial1 Radial
	radial2 Radial
	point3  Coordinate
	err     string
}{
	{Radial{Coordinate{0.6573, -2.1316}, 1.2392}, Radial{Coordinate{0.6568, -2.1109}, 5.4280}, Coordinate{0.6611492323068847, -2.117252771823951}, ""},
}

var crosstrack = []struct {
	point1   *Coordinate
	point2   *Coordinate
	point3   *Coordinate
	distance float64
}{
	{&Coordinate{0.592539, -2.066470}, &Coordinate{0.709186, -1.287762}, &Coordinate{0.6021386, -2.033309}, 0.0021674699088520496},
}

var alongtrack = []struct {
	point1   *Coordinate
	point2   *Coordinate
	point3   *Coordinate
	distance float64
}{
	{&Coordinate{0.592539, -2.066470}, &Coordinate{0.709186, -1.287762}, &Coordinate{0.6021386, -2.033309}, 0.028969025967186944},
}

var closestPoint = []struct {
	point1      *Coordinate
	point2      *Coordinate
	point3      *Coordinate
	coordinates Coordinate
}{
	{&Coordinate{0.592539, -2.066470}, &Coordinate{0.709186, -1.287762}, &Coordinate{0.6021386, -2.033309}, Coordinate{0.6041329655944052, -2.034339625370018}},
	{&Coordinate{0.6629, -2.1301}, &Coordinate{0.6717, -2.1132}, &Coordinate{0.6692, -2.1193}, Coordinate{0.6687501299912878, -2.1189211323650383}},
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
		if result != v.decimaldegrees {
			t.Fatalf("Expected: %v, received %v", v.decimaldegrees, result)
		}
	}
}

func TestNMToRadians(t *testing.T) {
	for _, v := range nauticalMilesRadiansStruct {
		result := NMToRadians(v.nauticalMiles)
		if result != v.radians {
			t.Fatalf("Expected: %v, received %v", v.radians, result)
		}
	}
}
func TestRadiansToNM(t *testing.T) {
	for _, v := range nauticalMilesRadiansStruct {
		result := RadiansToNM(v.radians)
		if result != v.nauticalMiles {
			t.Fatalf("Expected: %v, received %v", v.nauticalMiles, result)
		}
	}
}

func TestDistance(t *testing.T) {
	for _, v := range distanceStruct {
		result := Distance(v.point1, v.point2)
		if result != v.distance {
			t.Fatalf("Expected: %v, received %v", v.distance, result)
		}
	}
}

func TestInitialBearing(t *testing.T) {
	for _, v := range initialBearing {
		result := InitialBearing(v.point1, v.point2)
		if result != v.bearing {
			t.Fatalf("Expected: %v, received %v", v.bearing, result)
		}
	}
}

func TestIntersection(t *testing.T) {
	for _, v := range intersectionRadials {
		resCoordinate, reserr := IntersectionRadials(&v.radial1, &v.radial2)
		if resCoordinate != v.point3 && reserr == nil {
			t.Fatalf("Expected: latitude: %v longitude: %v err: %v, received latitude: %v longitude: %v err: %v ", v.point3.latitude, v.point3.longitude, v.err, resCoordinate.latitude, resCoordinate.longitude, reserr)
		}
	}
}

func TestCrossTrackError(t *testing.T) {
	for _, v := range crosstrack {
		result := CrossTrackError(v.point1, v.point2, v.point3)
		if result != v.distance {
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
	for _, v := range closestPoint {
		result := ClosestPoint(v.point1, v.point2, v.point3)
		if result != v.coordinates {
			t.Fatalf("Expected: %v, received %v", v.coordinates, result)
		}
	}
}

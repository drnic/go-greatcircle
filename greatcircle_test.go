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
	lat1     float64
	lon1     float64
	lat2     float64
	lon2     float64
	distance float64
}{
	{0.592539, -2.066470, 0.709186, -1.287762, 2143.727060139769},
	{0.65392, -2.13134, 0.65653, -2.11098, 56.218067123787776},
}

var initialBearing = []struct {
	lat1    float64
	lon1    float64
	lat2    float64
	lon2    float64
	bearing float64
}{
	{0.592539, -2.066470, 0.709186, -1.287762, 1.15003394270832},
	{0.65392, -2.13134, 0.65653, -2.11098, 1.404312223088645},
	{0.657782598, -2.126090282, 0.657302632, -2.131588069, -1.678971437808961},
	{0.657302632, -2.131588069, 0.657782598, -2.126090282, 1.459261107627339},
}

var intersectionRadials = []struct {
	lat1     float64
	lon1     float64
	bearing1 float64
	lat2     float64
	lon2     float64
	bearing2 float64
	lat3     float64
	lon3     float64
	err      string
}{
	{0.6573, -2.1316, 1.2392, 0.6568, -2.1109, 5.4280, 0.6611492323068847, -2.117252771823951, ""},
}

var crosstrack = []struct {
	lat1     float64
	lon1     float64
	lat2     float64
	lon2     float64
	lat3     float64
	lon3     float64
	distance float64
}{
	{0.592539, -2.066470, 0.709186, -1.287762, 0.6021386, -2.033309, 0.0021674699088520496},
}

var alongtrack = []struct {
	lat1     float64
	lon1     float64
	lat2     float64
	lon2     float64
	lat3     float64
	lon3     float64
	distance float64
}{
	{0.592539, -2.066470, 0.709186, -1.287762, 0.6021386, -2.033309, 0.005594254069336081},
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
		result := Distance(v.lat1, v.lon1, v.lat2, v.lon2)
		if result != v.distance {
			t.Fatalf("Expected: %v, received %v", v.distance, result)
		}
	}
}

func TestInitialBearing(t *testing.T) {
	for _, v := range initialBearing {
		result := InitialBearing(v.lat1, v.lon1, v.lat2, v.lon2)
		if result != v.bearing {
			t.Fatalf("Expected: %v, received %v", v.bearing, result)
		}
	}
}

func TestIntersection(t *testing.T) {
	for _, v := range intersectionRadials {
		reslat3, reslon3, reserr := IntersectionRadials(v.lat1, v.lon1, v.bearing1, v.lat2, v.lon2, v.bearing2)
		if reslat3 != v.lat3 && reslon3 != v.lon3 && reserr == nil {
			t.Fatalf("Expected: lat3: %v lon3: %v err: %v, received lat3: %v lon3: %v err: %v ", v.lat3, v.lon3, v.err, reslat3, reslon3, reserr)
		}
	}
}

func TestCrossTrackError(t *testing.T) {
	for _, v := range crosstrack {
		result := CrossTrackError(v.lat1, v.lon1, v.lat2, v.lon2, v.lat3, v.lon3)
		if result != v.distance {
			t.Fatalf("Expected: %v, received %v", v.distance, result)
		}
	}
}

func TestAlongTrackDistance(t *testing.T) {
	for _, v := range alongtrack {
		result := AlongTrackDistance(v.lat1, v.lon1, v.lat2, v.lon2, v.lat3, v.lon3)
		if result != v.distance {
			t.Fatalf("Expected: %v, received %v", v.distance, result)
		}
	}
}

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
	{0.592539, -2.066470, 0.709186, -1.287762, 65.89209121397664},
	{0.65392, -2.13134, 0.65653, -2.11098, 80.46116350161344},
	{0.657782598, -2.126090282, 0.657302632, -2.131588069, 263.80202269053495},
	{0.657302632, -2.131588069, 0.657782598, -2.126090282, 83.60950267463232},
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

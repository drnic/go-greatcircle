package gogreatcircle

import (
	"errors"
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

func approxDistance(lat1, lon1, lat2, lon2 float64) float64 {
	w := lon2 - lon1
	v := lat1 - lat2
	s := 2 * math.Asin(math.Sqrt((math.Sin(v/2)*math.Sin(v/2))+(math.Cos(lat1)*math.Cos(lat2)*math.Sin(w/2)*math.Sin(w/2))))
	return s
}

func modlon(x float64) float64 { //ensure longitude is +/-180
	return math.Mod(x+math.Pi, 2*math.Pi) - math.Pi
}

func CoordinateToDecimalDegree(degrees, minutes, seconds float64) float64 {
	return degrees + (minutes / 60) + (seconds / 3600)
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func RadiansToDegrees(radians float64) float64 {
	return radians * (180 / math.Pi)
}

func NMToRadians(nauticalMiles float64) float64 {
	return (math.Pi / (180 * 60)) * nauticalMiles
}

func RadiansToNM(radians float64) float64 {
	return ((180 * 60) / math.Pi) * radians
}

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	return (math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon1-lon2)) * 180 * 60) / math.Pi

}

func InitialBearing(lat1, lon1, lat2, lon2 float64) float64 {
	dLon := (lon2 - lon1)
	y := math.Sin(dLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLon)
	// bearing in radians
	bearing := math.Atan2(y, x)
	return bearing
}

func IntersectionRadials(lat1, lon1, bearing1, lat2, lon2, bearing2 float64) (lat3, lon3 float64, err error) {
	// adapted from http://williams.best.vwh.net/avform.htm#Intersection
	dLat := lat2 - lat1
	dLon := lon2 - lon1

	dist12 := 2 * math.Asin(math.Sqrt(math.Sin(dLat/2)*math.Sin(dLat/2)+math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)))
	if dist12 == 0 {
		return 0, 0, errors.New("dist 0")
	}

	// initial/final bearings between points
	brngA := math.Acos((math.Sin(lat2) - math.Sin(lat1)*math.Cos(dist12)) / (math.Sin(dist12) * math.Cos(lat1)))
	brngB := math.Acos((math.Sin(lat1) - math.Sin(lat2)*math.Cos(dist12)) / (math.Sin(dist12) * math.Cos(lat2)))

	var brng12 float64
	var brng21 float64
	if math.Sin(lon2-lon1) > 0 {
		brng12 = brngA
		brng21 = 2*math.Pi - brngB
	} else {
		brng12 = 2*math.Pi - brngA
		brng21 = brngB
	}

	alpha1 := math.Mod((bearing1-brng12+math.Pi), (2*math.Pi)) - math.Pi // angle 2-1-3
	alpha2 := math.Mod((brng21-bearing2+math.Pi), (2*math.Pi)) - math.Pi // angle 1-2-3

	if math.Sin(alpha1) == 0 && math.Sin(alpha2) == 0 {
		return 0, 0, errors.New("infinite intersections")
	}
	if math.Sin(alpha1)*math.Sin(alpha2) < 0 {
		return 0, 0, errors.New("ambiguous intersection")
	}

	alpha3 := math.Acos(-math.Cos(alpha1)*math.Cos(alpha2) + math.Sin(alpha1)*math.Sin(alpha2)*math.Cos(dist12))
	dist13 := math.Atan2(math.Sin(dist12)*math.Sin(alpha1)*math.Sin(alpha2), math.Cos(alpha2)+math.Cos(alpha1)*math.Cos(alpha3))
	lat3 = math.Asin(math.Sin(lat1)*math.Cos(dist13) + math.Cos(lat1)*math.Sin(dist13)*math.Cos(bearing1))
	dLon13 := math.Atan2(math.Sin(bearing1)*math.Sin(dist13)*math.Cos(lat1),
		math.Cos(dist13)-math.Sin(lat1)*math.Sin(lat3))
	lon3 = lon1 + dLon13
	lon3 = math.Mod((lon3+3*math.Pi), (2*math.Pi)) - math.Pi // normalise to -180..+180º

	return lat3, lon3, nil

}

func CrossTrackError(lat1, lon1, lat2, lon2, lat3, lon3 float64) float64 {
	dist_AD := NMToRadians(Distance(lat1, lon1, lat3, lon3))
	crs_AD := math.Acos((math.Sin(lat3) - math.Sin(lat1)*math.Cos(dist_AD)) / (math.Sin(dist_AD) * math.Cos(lat1)))
	initialBearing := InitialBearing(lat1, lon1, lat2, lon2)
	xtd := math.Asin(math.Sin(dist_AD) * math.Sin(crs_AD-initialBearing))
	return xtd
}

func AlongTrackDistance(lat1, lon1, lon2, lat2, lat3, lon3 float64) float64 {
	dist_AD := NMToRadians(Distance(lat1, lon1, lat3, lon3))
	xtd := CrossTrackError(lat1, lon1, lat2, lon2, lat3, lon3)
	atd := math.Asin(math.Sqrt(math.Pow((math.Sin(dist_AD)), 2)-math.Pow((math.Sin(xtd)), 2)) / math.Cos(xtd))
	return atd
}

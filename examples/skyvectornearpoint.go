package main

import (
	"fmt"

	"github.com/drnic/go-greatcircle"
)

func coord(name string, latitude string, longitude string) greatcircle.NamedCoordinate {
	nc, err := greatcircle.NewNamedCoordinate(name, latitude, longitude)
	if err != nil {
		panic(nil)
	}
	return nc
}

func main() {
	var coordsByName = map[string]greatcircle.NamedCoordinate{
		"KSFO": coord("KSFO", "37:37:00", "122:22:00"),
		"KSJC": coord("KSJC", "37:22:00", "121:55:00"),
		"E16":  coord("E16", "37:5", "121:35:20"),
		"KLAX": coord("KLAX", "33:57:00", "118:24:00"),
		"KJFK": coord("KJFK", "40:38:00", "73:47:00"),
		"KMOD": coord("KMOD", "37:37:33", "120:57:16"),
		"KMAE": coord("KMAE", "36:59:00", "120:7:00"),
		"KKIC": coord("KKIC", "36:13:50", "121:7:00"),
	}

	ksjcOnRoute := greatcircle.ClosestPoint(coordsByName["KSFO"].Coord, coordsByName["KLAX"].Coord, coordsByName["KSJC"].Coord)
	e16OnRoute := greatcircle.ClosestPoint(coordsByName["KSFO"].Coord, coordsByName["KLAX"].Coord, coordsByName["E16"].Coord)

	route := greatcircle.NewMultiPointRoute([]greatcircle.NamedCoordinate{
		coordsByName["KSFO"],
		ksjcOnRoute.ToNamedCoordinate(),
		e16OnRoute.ToNamedCoordinate(),
		coordsByName["KLAX"],
	})
	fmt.Println("Some selected closest points along a route:")
	fmt.Println("Visit http://skyvector.com/ and enter the following flight plan:")
	fmt.Println(route.ToSkyVector())

	coords := []greatcircle.Coordinate{}
	for _, coord := range coordsByName {
		coords = append(coords, coord.Coord)
	}
	coordsInReach := greatcircle.PointsInReach(coordsByName["KSFO"].Coord, coordsByName["KLAX"].Coord, 25, coords)

	route = greatcircle.NewMultiPointRoute([]greatcircle.NamedCoordinate{coordsByName["KSFO"]})
	for _, coordInReach := range coordsInReach {
		if !coordInReach.Equal(coordsByName["KLAX"].Coord) {
			route = append(route, coordInReach.ToNamedCoordinate())
		}
	}
	route = append(route, coordsByName["KLAX"])

	fmt.Println("\nA route of points that are within 25nM of KSFO-KLAX:")
	fmt.Println("Visit http://skyvector.com/ and enter the following flight plan:")
	fmt.Println(route.ToSkyVector())

	route = greatcircle.NewMultiPointRoute([]greatcircle.NamedCoordinate{coordsByName["KSFO"]})
	for _, coordInReach := range coordsInReach {
		if !coordInReach.Equal(coordsByName["KLAX"].Coord) {
			closestPoint := greatcircle.ClosestPoint(coordsByName["KSFO"].Coord, coordsByName["KLAX"].Coord, coordInReach)
			route = append(route, closestPoint.ToNamedCoordinate())
		}
	}
	route = append(route, coordsByName["KLAX"])

	fmt.Println("\nThe route KSFO-KLAX with waypoints near other airports along the route:")
	fmt.Println("Visit http://skyvector.com/ and enter the following flight plan:")
	fmt.Println(route.ToSkyVector())

}

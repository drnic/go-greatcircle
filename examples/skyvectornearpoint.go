package main

import (
	"fmt"

	gc "github.com/drnic/greatcircle"
)

func coord(name string, latitude string, longitude string) gc.NamedCoordinate {
	nc, err := gc.NewNamedCoordinate(name, latitude, longitude)
	if err != nil {
		panic(nil)
	}
	return nc
}

func main() {
	var coordsByName = map[string]gc.NamedCoordinate{
		"KSFO": coord("KSFO", "37:37:00", "122:22:00"),
		"KSJC": coord("KSJC", "37:22:00", "121:55:00"),
		"E16":  coord("E16", "37:5", "121:35:20"),
		"KLAX": coord("KLAX", "33:57:00", "118:24:00"),
		"KJFK": coord("KJFK", "40:38:00", "73:47:00"),
		"KMOD": coord("KMOD", "37:37:33", "120:57:16"),
		"KMAE": coord("KMAE", "36:59:00", "120:7:00"),
	}

	ksjcOnRoute := gc.ClosestPoint(coordsByName["KSFO"].Coord, coordsByName["KLAX"].Coord, coordsByName["KSJC"].Coord)
	e16OnRoute := gc.ClosestPoint(coordsByName["KSFO"].Coord, coordsByName["KLAX"].Coord, coordsByName["E16"].Coord)

	route := gc.NewMultiPointRoute([]gc.NamedCoordinate{
		coordsByName["KSFO"],
		ksjcOnRoute.ToNamedCoordinate(),
		e16OnRoute.ToNamedCoordinate(),
		coordsByName["KLAX"],
	})
	fmt.Println("Visit http://skyvector.com/ and enter the following flight plan:")
	fmt.Println(route.ToSkyVector())
}

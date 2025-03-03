package dungeon

import "codeberg.org/anaseto/gruid/rl"

type cell rl.Cell

// These constants represent the different kind of map tiles.
const (
	WallCell cell = iota
	FloorCell
	DoorCell
	RoadCell
	FoliageCell

	Explored = 0b10000000
)

func terrain(c cell) cell {
	return c &^ Explored
}

func explored(c cell) bool {
	return c&Explored != 0
}

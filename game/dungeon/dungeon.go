// This file contains map-related code.

package dungeon

import (
	"math/rand"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/rl"
)

// These constants represent the different kind of map tiles.
const (
	WallCell rl.Cell = iota
	FloorCell
	DoorCell
	RoadCell
	FoliageCell
)

// Map represents the rectangular map of the game's level.
type Map struct {
	Grid     rl.Grid
	rand     *rand.Rand           // random number generator
	Explored map[gruid.Point]bool // explored cells
}

func (m *Map) GetRand() *rand.Rand {
	return m.rand
}

func (m *Map) SeedRand(seed int64) {
	m.rand = rand.New(rand.NewSource(seed))
}

func (m *Map) At(p gruid.Point) rl.Cell {
	return m.Grid.At(p)
}

// Walkable returns true if at the given position there is a floor tile.
func (m *Map) Walkable(p gruid.Point) bool {
	return m.Grid.At(p) == FloorCell
}

// Rune returns the character rune representing a given terrain.
func (m *Map) Rune(c rl.Cell) (r rune) {
	switch c {
	case WallCell:
		r = '#'
	case FloorCell:
		r = '.'
	case DoorCell:
		r = '+'
	case RoadCell:
		r = '='
	case FoliageCell:
		r = '"'
	}
	return r
}

// RandomFloor returns a random floor cell in the map. It assumes that such a
// floor cell exists (otherwise the function does not end).
func (m *Map) RandomFloor() gruid.Point {
	size := m.Grid.Size()
	for {
		freep := gruid.Point{m.rand.Intn(size.X), m.rand.Intn(size.Y)}
		if m.Grid.At(freep) == FloorCell {
			return freep
		}
	}
}

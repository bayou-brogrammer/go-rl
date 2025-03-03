// This file contains map-related code.

package dungeon

import (
	"math/rand"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/rl"
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

func (m *Map) Cell(p gruid.Point) cell {
	return cell(m.Grid.At(p))
}

func (m *Map) SetCell(p gruid.Point, c cell) {
	oc := m.Cell(p)
	m.Grid.Set(p, rl.Cell(c|oc&Explored))
}

func (m *Map) SetExplored(p gruid.Point) {
	oc := m.Cell(p)
	m.Grid.Set(p, rl.Cell(oc|Explored))
}

func (m *Map) At(p gruid.Point) rl.Cell {
	return m.Grid.At(p)
}

// Walkable returns true if at the given position there is a floor tile.
func (m *Map) Walkable(p gruid.Point) bool {
	return m.Cell(p) == FloorCell
}

// Rune returns the character rune representing a given terrain.
func (m *Map) Rune(c rl.Cell) (r rune) {
	switch cell(c) {
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
		if m.Cell(freep) == FloorCell {
			return freep
		}
	}
}

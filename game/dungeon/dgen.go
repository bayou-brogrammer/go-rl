package dungeon

import (
	"fmt"
	"math/rand"
	"time"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
	"codeberg.org/anaseto/gruid/rl"
)

type maplayout int

const (
	Town maplayout = iota
	AutomataCave
	// RandomWalkCave
	// RandomWalkTreeCave
	// RandomSmallWalkCaveUrbanised
	// NaturalCave
)

type dgen struct {
	m *Map

	layout maplayout

	rand *rand.Rand
	// neighbors paths.Neighbors
	PR *paths.PathRange
}

func NewMap(size gruid.Point, ml maplayout) *Map {
	m := &Map{
		Grid:     rl.NewGrid(size.X, size.Y),
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		Explored: make(map[gruid.Point]bool),
	}

	d := &dgen{
		m:      m,
		layout: ml,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
		PR:     paths.NewPathRange(m.Grid.Range()),
	}

	switch ml {
	case Town:
		d.GenerateTown()
	case AutomataCave:
		d.Generate()
	}

	fmt.Println("Done generating town")
	return m
}

// Generate fills the Grid attribute of m with a procedurally generated map.
func (dg *dgen) Generate() {
	// map generator using the rl package from gruid
	mgen := rl.MapGen{Rand: dg.rand, Grid: dg.m.Grid}
	// cellular automata map generation with rules that give a cave-like
	// map.
	rules := []rl.CellularAutomataRule{
		{WCutoff1: 5, WCutoff2: 2, Reps: 4, WallsOutOfRange: true},
		{WCutoff1: 5, WCutoff2: 25, Reps: 3, WallsOutOfRange: true},
	}

	for {
		mgen.CellularAutomataCave(rl.Cell(WallCell), rl.Cell(FloorCell), 0.42, rules)
		freep := dg.m.RandomFloor()
		// We put walls in floor cells non reachable from freep, to ensure that
		// all the cells are connected (which is not guaranteed by cellular
		// automata map generation).
		pr := paths.NewPathRange(dg.m.Grid.Range())
		pr.CCMap(&path{m: dg.m}, freep)
		ntiles := mgen.KeepCC(pr, freep, rl.Cell(WallCell))
		const minCaveSize = 400
		if ntiles > minCaveSize {
			break
		}
		// If there were not enough free tiles, we run the map
		// generation again.
	}
}

// Foliage generates foliage across the entire map
func (dg *dgen) Foliage(less bool) {
	// Call FoliageInRange with the entire map range
	dg.FoliageInRange(less, dg.m.Grid.Range())
}

// FoliageInRange generates foliage only within the specified range of the map
func (dg *dgen) FoliageInRange(less bool, rg gruid.Range) {
	// Create a grid with the same size as the range
	width := rg.Max.X - rg.Min.X
	height := rg.Max.Y - rg.Min.Y
	gd := rl.NewGrid(width, height)

	rules := []rl.CellularAutomataRule{
		{WCutoff1: 5, WCutoff2: 2, Reps: 4, WallsOutOfRange: true},
		{WCutoff1: 5, WCutoff2: 25, Reps: 2, WallsOutOfRange: true},
	}
	mg := rl.MapGen{Rand: dg.rand, Grid: gd}
	winit := 0.55
	if less {
		winit = 0.53
	}

	mg.CellularAutomataCave(rl.Cell(WallCell), rl.Cell(FoliageCell), winit, rules)

	// Iterate over the specified range in the main grid
	for y := rg.Min.Y; y < rg.Max.Y; y++ {
		for x := rg.Min.X; x < rg.Max.X; x++ {
			p := gruid.Point{X: x, Y: y}
			// Calculate the corresponding position in the generated grid
			gp := gruid.Point{X: x - rg.Min.X, Y: y - rg.Min.Y}

			// Only place foliage on floor cells
			if dg.m.Cell(p) == FloorCell && gd.At(gp) == rl.Cell(FoliageCell) {
				dg.m.SetCell(p, FoliageCell)
			}
		}
	}
}

// path implements the paths.Pather interface and is used to provide pathing
// information in map generation.
type path struct {
	m  *Map
	nb paths.Neighbors
}

// Neighbors returns the list of walkable neighbors of q in the map using 4-way
// movement along cardinal directions.
func (p *path) Neighbors(q gruid.Point) []gruid.Point {
	return p.nb.Cardinal(q,
		func(r gruid.Point) bool { return p.m.Walkable(r) })
}

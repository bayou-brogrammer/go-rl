package main

import (
	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
)

// aiPath implements the paths.Astar interface for use in AI pathfinding.
type aiPath struct {
	g  *game
	nb paths.Neighbors
}

// Neighbors returns the list of walkable neighbors of q in the map using 4-way
// movement along cardinal directions.
func (aip *aiPath) Neighbors(q gruid.Point) []gruid.Point {
	return aip.nb.Cardinal(q,
		func(r gruid.Point) bool {
			return aip.g.Map.Walkable(r)
		})
}

// Cost implements paths.Astar.Cost.
func (aip *aiPath) Cost(p, q gruid.Point) int {
	if !aip.g.ECS.NoBlockingEntityAt(q) {
		// Extra cost for blocked positions: this encourages the
		// pathfinding algorithm to take another path to reach the
		// player.
		return 8
	}
	return 1
}

// Estimation implements paths.Astar.Estimation. For 4-way movement, we use the
// Manhattan distance.
func (aip *aiPath) Estimation(p, q gruid.Point) int {
	return paths.DistanceManhattan(p, q)
}

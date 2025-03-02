// game/dungeon/bsp.go
package dungeon

import (
	"fmt"
	"math/rand"

	"codeberg.org/anaseto/gruid"
)

// BSPRoom represents a room in the dungeon
type BSPRoom struct {
	X, Y          int // Top-left corner
	Width, Height int
}

func (r *BSPRoom) Carve(m *Map) {
	for x := r.X; x < r.X+r.Width; x++ {
		for y := r.Y; y < r.Y+r.Height; y++ {
			m.Grid.Set(gruid.Point{X: x, Y: y}, Floor)
		}
	}
}

// BSPNode represents a node in the BSP tree
type BSPNode struct {
	X, Y          int
	Width, Height int
	Room          *BSPRoom
	Left, Right   *BSPNode
}

// CreateRoom creates a room within the given node with some padding
func (node *BSPNode) CreateRoom(rng *rand.Rand) {
	// Minimum room size
	minSize := 4
	// Padding from the node edges
	padding := 1

	// Calculate maximum possible room dimensions
	maxWidth := node.Width - (padding * 2)
	maxHeight := node.Height - (padding * 2)

	if maxWidth < minSize || maxHeight < minSize {
		return
	}

	// Generate random room dimensions
	width := minSize + rng.Intn(maxWidth-minSize+1)
	height := minSize + rng.Intn(maxHeight-minSize+1)

	// Calculate room position (centered in node)
	rWidth := node.Width - width - (padding * 2)
	rHeight := node.Height - height - (padding * 2)
	if rWidth <= 0 || rHeight <= 0 {
		return
	}

	fmt.Printf("rWidth: %d, rHeight: %d\n", rWidth, rHeight)

	x := node.X + padding + rng.Intn(rWidth)
	y := node.Y + padding + rng.Intn(rHeight)

	node.Room = &BSPRoom{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// SplitNode recursively splits the node into two child nodes
func (node *BSPNode) SplitNode(rng *rand.Rand, minSize int, depth int) {
	if depth <= 0 {
		return
	}

	// Decide split direction (horizontal or vertical)
	horizontal := rng.Float64() < 0.5

	if horizontal {
		if node.Height < minSize*2 {
			return
		}
		// Calculate split point
		split := minSize + rng.Intn(node.Height-minSize*2)

		node.Left = &BSPNode{
			X:      node.X,
			Y:      node.Y,
			Width:  node.Width,
			Height: split,
		}
		node.Right = &BSPNode{
			X:      node.X,
			Y:      node.Y + split,
			Width:  node.Width,
			Height: node.Height - split,
		}
	} else {
		if node.Width < minSize*2 {
			return
		}

		// Calculate split point
		minWidth := node.Width - minSize*2
		if minWidth <= 0 {
			return
		}

		split := minSize + rng.Intn(minWidth)

		node.Left = &BSPNode{
			X:      node.X,
			Y:      node.Y,
			Width:  split,
			Height: node.Height,
		}
		node.Right = &BSPNode{
			X:      node.X + split,
			Y:      node.Y,
			Width:  node.Width - split,
			Height: node.Height,
		}
	}

	// Recursively split children
	node.Left.SplitNode(rng, minSize, depth-1)
	node.Right.SplitNode(rng, minSize, depth-1)
}

// ConnectRooms creates corridors between rooms
func (node *BSPNode) ConnectRooms(m *Map) {
	if node.Left == nil || node.Right == nil {
		return
	}

	// Get centers of rooms in left and right nodes
	var leftCenter, rightCenter gruid.Point
	if node.Left.Room != nil {
		leftCenter = gruid.Point{
			X: node.Left.Room.X + node.Left.Room.Width/2,
			Y: node.Left.Room.Y + node.Left.Room.Height/2,
		}
	} else if node.Left.Left != nil && node.Left.Left.Room != nil {
		room := node.Left.Left.Room
		leftCenter = gruid.Point{
			X: room.X + room.Width/2,
			Y: room.Y + room.Height/2,
		}
	}

	if node.Right.Room != nil {
		rightCenter = gruid.Point{
			X: node.Right.Room.X + node.Right.Room.Width/2,
			Y: node.Right.Room.Y + node.Right.Room.Height/2,
		}
	} else if node.Right.Right != nil && node.Right.Right.Room != nil {
		room := node.Right.Right.Room
		rightCenter = gruid.Point{
			X: room.X + room.Width/2,
			Y: room.Y + room.Height/2,
		}
	}

	// Create L-shaped corridor
	createCorridor(m, leftCenter, rightCenter)

	// Recursively connect rooms in children
	node.Left.ConnectRooms(m)
	node.Right.ConnectRooms(m)
}

// createCorridor creates an L-shaped corridor between two points
func createCorridor(m *Map, start, end gruid.Point) {
	// Create horizontal corridor
	x1, x2 := start.X, end.X
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	for x := x1; x <= x2; x++ {
		m.Grid.Set(gruid.Point{X: x, Y: start.Y}, Floor)
	}

	// Create vertical corridor
	y1, y2 := start.Y, end.Y
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for y := y1; y <= y2; y++ {
		m.Grid.Set(gruid.Point{X: end.X, Y: y}, Floor)
	}
}

func (m *Map) GenerateBSPMap() {
	fmt.Println("generating bsp map")

	// Fill the map with walls initially
	size := m.Grid.Size()
	for x := range size.X {
		for y := range size.Y {
			m.Grid.Set(gruid.Point{X: x, Y: y}, Wall)
		}
	}

	// Create root node covering the whole map
	root := &BSPNode{
		X:      0,
		Y:      0,
		Width:  size.X,
		Height: size.Y,
	}

	// Split the dungeon recursively
	minSize := 6  // Minimum size of a partition
	maxDepth := 5 // Maximum depth of BSP tree
	root.SplitNode(m.rand, minSize, maxDepth)

	// Create rooms in the leaf nodes
	var createRooms func(*BSPNode)
	createRooms = func(node *BSPNode) {
		if node.Left == nil && node.Right == nil {
			node.CreateRoom(m.rand)
			if node.Room != nil {
				// Carve out the room
				node.Room.Carve(m)
			}
			return
		}
		if node.Left != nil {
			createRooms(node.Left)
		}
		if node.Right != nil {
			createRooms(node.Right)
		}
	}
	createRooms(root)

	// Connect the rooms with corridors
	root.ConnectRooms(m)

	// Verify connectivity (using existing code)
	// freep := m.RandomFloor()
	// pr := paths.NewPathRange(m.Grid.Range())
	// pr.CCMap(&path{m: m}, freep)
	// ntiles := 0
	// m.Grid.Range().Iter(func(p gruid.Point) {
	// 	if m.Grid.At(p) == Floor {
	// 		ntiles++
	// 	}
	// })

	// If the map isn't well connected, generate a new one
	// const minRoomSize = 100
	// if ntiles < minRoomSize {
	// 	m.Generate()
	// }
}

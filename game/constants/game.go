package constants

import "codeberg.org/anaseto/gruid"

var Version string = "v0.1.0"

var (
	CustomKeys             = false
	DisableAnimations bool = false
	LogGame                = false
	Terminal               = false
	Xterm256Color          = false
	Testing                = true
)

const (
	// UI constants
	UIWidth  = 80
	UIHeight = 24
	LogLines = 2

	// Dungeon dimensions
	MapWidth  = UIWidth
	MapHeight = UIHeight - 1 - LogLines
	MapNCells = MapWidth * MapHeight

	// Win constants
	WinDepth = 8
	MaxDepth = 11

	// Vision constants
	TreeRange              = 50
	LightRange             = 6
	DefaultLOSRange        = 12
	DefaultMonsterLOSRange = 12
)

var InvalidPos = gruid.Point{-1, -1}

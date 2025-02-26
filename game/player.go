package main

import (
	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/rl"
)

// Player contains information relevant to the player.
type Player struct {
	FOV *rl.FOV // player's field of view
}

// maxLOS is the maximum distance in player's field of view.
const maxLOS = 10

// NewPlayer returns a new Player entity at a given position.
func NewPlayer() *Player {
	player := &Player{}
	player.FOV = rl.NewFOV(gruid.NewRange(-maxLOS, -maxLOS, maxLOS+1, maxLOS+1))
	return player
}

package main

import "codeberg.org/anaseto/gruid"

// Color definitions. We start from 1, because 0 is gruid.ColorDefault, which
// we use for default foreground and background.
const (
	ColorFOV gruid.Color = iota + 1
	ColorPlayer
	ColorMonster
	ColorLogPlayerAttack
	ColorLogItemUse
	ColorLogMonsterAttack
	ColorLogSpecial
	ColorStatusHealthy
	ColorStatusWounded
	ColorConsumable
	ColorMenuActive
)

const (
	AttrReverse = 1 << iota
)

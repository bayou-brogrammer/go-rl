package game

import "codeberg.org/anaseto/gruid"

var NORMAL_KEYS = map[gruid.Key]action{
	gruid.KeyArrowLeft:  {Type: ActionBump, Delta: gruid.Point{-1, 0}},
	gruid.KeyArrowDown:  {Type: ActionBump, Delta: gruid.Point{0, 1}},
	gruid.KeyArrowUp:    {Type: ActionBump, Delta: gruid.Point{0, -1}},
	gruid.KeyArrowRight: {Type: ActionBump, Delta: gruid.Point{1, 0}},
	// "h":                 ActionW,
	// "j":                 ActionS,
	// "k":                 ActionN,
	// "l":                 ActionE,
	// "a":                 ActionW,
	// "s":                 ActionS,
	// "w":                 ActionN,
	// "d":                 ActionE,
	"x": {Type: ActionExamine},
}

var TARGET_KEYS = map[gruid.Key]action{
	// gruid.KeyArrowLeft:  ActionW,
	// gruid.KeyArrowDown:  ActionS,
	// gruid.KeyArrowUp:    ActionN,
	// gruid.KeyArrowRight: ActionE,
	// "h":                 ActionW,
	// "j":                 ActionS,
	// "k":                 ActionN,
	// "l":                 ActionE,
	// "a":                 ActionW,
	// "s":                 ActionS,
	// "w":                 ActionN,
	// "d":                 ActionE,
}

func (md *model) initKeys() {
	md.keysNormal = NORMAL_KEYS
	md.keysTarget = TARGET_KEYS
}

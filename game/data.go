package game

import "codeberg.org/anaseto/gruid"

// case gruid.KeyEnter, ".":
// 	m.action = action{Type: ActionWait}
// case "Q":
// 	m.action = action{Type: ActionQuit}
// case "S":
// 	m.action = action{Type: ActionSave}
// case "m":
// 	m.action = action{Type: ActionViewMessages}
// case "i":
// 	m.action = action{Type: ActionInventory}
// case "d":
// 	m.action = action{Type: ActionDrop}
// case "g":
// 	m.action = action{Type: ActionPickup}
// case "x":
// 	m.action = action{Type: ActionExamine}

var NORMAL_KEYS = map[gruid.Key]action{
	gruid.KeyArrowLeft:  {Type: ActionBump, Delta: gruid.Point{-1, 0}},
	gruid.KeyArrowDown:  {Type: ActionBump, Delta: gruid.Point{0, 1}},
	gruid.KeyArrowUp:    {Type: ActionBump, Delta: gruid.Point{0, -1}},
	gruid.KeyArrowRight: {Type: ActionBump, Delta: gruid.Point{1, 0}},

	"h": {Type: ActionBump, Delta: gruid.Point{-1, 0}},
	"j": {Type: ActionBump, Delta: gruid.Point{0, 1}},
	"k": {Type: ActionBump, Delta: gruid.Point{0, -1}},
	"l": {Type: ActionBump, Delta: gruid.Point{1, 0}},
	"a": {Type: ActionBump, Delta: gruid.Point{-1, 0}},
	"s": {Type: ActionBump, Delta: gruid.Point{0, 1}},
	"w": {Type: ActionBump, Delta: gruid.Point{0, -1}},

	"d": {Type: ActionDrop},
	"g": {Type: ActionPickup},
	"x": {Type: ActionExamine},

	"i": {Type: ActionInventory},
	"m": {Type: ActionViewMessages},

	"Q": {Type: ActionQuit},
	"S": {Type: ActionSave},
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

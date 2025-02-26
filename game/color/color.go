package color

import (
	"image/color"

	"codeberg.org/anaseto/gruid"
)

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

func GetCellColor(c gruid.Cell) (color.RGBA, color.RGBA) {
	// We use some colors from https://github.com/jan-warchol/selenized,
	// using the palette variant with dark backgound and light foreground.
	fg := color.RGBA{0xad, 0xbc, 0xbc, 255}
	bg := color.RGBA{0x10, 0x3c, 0x48, 255}

	// We define non default-colors (for FOV, ...).
	switch c.Style.Bg {
	case ColorFOV:
		bg = color.RGBA{0x18, 0x49, 0x56, 255}
	}

	switch c.Style.Fg {
	case ColorPlayer, ColorLogItemUse:
		fg = color.RGBA{0x46, 0x95, 0xf7, 255}
	case ColorMonster:
		fg = color.RGBA{0xfa, 0x57, 0x50, 255}
	case ColorLogPlayerAttack, ColorStatusHealthy:
		fg = color.RGBA{0x75, 0xb9, 0x38, 255}
	case ColorLogMonsterAttack, ColorStatusWounded:
		fg = color.RGBA{0xed, 0x86, 0x49, 255}
	case ColorLogSpecial:
		fg = color.RGBA{0xf2, 0x75, 0xbe, 255}
	case ColorConsumable, ColorMenuActive:
		fg = color.RGBA{0xdb, 0xb3, 0x2d, 255}
	}

	return fg, bg
}

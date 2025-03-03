package dungeon

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/rl"
	"github.com/bayou-brogrammer/go-rl/io"
)

func (dg *dgen) GenerateTown() {
	level, err := io.ReadFile("/Users/jacoblecoq/Downloads/REXPaint-v1.70/images/prefab2.txt")
	if err != nil {
		fmt.Println(err)
	}

	dg.GeneratePrefab(string(level))

	// Generate foliage in the right half of the map
	size := dg.m.Grid.Size()
	rightHalf := gruid.NewRange(size.X/2+1, 0, size.X, size.Y)
	dg.FoliageInRange(true, rightHalf)

}

func (dg *dgen) GeneratePrefab(prefab string) {
	vault := rl.Vault{}
	err := vault.Parse(prefab)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(vault.Content())

	vault.Iter(func(p gruid.Point, c rune) {
		switch c {
		case '#':
			dg.m.SetCell(p, WallCell)
		case '.':
			dg.m.SetCell(p, FloorCell)
		case '+':
			dg.m.SetCell(p, DoorCell)
		case '=':
			dg.m.SetCell(p, RoadCell)
		}
	})
}

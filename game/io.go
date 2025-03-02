package game

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bayou-brogrammer/go-rl/game/constants"
	"github.com/bayou-brogrammer/go-rl/io"
)

/////////////
// Config
/////////////

func SaveConfig() error {
	return io.SaveGob("config.gob", &GameConfig)
}

func LoadConfig() (bool, error) {
	c := &config{}
	err := io.LoadGob("config.gob", &c)
	if err != nil {
		return false, err
	}

	if c.Version != GameConfig.Version {
		return false, nil
	}

	GameConfig = *c
	return true, nil
}

/////////////
// Game
/////////////

func (g *game) LoadGame() (bool, error) {
	dataDir, err := io.GetDataDir()
	if err != nil {
		return false, err
	}

	log.Println("Loading game from", dataDir)

	saveFile := filepath.Join(dataDir, "save")
	_, err = os.Stat(saveFile)
	if err != nil {
		// no save file, new game
		return false, nil
	}

	data, err := os.ReadFile(saveFile)
	if err != nil {
		return false, err
	}

	lg, err := io.DecodeSave[game](data)
	if err != nil {
		return false, err
	}

	if lg.Version != constants.Version {
		return false, nil
	}

	*g = *lg
	return true, nil
}

func (g *game) CheckSave() {
	if !constants.Testing {
		return
	}

	// for _, m := range g.Monsters {
	// 	mons := g.MonsterAt(m.P)
	// 	if !mons.Exists() && m.Exists() {
	// 		log.Printf("does not exist")
	// 		continue
	// 	}
	// 	if mons != m {
	// 		log.Printf("bad monster: %v vs %v", mons.Index, m.Index)
	// 	}
	// }
}

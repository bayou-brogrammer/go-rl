package main

import (
	"fmt"
	"log"

	"codeberg.org/anaseto/gruid"
	"github.com/bayou-brogrammer/go-rl/game/constants"
)

type config struct {
	Version string

	Tiles       bool
	ShowNumbers bool

	NormalModeKeys map[gruid.Key]action
	TargetModeKeys map[gruid.Key]action
}

var GameConfig config

func initConfig() error {
	GameConfig.Version = constants.Version
	GameConfig.Tiles = true

	load, err := LoadConfig()
	if err != nil {
		err = fmt.Errorf("error loading config: %v", err)
		saverr := SaveConfig()
		if saverr != nil {
			log.Printf("error resetting badly loaded config: %v", saverr)
		}
		return err
	}

	if load {
		constants.CustomKeys = true
	}

	return err
}

func (md *model) applyConfig() {
	if GameConfig.NormalModeKeys != nil {
		// md.keysNormal = GameConfig.NormalModeKeys
	}

	if GameConfig.TargetModeKeys != nil {
		// md.keysTarget = GameConfig.TargetModeKeys
	}
}

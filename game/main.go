// This is the main file of the tutorial. It defines the main routine.
package game

import (
	"context"
	"log"

	"codeberg.org/anaseto/gruid"
	sdl "codeberg.org/anaseto/gruid-sdl"
)

const (
	UIWidth   = 80
	UIHeight  = 24
	LogLines  = 2
	MapWidth  = UIWidth
	MapHeight = UIHeight - 1 - LogLines
)

func RunGame() {
	// Create a new grid with standard 80x24 size.
	gd := gruid.NewGrid(UIWidth, UIHeight)

	// Create the main application's model, using grid gd.
	m := &model{grid: gd, game: &game{}}

	// Get a TileManager for drawing fonts on the screen.
	t, err := GetTileDrawer()
	if err != nil {
		log.Fatal(err)
	}

	// Use the SDL2 driver from gruid-sdl, using the previously defined
	// TileManager.
	dr := sdl.NewDriver(sdl.Config{
		TileManager: t,
	})

	// Define a new application using the SDL2 gruid driver and our model.
	app := gruid.NewApp(gruid.AppConfig{
		Driver: dr,
		Model:  m,
	})

	// Start the application.
	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

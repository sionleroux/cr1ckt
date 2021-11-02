package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	gameWidth, gameHeight := 640, 480

	ebiten.SetWindowSize(gameWidth, gameHeight)
	ebiten.SetWindowTitle("cr1ck_t")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	game := &Game{
		Width:  gameWidth,
		Height: gameHeight,
	}

	// go NewGame(game)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Game represents the main game state
type Game struct {
	Width  int
	Height int
}

// Update calculates game logic
func (g *Game) Update() error {
	return nil
}

// Draw handles rendering the sprites
func (g *Game) Draw(screen *ebiten.Image) {
}

// Layout is hardcoded for now, may be made dynamic in future
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Width, g.Height
}

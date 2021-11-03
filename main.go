// Copyright 2021 Siôn le Roux.  All rights reserved.

package main

import (
	"embed"
	"errors"
	"image"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/*.png
var assets embed.FS

func main() {
	gameWidth, gameHeight := 640, 480

	ebiten.SetWindowSize(gameWidth, gameHeight)
	ebiten.SetWindowTitle("cr1ck_t")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	cricket := &Cricket{Object: NewObjectFromImage(loadImage("assets/cricket.png"))}

	game := &Game{
		Width:   gameWidth,
		Height:  gameHeight,
		Cricket: cricket,
		Wait:    10,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Game represents the main game state
type Game struct {
	Width   int
	Height  int
	Cricket *Cricket
	Wait    int
}

// Update calculates game logic
func (g *Game) Update() error {

	// Pressing Esc any time quits immediately
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("game quit by player")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		} else {
			ebiten.SetFullscreen(true)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Cricket.Velocity = g.Cricket.Velocity + 10
	}

	// Move the cricket
	if g.Wait++; g.Cricket.Velocity > -5 && g.Wait%10 == 0 {
		g.Cricket.Velocity--
	}
	if g.Cricket.Position.Y < g.Height-g.Cricket.Image.Bounds().Dy() || g.Cricket.Velocity > 0 {
		g.Cricket.Position.Y = g.Cricket.Position.Y - g.Cricket.Velocity
		g.Cricket.Op.GeoM.Reset()
		g.Cricket.Op.GeoM.Translate(0, float64(g.Cricket.Position.Y))
	}

	return nil
}

// Draw handles rendering the sprites
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.Cricket.Image, g.Cricket.Op)
}

// Layout is hardcoded for now, may be made dynamic in future
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Width, g.Height
}

// An Object is something that can be seen and positioned in the game
type Object struct {
	Image  *ebiten.Image
	Op     *ebiten.DrawImageOptions
	Center image.Point
}

// NewObjectFromImage makes a new game Object with fields calculated from an
// already loaded image
func NewObjectFromImage(img *ebiten.Image) *Object {
	return &Object{
		Image:  img,
		Op:     &ebiten.DrawImageOptions{},
		Center: image.Pt(0, 0),
	}
}

type Cricket struct {
	*Object
	Position image.Point
	Velocity int
}

// Load an image from embedded FS into an ebiten Image object
func loadImage(name string) *ebiten.Image {
	log.Printf("loading %s\n", name)

	file, err := assets.Open(name)
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", name, err)
	}
	defer file.Close()

	raw, err := png.Decode(file)
	if err != nil {
		log.Fatalf("error decoding file %s as PNG: %v\n", name, err)
	}

	return ebiten.NewImageFromImage(raw)
}

// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"embed"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/*.png
var assets embed.FS

func main() {
	gameWidth, gameHeight := 640, 480

	ebiten.SetWindowSize(gameWidth, gameHeight)
	ebiten.SetWindowTitle("cr1ck_t")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	cricket := &Cricket{
		Object:  NewObjectFromImage(loadImage("assets/cricket.png")),
		Jumping: true,
	}

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

	if !g.Cricket.Jumping && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Cricket.Velocity.Y += 7
		g.Cricket.Velocity.X += 5
		g.Cricket.Jumping = true
	}

	g.Wait++

	// Move the cricket
	if g.Wait%10 == 0 {
		if g.Cricket.Velocity.Y > -5 {
			g.Cricket.Velocity.Y--
		}
		if g.Cricket.Velocity.X > 0 {
			g.Cricket.Velocity.X--
		}
	}

	if g.Cricket.Position.Y < g.Height-g.Cricket.Image.Bounds().Dy() || g.Cricket.Velocity.Y > 0 {
		g.Cricket.Position.Y = g.Cricket.Position.Y - g.Cricket.Velocity.Y
		g.Cricket.Position.X = g.Cricket.Position.X + g.Cricket.Velocity.X
		g.Cricket.Op.GeoM.Reset()
		g.Cricket.Op.GeoM.Translate(float64(g.Cricket.Position.X), float64(g.Cricket.Position.Y))
	} else {
		g.Cricket.Jumping = false
	}

	return nil
}

// Draw handles rendering the sprites
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.Cricket.Image, g.Cricket.Op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("p%v - v%v\n",
		g.Cricket.Position,
		g.Cricket.Velocity,
	))
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

// Cricket is a small, jumping insect, the main character of the game
type Cricket struct {
	*Object
	Position image.Point
	Velocity image.Point
	Jumping  bool
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

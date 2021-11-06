// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"embed"
	"errors"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/ldtkgo"
	renderer "github.com/solarlune/ldtkgo/ebitenrenderer"
)

//go:embed assets/*
var assets embed.FS

func main() {
	gameWidth, gameHeight := 640, 480

	ebiten.SetWindowSize(gameWidth, gameHeight)
	ebiten.SetWindowTitle("cr1ck_t")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	cricket := &Cricket{
		Object:    NewObjectFromImage(loadImage("assets/cricket.png")),
		Jumping:   true,
		Direction: -1,
	}

	ldtkProject := loadMaps("assets/maps.ldtk")
	ebitenRenderer := renderer.NewEbitenRenderer(&EmbedLoader{"assets"})

	game := &Game{
		Width:        gameWidth,
		Height:       gameHeight,
		Cricket:      cricket,
		Wait:         10,
		TileRenderer: ebitenRenderer,
		LDTKProject:  ldtkProject,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Game represents the main game state
type Game struct {
	Width        int
	Height       int
	Cricket      *Cricket
	Wait         int
	TileRenderer *renderer.EbitenRenderer
	LDTKProject  *ldtkgo.Project
}

// Update calculates game logic
func (g *Game) Update() error {

	// Pressing Esc any time quits immediately
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errors.New("game quit by player")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		} else {
			ebiten.SetFullscreen(true)
		}
	}

	// Render map
	g.TileRenderer.Render(g.LDTKProject.Levels[0])

	// Jump
	if !g.Cricket.Jumping && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Cricket.Velocity.Y += 7
		g.Cricket.Velocity.X = (5 * g.Cricket.Direction)
		g.Cricket.Jumping = true
		g.Cricket.Direction = -g.Cricket.Direction
	}

	g.Wait++

	// Move the cricket
	if g.Wait%10 == 0 {
		if g.Cricket.Velocity.Y > -5 {
			g.Cricket.Velocity.Y--
		}
		if g.Cricket.Velocity.X < 0 {
			g.Cricket.Velocity.X++
		}
		if g.Cricket.Velocity.X > 0 {
			g.Cricket.Velocity.X--
		}
	}

	layer := g.LDTKProject.Levels[0].Layers[0]
	tile := layer.TileAt(layer.ToGridPosition(g.Cricket.Position.X, g.Cricket.Position.Y))
	if tile == nil || g.Cricket.Velocity.Y > 0 {
		g.Cricket.Position.Y = g.Cricket.Position.Y - g.Cricket.Velocity.Y
		g.Cricket.Position.X = g.Cricket.Position.X - g.Cricket.Velocity.X
		g.Cricket.Op.GeoM.Reset()
		g.Cricket.Op.GeoM.Translate(float64(g.Cricket.Position.X), float64(g.Cricket.Position.Y))
	} else {
		g.Cricket.Jumping = false
	}

	return nil
}

// Draw handles rendering the sprites
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.LDTKProject.Levels[0].BGColor)
	for _, layer := range g.TileRenderer.RenderedLayers {
		screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
	}
	screen.DrawImage(g.Cricket.Image, g.Cricket.Op)
	layer := g.LDTKProject.Levels[0].Layers[0]
	ebitenutil.DebugPrint(screen, fmt.Sprintf("p%v - v%v: %v",
		g.Cricket.Position,
		g.Cricket.Velocity,
		layer.TileAt(layer.ToGridPosition(g.Cricket.Position.X, g.Cricket.Position.Y)),
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
	Position  image.Point
	Velocity  image.Point
	Jumping   bool
	Direction int
}

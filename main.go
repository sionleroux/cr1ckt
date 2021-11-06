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

// TILELAYER is the layer to check for tile collisions
const TILELAYER int = 2

func main() {
	gameWidth, gameHeight := 640, 480

	ebiten.SetWindowSize(gameWidth, gameHeight)
	ebiten.SetWindowTitle("cr1ck_t")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	ldtkProject, err := ldtkgo.Open("maps.ldtk")
	var ebitenRenderer *renderer.EbitenRenderer
	if err == nil {
		log.Println("Found local map override, using that instead!")
		log.Println("Looking for local tileset...")
		ebitenRenderer = renderer.NewEbitenRenderer(renderer.NewDiskLoader("assets"))
	} else {
		log.Println("Using embedded map data...")
		ldtkProject = loadMaps("assets/maps.ldtk")
		ebitenRenderer = renderer.NewEbitenRenderer(&EmbedLoader{"assets"})
	}

	cricketPos := ldtkProject.Levels[0].Layers[0].EntityByIdentifier("Cricket").Position
	log.Println("Cricket starting position", cricketPos)
	cricket := &Cricket{
		Object:    NewObjectFromImage(loadImage("assets/cricket.png")),
		Jumping:   true,
		Direction: -1,
		Position:  image.Pt(cricketPos[0], cricketPos[1]),
	}

	game := &Game{
		Width:        gameWidth,
		Height:       gameHeight,
		Cricket:      cricket,
		Wait:         0,
		WaitTime:     10,
		TileRenderer: ebitenRenderer,
		LDTKProject:  ldtkProject,
		Level:        0,
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
	WaitTime     int
	TileRenderer *renderer.EbitenRenderer
	LDTKProject  *ldtkgo.Project
	Level        int
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
	g.TileRenderer.Render(g.LDTKProject.Levels[g.Level])

	// Jump
	if !g.Cricket.Jumping {
		if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
			g.Cricket.Jumping = true
			g.Cricket.Direction = -g.Cricket.Direction
			g.Cricket.Velocity.Y = g.Cricket.PrimeDuration
			g.Cricket.Velocity.X = 2 * g.Cricket.PrimeDuration * g.Cricket.Direction
		}
		g.Cricket.PrimeDuration = inpututil.KeyPressDuration(ebiten.KeySpace) / g.WaitTime
	}

	g.Wait = (g.Wait + 1) % g.WaitTime

	// Move the cricket
	if g.Wait%g.WaitTime == 0 {
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

	layer := g.LDTKProject.Levels[g.Level].Layers[TILELAYER]
	tile := layer.TileAt(layer.ToGridPosition(g.Cricket.Position.X, g.Cricket.Position.Y+g.Cricket.Image.Bounds().Dy()))
	if tile == nil || g.Cricket.Velocity.Y > 0 {
		g.Cricket.Position.X = g.Cricket.Position.X - g.Cricket.Velocity.X
		// keep within the map
		if g.Cricket.Position.X < 0 {
			g.Cricket.Position.X = 0
		}
		if g.Cricket.Position.X+g.Cricket.Image.Bounds().Dx() > g.Width {
			g.Cricket.Position.X = g.Width - g.Cricket.Image.Bounds().Dx()
		}
		g.Cricket.Position.Y = g.Cricket.Position.Y - g.Cricket.Velocity.Y
		g.Cricket.Op.GeoM.Reset()
		g.Cricket.Op.GeoM.Translate(float64(g.Cricket.Position.X), float64(g.Cricket.Position.Y))
	} else {
		g.Cricket.Jumping = false
	}

	return nil
}

// Draw handles rendering the sprites
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.LDTKProject.Levels[g.Level].BGColor)
	for _, layer := range g.TileRenderer.RenderedLayers {
		screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
	}
	screen.DrawImage(g.Cricket.Image, g.Cricket.Op)
	layer := g.LDTKProject.Levels[g.Level].Layers[TILELAYER]
	ebitenutil.DebugPrint(screen, fmt.Sprintf("p%v - v%v: %v\n%v/%v",
		g.Cricket.Position,
		g.Cricket.Velocity,
		layer.TileAt(layer.ToGridPosition(g.Cricket.Position.X, g.Cricket.Position.Y)),
		inpututil.KeyPressDuration(ebiten.KeySpace),
		g.Cricket.PrimeDuration,
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
	Position      image.Point
	Velocity      image.Point
	Jumping       bool
	PrimeDuration int
	Direction     int
}

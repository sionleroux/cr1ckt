// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"embed"
	"errors"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/ldtkgo"
	renderer "github.com/solarlune/ldtkgo/ebitenrenderer"
	"gopkg.in/ini.v1"
)

//go:embed assets/*
var assets embed.FS

// LayerEntities is the layer to use for entity positions
const LayerEntities int = 0

// LayerAuto is the layer to check for auto-tile collisions
const LayerAuto int = 1

// LayerTile is the layer to check for tile collisions
const LayerTile int = 2

// VelocityDenominator is by how much to divide the time the jump was primed to
// get the jump velocity
var VelocityDenominator int = 10

// VelocityXMultiplier is by how much to multiply the Y velocity to get the
// velocity for the X axis, it's usually bigger
var VelocityXMultiplier int = 2

// MaxPrime is the maximum jump level (after division) you can prime the cricket
// to jump for, it avoids you jumping off the screen
var MaxPrime int = 5

// DebugMode sets whether to display additional debugging info on the screen
// during playing the game or not
var DebugMode bool = false

func main() {
	gameWidth, gameHeight := 640, 480

	ebiten.SetWindowSize(gameWidth, gameHeight)
	ebiten.SetWindowTitle("cr1ck_t")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	game := &Game{
		Width:    gameWidth,
		Height:   gameHeight,
		Wait:     0,
		WaitTime: 10,
		Level:    0,
		Loading:  true,
	}

	go NewGame(game)

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
	Loading      bool
}

// NewGame populates a default game object with game data
func NewGame(game *Game) {
	log.Println("Loading game...")

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

	cricketPos := ldtkProject.Levels[0].Layers[LayerEntities].EntityByIdentifier("Cricket").Position
	log.Println("Cricket starting position", cricketPos)
	cricket := &Cricket{
		Object:    NewObjectFromImage(loadImage("assets/cricket.png")),
		Hitbox:    image.Rect(7, 24, 30, 36).Inset(1),
		Jumping:   true,
		Direction: 1,
		Position:  image.Pt(cricketPos[0], cricketPos[1]),
		Frame:     1,
		Width:     37,
	}

	game.Cricket = cricket
	game.TileRenderer = ebitenRenderer
	game.LDTKProject = ldtkProject
	game.Loading = false
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

	// Skip logic while game is loading
	if g.Loading {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.Reset(g.Level + 1)
	}

	// Render map
	g.TileRenderer.Render(g.LDTKProject.Levels[g.Level])

	// Jump
	if !g.Cricket.Jumping {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Cricket.PrimeDuration++
		} else if g.Cricket.PrimeDuration > 0 {
			g.Cricket.PrimeDuration /= VelocityDenominator
			if g.Cricket.PrimeDuration > MaxPrime {
				g.Cricket.PrimeDuration = MaxPrime
			}
			g.Cricket.Jumping = true
			g.Cricket.State = Jumping
			g.Cricket.Velocity.Y = g.Cricket.PrimeDuration
			g.Cricket.Velocity.X =
				VelocityXMultiplier * g.Cricket.PrimeDuration * g.Cricket.Direction
			g.Cricket.PrimeDuration = 0
		}
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

	// Animation ...these magic numbers refer to frames in cricket.png
	switch g.Cricket.State {
	case Idle:
		if g.Wait%g.WaitTime == 0 {
			g.Cricket.Frame = (g.Cricket.Frame + 1) % 5
		}
	case Jumping:
		if g.Cricket.Frame < 5 || g.Cricket.Frame > 8 {
			g.Cricket.Frame = 4
		}
		if g.Cricket.Frame < 8 {
			g.Cricket.Frame++
		}
	case Landing:
		if g.Cricket.Frame < 9 {
			g.Cricket.Frame = 8
		}
		if g.Cricket.Frame <= 11 {
			g.Cricket.Frame++
		}
	}

	// Landing state
	if g.Cricket.Jumping && g.Cricket.Velocity.Y <= 0 {
		g.Cricket.State = Landing
	}

	// Collision
	level := g.LDTKProject.Levels[g.Level]
	tiles := level.Layers[LayerTile]
	auto := level.Layers[LayerAuto].AllTiles()
	hitbox := g.Cricket.Hitbox.Add(image.Pt(
		g.Cricket.Position.X,
		g.Cricket.Position.Y,
	))
	collides := func(ts []*ldtkgo.Tile) bool {
		for _, v := range ts {
			if v != nil && image.Rect(
				v.Position[0], v.Position[1],
				v.Position[0]+tiles.GridSize, v.Position[1]+tiles.GridSize,
			).Overlaps(hitbox) {
				if v.ID == 11 || v.ID == 12 {
					log.Println("Hit water, restarting level")
					g.Reset(g.Level)
				}
				return true
			}
		}
		return false
	}
	colliding := collides(tiles.AllTiles()) || collides(auto)

	// Jump arc
	if !colliding || g.Cricket.Velocity.Y > 0 {
		g.Cricket.Position.X = g.Cricket.Position.X - g.Cricket.Velocity.X
		// keep within the map
		if g.Cricket.Position.X < 0 {
			g.Cricket.Position.X = 0
		}
		if g.Cricket.Position.X+g.Cricket.Width > g.Width {
			g.Cricket.Position.X = g.Width - g.Cricket.Width
		}
		g.Cricket.Position.Y = g.Cricket.Position.Y - g.Cricket.Velocity.Y
	} else if g.Cricket.Jumping {
		g.Cricket.Jumping = false
		g.Cricket.State = Idle
		g.Cricket.Direction = -g.Cricket.Direction
	}

	// Update GeoM
	g.Cricket.Op.GeoM.Reset()
	// Flip cricket direction
	g.Cricket.Op.GeoM.Scale(float64(-g.Cricket.Direction), 1)
	if g.Cricket.Direction > 0 {
		g.Cricket.Op.GeoM.Translate(float64(g.Cricket.Width), 0)
	}
	// Position cricket
	g.Cricket.Op.GeoM.Translate(float64(g.Cricket.Position.X), float64(g.Cricket.Position.Y))

	return nil
}

// Draw handles rendering the sprites
func (g *Game) Draw(screen *ebiten.Image) {
	if g.Loading {
		ebitenutil.DebugPrint(screen, "Loading...")
		return
	}

	screen.Fill(g.LDTKProject.Levels[g.Level].BGColor)
	for _, layer := range g.TileRenderer.RenderedLayers {
		screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
	}
	frameSize := g.Cricket.Width
	screen.DrawImage(g.Cricket.Image.SubImage(image.Rect(
		g.Cricket.Frame*frameSize, 0, (1+g.Cricket.Frame)*frameSize, frameSize,
	)).(*ebiten.Image), g.Cricket.Op)

	if DebugMode {
		debug(screen, g)
	}
}

// Layout is hardcoded for now, may be made dynamic in future
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Width, g.Height
}

// Reset resets the game level and cricket states to defaults for a provided
// game level
func (g *Game) Reset(level int) {
	g.Level = (level) % len(g.LDTKProject.Levels)
	log.Println("Switching to Level", g.Level)
	cricketPos := g.LDTKProject.Levels[g.Level].Layers[LayerEntities].EntityByIdentifier("Cricket").Position
	log.Println("Cricket starting position", cricketPos)
	cricket := &Cricket{
		Object:    g.Cricket.Object,
		Hitbox:    g.Cricket.Hitbox,
		Jumping:   true,
		Direction: 1,
		Position:  image.Pt(cricketPos[0], cricketPos[1]),
		Frame:     1,
		Width:     37,
	}
	g.Cricket = cricket
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

// CricketState are the different animation states a Cricket can be in
type CricketState int

const (
	// Idle is the animation state when the Cricket is not moving
	Idle CricketState = iota
	// Jumping is the animation state on the way up
	Jumping
	// Landing is the animation state on the way down
	Landing
)

// Cricket is a small, jumping insect, the main character of the game
type Cricket struct {
	*Object
	Hitbox        image.Rectangle
	Position      image.Point
	Velocity      image.Point
	Jumping       bool
	PrimeDuration int
	Direction     int
	Frame         int
	Width         int
	State         CricketState
}

func applyConfigs() {
	cfg, err := ini.Load("cr1ck_t.ini")
	log.Println(err)
	if err == nil {
		VelocityDenominator, _ = cfg.Section("").Key("VelocityDenominator").Int()
		VelocityXMultiplier, _ = cfg.Section("").Key("VelocityXMultiplier").Int()
		MaxPrime, _ = cfg.Section("").Key("MaxPrime").Int()
		DebugMode, _ = cfg.Section("").Key("DebugMode").Bool()
	}
}

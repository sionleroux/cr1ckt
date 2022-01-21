// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package cr1ckt

import (
	"embed"
	"errors"
	"image"
	"image/color"
	"log"
	"math/rand"

	"golang.org/x/image/font"
	"gopkg.in/ini.v1"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/solarlune/ldtkgo"
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

// MinPrime is the minimum jump level (after division) you can prime the cricket
// to jump for, it ensures you jump somewhat even just for a short tap
var MinPrime int = 1

// MaxPrime is the maximum jump level (after division) you can prime the cricket
// to jump for, it avoids you jumping off the screen
var MaxPrime int = 5

// DebugMode sets whether to display additional debugging info on the screen
// during playing the game or not
var DebugMode bool = false

// Double the amount of blackness that occurs after this many jumps
var BlacknessFactor int = 10

// JumpPress are the different jump states for controls
const (
	JumpPressNone int = iota
	JumpPressLeft
	JumpPressRight
	JumpPressCancel
)

// Game represents the main game state
type Game struct {
	Width        int
	Height       int
	Cricket      *Cricket
	Wait         int
	WaitTime     int
	TileRenderer *TileRenderer
	LDTKProject  *ldtkgo.Project
	Level        int
	Loading      bool
	touchIDs     []ebiten.TouchID
	blackness    Blackness
	blackFactor  int
	bg, fruit    *ebiten.Image
	cam          *camera.Camera
	win          bool
	fontBig      font.Face
	fontSmall    font.Face
}

// NewGame populates a default game object with game data
func NewGame(game *Game) {
	log.Println("Loading game...")
	// 	ldtkProject, err := ldtkgo.Open("maps.ldtk")
	var renderer *TileRenderer
	// 	if err == nil {
	// 		log.Println("Found local map override, using that instead!")
	// 		log.Println("Looking for local tileset...")
	// 		ebitenRenderer = renderer.NewEbitenRenderer(renderer.NewDiskLoader("assets"))
	// 	} else {
	log.Println("Using embedded map data...")
	ldtkProject := loadMaps("assets/maps.ldtk")
	renderer = NewTileRenderer(&EmbedLoader{"assets"})
	// }

	game.TileRenderer = renderer
	game.LDTKProject = ldtkProject
	game.fruit = loadImage("assets/fruit.png")
	game.cam = camera.NewCamera(game.Width, game.Height, 0, 0, 0, 1)
	game.Cricket = NewCricket(game.EntityByIdentifier("Cricket").Position)
	game.blackness = make(map[image.Point]bool)
	game.fontBig = loadFont(32)
	game.fontSmall = loadFont(16)

	background := loadImage("assets/background.png")
	bg := ebiten.NewImage(
		game.LDTKProject.Levels[game.Level].Width,
		game.LDTKProject.Levels[game.Level].Height,
	)
	bg.Fill(game.LDTKProject.Levels[game.Level].BGColor)
	bg.DrawImage(background, &ebiten.DrawImageOptions{})

	// Render map
	game.TileRenderer.Render(game.LDTKProject.Levels[game.Level])
	for _, layer := range game.TileRenderer.RenderedLayers {
		bg.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
	}
	for _, v := range game.LDTKProject.Levels[game.Level].Layers[LayerEntities].Entities {
		if v.Identifier == "Exit" {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(v.Position[0]), float64(v.Position[1]))
			bg.DrawImage(game.fruit, op)
		}
		// 23
	}
	game.bg = bg

	// Music
	sampleRate := 44100
	musicName := "assets/music.ogg"
	audioConext := audio.NewContext(sampleRate)
	musicFile := loadSoundFile(musicName, audioConext)
	music, err := vorbis.DecodeWithSampleRate(sampleRate, musicFile)
	if err != nil {
		log.Fatalf("error decoding file %s as Vorbis: %v\n", musicName, err)
	}
	const introLength int64 = 10113930 // pre-calculated from music editor
	musicLoop := audio.NewInfiniteLoopWithIntro(music, introLength, music.Length())
	musicPlayer, err := audio.NewPlayer(audioConext, musicLoop)
	if err != nil {
		log.Fatalf("error making music player: %v\n", err)
	}
	musicPlayer.SetVolume(0.5)
	musicPlayer.Play()

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

	// No more input when you've won
	if g.win {
		return nil
	}

	// Skip to next level
	if DebugMode && inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.Reset(g.Level + 1)
		g.win = true
	}

	// Reset jump counter
	if DebugMode && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		debugNumberOfJumps = 0
	}

	// Controls
	var JumpPress = JumpPressNone
	func() {
		// Keyboard input
		if ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyRight) {
			JumpPress = JumpPressCancel
			return
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) && ebiten.IsKeyPressed(ebiten.KeyD) {
			JumpPress = JumpPressCancel
			return
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			JumpPress = JumpPressLeft
			return
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			JumpPress = JumpPressLeft
			return
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			JumpPress = JumpPressRight
			return
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			JumpPress = JumpPressRight
			return
		}

		// Touch input
		g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])
		if len(g.touchIDs) < 1 {
			return
		}
		if len(g.touchIDs) > 2 {
			JumpPress = JumpPressCancel
			return
		}
		touchX, _ := ebiten.TouchPosition(g.touchIDs[0])
		if touchX == 0 {
			return
		}
		if touchX < g.Width/2 {
			JumpPress = JumpPressLeft
			return
		}
		if touchX >= g.Width/2 {
			JumpPress = JumpPressRight
			return
		}
	}()

	// Jump
	func() {
		if !g.Cricket.Jumping {
			// Why would you press both at once?
			if JumpPress == JumpPressCancel {
				g.Cricket.PrimeDuration = 0
				return
			}
			if JumpPress == JumpPressLeft {
				g.Cricket.Direction = 1
				g.Cricket.PrimeDuration++
				return
			}
			if JumpPress == JumpPressRight {
				g.Cricket.Direction = -1
				g.Cricket.PrimeDuration++
				return
			}
			if g.Cricket.PrimeDuration > 0 {
				g.Cricket.PrimeDuration /= VelocityDenominator
				if g.Cricket.PrimeDuration > MaxPrime {
					g.Cricket.PrimeDuration = MaxPrime
				}
				if g.Cricket.PrimeDuration < MinPrime {
					g.Cricket.PrimeDuration = MinPrime
				}
				g.Cricket.Jumping = true
				g.Cricket.State = Jumping
				debugNumberOfJumps++
				debugLastJumpStrength = g.Cricket.PrimeDuration
				g.Cricket.Velocity.Y = g.Cricket.PrimeDuration
				g.Cricket.Velocity.X =
					VelocityXMultiplier * g.Cricket.PrimeDuration * g.Cricket.Direction
				g.Cricket.PrimeDuration = 0
				g.blackFactor = debugNumberOfJumps / BlacknessFactor
				for i := 0; i < 2^g.blackFactor; i++ {
					g.blackness[image.Pt(
						rand.Intn(g.Width/16),
						rand.Intn(g.Height/16),
					)] = true
				}
			}
		}
	}()

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

	// Save pos for after collision
	oldPos := g.Cricket.Position

	// Jump arc
	if g.Cricket.Jumping {
		g.Cricket.Position.X = g.Cricket.Position.X - g.Cricket.Velocity.X
		// keep within the map
		if g.Cricket.Position.X < 0 {
			g.Cricket.Position.X = 0
		}
		if g.Cricket.Position.X+g.Cricket.Hitbox().Dx() > g.LDTKProject.Levels[g.Level].Width {
			g.Cricket.Position.X = g.Width - g.Cricket.Width
		}
		g.Cricket.Position.Y = g.Cricket.Position.Y - g.Cricket.Velocity.Y
	}

	// Collision response
	if v := Collides(g); v != nil {
		for _, w := range TilesWater {
			if v.ID == w {
				log.Println("Hit water, restarting level")
				g.Reset(g.Level)
				return nil
			}
		}
		tiles := g.LDTKProject.Levels[g.Level].Layers[LayerTile]
		exit := g.EntityByIdentifier("Exit")
		exitbox := image.Rect(
			exit.Position[0], exit.Position[1],
			exit.Position[0]+tiles.GridSize, exit.Position[1]+tiles.GridSize,
		)
		if exitbox.Overlaps(g.Cricket.Hitbox()) {
			log.Println("Found the exit, you win!")
			g.win = true
			return nil
		}
		if g.Cricket.Velocity.Y > 0 {
			g.Cricket.Velocity.Y *= -1 // Invert on hit
		} else {
			g.Cricket.Jumping = false
			g.Cricket.State = Idle
		}
		// Hop onto squishy
		if Squishy(v) {
			g.Cricket.Position = image.Pt(
				v.Position[0]+16/2-g.Cricket.Width/2,
				v.Position[1]-g.Cricket.Image.Bounds().Dy(),
			)
		}
		// Collide into impassible
		if Impassible(v) {
			g.Cricket.Position = oldPos
		}
	}
	// Landing state
	if g.Cricket.Jumping && g.Cricket.Velocity.Y <= 0 {
		g.Cricket.State = Landing
	}

	// Update GeoM
	g.Cricket.Op.GeoM.Reset()
	// Flip cricket direction
	g.Cricket.Op.GeoM.Scale(float64(-g.Cricket.Direction), 1)
	if g.Cricket.Direction > 0 {
		g.Cricket.Op.GeoM.Translate(float64(g.Cricket.Width), 0)
	}

	// Position camera
	camX, camY := 0, 0
	// Clamp the Camera to the Map dimensions
	// Surely there is an easier way to do this with maths... ಠ_ಠ
	func() {
		level := g.LDTKProject.Levels[g.Level]
		cpos := g.Cricket.Position
		cpos.X, cpos.Y = cpos.X+g.Cricket.Width/2, cpos.Y+g.Cricket.Image.Bounds().Dy()
		if cpos.X-g.Width/2 < 0 {
			camX = g.Width / 2
		} else if cpos.X+g.Width/2 > level.Width {
			camX = level.Width - g.Width/2
		} else {
			camX = cpos.X
		}
		if cpos.Y-g.Height/2 < 0 {
			camY = g.Height / 2
		} else if cpos.Y+g.Height/2 > level.Height {
			camY = level.Height - g.Height/2
		} else {
			camY = cpos.Y
		}
	}()
	g.cam.SetPosition(float64(camX), float64(camY))

	return nil
}

// Draw handles rendering the sprites
func (g *Game) Draw(screen *ebiten.Image) {
	if g.Loading {
		ebitenutil.DebugPrint(screen, "Loading...")
		return
	}

	if g.win {
		w := WinScreen(debugNumberOfJumps)
		w.Draw(g, screen)
		return
	}

	g.cam.Surface.Clear()
	g.cam.Surface.DrawImage(g.bg, g.cam.GetTranslation(0, 0))

	frameSize := g.Cricket.Width
	g.Cricket.Op.GeoM.Concat(g.cam.GetTranslation(
		float64(g.Cricket.Position.X), float64(g.Cricket.Position.Y),
	).GeoM)
	g.cam.Surface.DrawImage(g.Cricket.Image.SubImage(image.Rect(
		g.Cricket.Frame*frameSize, 0, (1+g.Cricket.Frame)*frameSize, frameSize,
	)).(*ebiten.Image), g.Cricket.Op)

	g.cam.Blit(screen)

	for b := range g.blackness {
		ebitenutil.DrawRect(screen,
			float64(b.X*16), float64(b.Y*16),
			16, 16,
			color.Black,
		)
	}

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
	g.Cricket = NewCricket(g.EntityByIdentifier("Cricket").Position)
	g.blackness = make(map[image.Point]bool)
	debugNumberOfJumps = 0
}

// EntityByIdentifier is a convenience function for the same thing in ldtkgo but
// defaulting to checking the Entities layer of the current level
func (g *Game) EntityByIdentifier(identifier string) *ldtkgo.Entity {
	return g.LDTKProject.Levels[g.Level].
		Layers[LayerEntities].
		EntityByIdentifier(identifier)
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
	hitbox        image.Rectangle
	Position      image.Point
	Velocity      image.Point
	Jumping       bool
	PrimeDuration int
	Direction     int
	Frame         int
	Width         int
	State         CricketState
}

// NewCricket returns a new Cricket object at the given position
func NewCricket(cricketPos []int) *Cricket {
	log.Println("Cricket starting position", cricketPos)
	return &Cricket{
		Object:    NewObjectFromImage(loadImage("assets/cricket.png")),
		hitbox:    image.Rect(7, 24, 30, 36).Inset(1),
		Jumping:   true,
		Direction: 1,
		Position:  image.Pt(cricketPos[0], cricketPos[1]),
		Frame:     1,
		Width:     37,
	}
}

// Hitbox returns a correctly positioned rectangular hitbox for collision
// detection with the Cricket
func (c *Cricket) Hitbox() image.Rectangle {
	return c.hitbox.Add(image.Pt(
		c.Position.X,
		c.Position.Y,
	))
}

// ApplyConfigs overrides default values with a config file if available
func ApplyConfigs() {
	log.Println("Looking for INI file...")
	cfg, err := ini.Load("cr1ckt.ini")
	log.Println(err)
	if err == nil {
		VelocityDenominator, _ = cfg.Section("").Key("VelocityDenominator").Int()
		VelocityXMultiplier, _ = cfg.Section("").Key("VelocityXMultiplier").Int()
		MaxPrime, _ = cfg.Section("").Key("MaxPrime").Int()
		MinPrime, _ = cfg.Section("").Key("MinPrime").Int()
		DebugMode, _ = cfg.Section("").Key("DebugMode").Bool()
	}
}

type Blackness map[image.Point]bool

func (b Blackness) Has(v image.Point) bool {
	_, ok := b[v]
	return ok
}

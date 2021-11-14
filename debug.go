// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func debug(screen *ebiten.Image, g *Game) {
	layer := g.LDTKProject.Levels[g.Level].Layers[LayerTile]
	hitbox := g.Cricket.Hitbox.Add(image.Pt(g.Cricket.Position.X, g.Cricket.Position.Y))

	var state string
	switch g.Cricket.State {
	case Idle:
		state = "idle"
	case Jumping:
		state = "jumping"
	case Landing:
		state = "landing"
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("p%v - v%v: %v [%v]\n%v/%v\nl:%d\n%v",
		g.Cricket.Position,
		g.Cricket.Velocity,
		hitbox,
		layer.TileAt(layer.ToGridPosition(g.Cricket.Position.X, g.Cricket.Position.Y)),
		inpututil.KeyPressDuration(ebiten.KeySpace),
		g.Cricket.PrimeDuration,
		g.Level,
		state,
	))
}

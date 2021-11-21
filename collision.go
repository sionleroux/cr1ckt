// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"image"

	"github.com/solarlune/ldtkgo"
)

// TilesImpassible is a list of tiles you can't pass through while jumping
var TilesImpassible = []int{
	0, 1, 32, 64, 65, // Earth top
	17, 21, 81, 85, // Water bank
	128,                // Earth inner
	194, 256, 260, 322, // Cave walls
	// Slops excluded intentionally
}

// TilesWater is a list of tiles that should behave like water
var TilesWater = []int{18, 19, 20, 82, 83, 84, 114, 145, 156, 146}

// Collides checks whether the Cricket is colliding with a tile
func Collides(g *Game) *ldtkgo.Tile {
	level := g.LDTKProject.Levels[g.Level]
	tiles := level.Layers[LayerTile]
	auto := level.Layers[LayerAuto]

	// This inner function is a workaround because we need to loop through both
	// Tiles and AutoTiles in exactly the same way
	overlapsTiles := func(ts []*ldtkgo.Tile) *ldtkgo.Tile {
		for _, v := range ts {
			if v != nil && image.Rect(
				v.Position[0], v.Position[1],
				v.Position[0]+tiles.GridSize, v.Position[1]+tiles.GridSize,
			).Overlaps(g.Cricket.Hitbox()) {
				return v
			}
		}
		return nil
	}

	if c := overlapsTiles(tiles.AllTiles()); c != nil {
		return c
	}
	return overlapsTiles(auto.AllTiles())
}

// XXX: impassible

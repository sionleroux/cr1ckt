// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"image"

	"github.com/solarlune/ldtkgo"
)

// ImpassibleTiles is a list of tiles you can't pass through while jumping
var ImpassibleTiles = []int{
	0,  // Earth top
	1,  // Earth top slope right
	2,  // Earth top slope left
	3,  // Stone
	10, // Water bank
	11, // Water pool
	12, // Water plant
	13, // Water stone
	32, // Earth middle
	33, // Earth middle slope right
	34, // Earth middle slope left
}

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

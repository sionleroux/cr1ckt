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

// TileSquishy is a list of tiles you should move on top of when you them
var TileSquishy = []int{
	9, 41, 73, // Rocks
	13, 45, 77, // Flowers
	15, 47, 79, // Mushrooms
}

// Collides checks whether the Cricket is colliding with a tile
func Collides(g *Game) *ldtkgo.Tile {
	level := g.LDTKProject.Levels[g.Level]
	tiles := level.Layers[LayerTile]
	auto := level.Layers[LayerAuto]
	hitbox := g.Cricket.Hitbox()
	if c := OverlapsTiles(tiles.AllTiles(), hitbox, tiles.GridSize); c != nil {
		return c
	}
	return OverlapsTiles(auto.AllTiles(), hitbox, tiles.GridSize)
}

// This inner function is a workaround because we need to loop through both
// Tiles and AutoTiles in exactly the same way
func OverlapsTiles(ts []*ldtkgo.Tile, hitbox image.Rectangle, gridSize int) *ldtkgo.Tile {
	for _, v := range ts {
		if v != nil && image.Rect(
			v.Position[0], v.Position[1],
			v.Position[0]+gridSize, v.Position[1]+gridSize,
		).Overlaps(hitbox) {
			return v
		}
	}
	return nil
}

// Impassible checks whether a tile is impassible (true) or passible (false)
func Impassible(tile *ldtkgo.Tile) bool {
	for _, t := range TilesImpassible {
		if tile.ID == t {
			return true
		}
	}
	return false
}

// Squishy checks whether a tile is squishy (true) or hard (false)
func Squishy(tile *ldtkgo.Tile) bool {
	for _, t := range TileSquishy {
		if tile.ID == t {
			return true
		}
	}
	return false
}

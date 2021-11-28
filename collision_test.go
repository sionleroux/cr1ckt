package main

import (
	"testing"

	"github.com/solarlune/ldtkgo"
)

func TestImpassible(t *testing.T) {
	IdWater, IdEarth := 114, 0
	if Impassible(&ldtkgo.Tile{ID: IdWater}) {
		t.Error("Water should be impassible")
	}
	if !Impassible(&ldtkgo.Tile{ID: IdEarth}) {
		t.Error("Earth should be impassible")
	}
}

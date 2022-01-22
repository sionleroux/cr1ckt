package cr1ckt

import (
	"image"
	"testing"

	"github.com/solarlune/ldtkgo"
)

func TestImpassible(t *testing.T) {
	IDWater, IDEarth := 114, 0
	if Impassible(&ldtkgo.Tile{ID: IDWater}) {
		t.Error("Water should be impassible")
	}
	if !Impassible(&ldtkgo.Tile{ID: IDEarth}) {
		t.Error("Earth should be impassible")
	}
}

func TestSquishy(t *testing.T) {
	IDMushroom, IDEarth := 15, 0
	if !Squishy(&ldtkgo.Tile{ID: IDMushroom}) {
		t.Error("Mushroom should be squishy")
	}
	if Squishy(&ldtkgo.Tile{ID: IDEarth}) {
		t.Error("Earth should be hard")
	}
}

// rect16 returns a 16x16 rectangle at the given coordinates
func rect16(x, y int) image.Rectangle {
	p := image.Pt(x, y)
	return image.Rectangle{p, p.Add(image.Pt(16, 16))}
}

func TestOverlapsTiles(t *testing.T) {
	const gridSize int = 16
	tiles := []*ldtkgo.Tile{
		{ID: 0, Position: []int{0, 0}},
		{ID: 114, Position: []int{32, 48}},
		{ID: 0, Position: []int{48, 48}},
	}
	cases := []struct {
		hitbox  image.Rectangle
		want    *ldtkgo.Tile
		comment string
	}{
		{rect16(0, 0), tiles[0], "inside earth"},
		{rect16(15, 15), tiles[0], "just inside earth"},
		{rect16(16, 16), nil, "next to earth"},
		{rect16(33, 33), tiles[1], "touching water"},
		{rect16(40, 43), tiles[1], "touching water before earth"},
		{rect16(50, 50), tiles[2], "touching other earth"},
		{rect16(96, 96), nil, "far away"},
	}
	for _, c := range cases {
		if res := OverlapsTiles(tiles, c.hitbox, gridSize); res != c.want {
			t.Errorf("Collision %s is %v, want %v", c.comment, res, c.want)
		}
	}
}

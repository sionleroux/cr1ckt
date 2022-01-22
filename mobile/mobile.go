// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"

	cr1ckt "github.com/sinisterstuf/cr1ckt/internal"
)

func init() {
	gameWidth, gameHeight := 640, 480

	game := &cr1ckt.Game{
		Width:    gameWidth,
		Height:   gameHeight,
		Wait:     0,
		WaitTime: 10,
		Level:    0,
		Loading:  true,
	}

	// TODO: try if this works as a goroutine
	cr1ckt.NewGame(game)

	mobile.SetGame(game)
}

// Dummy forces gomobile to compile this package.
func Dummy() {}

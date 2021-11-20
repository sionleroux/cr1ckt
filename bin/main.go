// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sinisterstuf/cr1ckt"
)

func main() {
	gameWidth, gameHeight := 640, 480

	ebiten.SetWindowSize(gameWidth, gameHeight)
	ebiten.SetWindowTitle("cr1ck_t")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	rand.Seed(time.Now().UnixNano())
	cr1ckt.ApplyConfigs()

	game := &cr1ckt.Game{
		Width:    gameWidth,
		Height:   gameHeight,
		Wait:     0,
		WaitTime: 10,
		Level:    0,
		Loading:  true,
	}

	go cr1ckt.NewGame(game)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

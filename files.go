// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"image/png"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ldtkgo"
)

// Load an project from embedded FS into an LDtk Project object
func loadMaps(name string) *ldtkgo.Project {
	log.Printf("loading %s\n", name)

	file, err := assets.Open(name)
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", name, err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("error reading from file %s: %v\n", name, err)
	}

	// Load the LDtk Project
	maps, err := ldtkgo.Read(data)
	if err != nil {
		log.Fatalf("error parsing file %s as LDtk Project: %v\n", name, err)
	}

	return maps
}

// Load an image from embedded FS into an ebiten Image object
func loadImage(name string) *ebiten.Image {
	log.Printf("loading %s\n", name)

	file, err := assets.Open(name)
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", name, err)
	}
	defer file.Close()

	raw, err := png.Decode(file)
	if err != nil {
		log.Fatalf("error decoding file %s as PNG: %v\n", name, err)
	}

	return ebiten.NewImageFromImage(raw)
}

// EmbedLoader is a TilesetLoader for the embedded FS
type EmbedLoader struct {
	BasePath string
}

// LoadTileset loads an LDtk tileset image from the embedded FS
func (l *EmbedLoader) LoadTileset(tileSetPath string) *ebiten.Image {
	return loadImage(filepath.Join(l.BasePath, tileSetPath))
}

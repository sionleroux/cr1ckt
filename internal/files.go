// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package cr1ckt

import (
	"image/png"
	"io/fs"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/solarlune/ldtkgo"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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

// LoadImage loads an ebiten.Image given the name of a file in the embedded fs
func LoadImage(name string) *ebiten.Image {
	return loadImage(name)
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

func loadSoundFile(name string, context *audio.Context) fs.File {
	log.Printf("loading %s\n", name)

	file, err := assets.Open(name)
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", name, err)
	}
	return file
}

func loadFont(size int) font.Face {
	fontdata, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	fontface, err := opentype.NewFace(fontdata, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	return fontface
}

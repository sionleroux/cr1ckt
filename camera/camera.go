// MIT License
//
// Copyright (c) 2021 Thibaud Ducasse
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package camera

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera can look at positions, zoom and rotate.
type Camera struct {
	X, Y, Rot, Scale float64
	WindowSize       image.Point
	Surface          *ebiten.Image
}

// NewCamera returns a new Camera
func NewCamera(x, y, rotation, zoom float64, windowSize image.Point) *Camera {
	return &Camera{
		X:          x,
		Y:          y,
		Rot:        rotation,
		Scale:      zoom,
		WindowSize: windowSize,
		Surface:    ebiten.NewImage(windowSize.X, windowSize.Y),
	}
}

// SetPosition looks at a position
func (c *Camera) SetPosition(x, y float64) *Camera {
	c.X = x
	c.Y = y
	return c
}

// MovePosition moves the Camera by x and y.
// Use SetPosition if you want to set the position
func (c *Camera) MovePosition(x, y float64) *Camera {
	c.X += x
	c.Y += y
	return c
}

// Rotate rotates by phi
func (c *Camera) Rotate(phi float64) *Camera {
	c.Rot += phi
	return c
}

// SetRotation sets the rotation to rot
func (c *Camera) SetRotation(rot float64) *Camera {
	c.Rot = rot
	return c
}

// Zoom *= the current zoom
func (c *Camera) Zoom(mul float64) *Camera {
	c.Scale *= mul
	if c.Scale <= 0.01 {
		c.Scale = 0.01
	}
	c.Resize(c.WindowSize.X, c.WindowSize.Y)
	return c
}

// SetZoom sets the zoom
func (c *Camera) SetZoom(zoom float64) *Camera {
	c.Scale = zoom
	if c.Scale <= 0.01 {
		c.Scale = 0.01
	}
	c.Resize(ebiten.WindowSize())
	return c
}

// Resize resizes the camera Surface
func (c *Camera) Resize(w, h int) *Camera {
	newW := int(float64(w) * 1.0 / c.Scale)
	newH := int(float64(h) * 1.0 / c.Scale)
	if newW <= 16384 && newH <= 16384 {
		c.Surface.Dispose()
		c.Surface = ebiten.NewImage(newW, newH)
	}
	return c
}

// GetTranslation returns the coordinates based on the given x,y offset and the
// camera's position
func (c *Camera) GetTranslation(x, y float64) *ebiten.DrawImageOptions {
	w, h := c.Surface.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w)/2, float64(h)/2)
	op.GeoM.Translate(-c.X+x, -c.Y+y)
	return op
}

// Blit draws the camera's surface to the screen and applies zoom
func (c *Camera) Blit(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := c.Surface.Size()
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0

	op.GeoM.Rotate(c.Rot)
	op.GeoM.Translate(-cx, -cy)
	op.GeoM.Scale(c.Scale, c.Scale)
	op.GeoM.Translate(cx*c.Scale, cy*c.Scale)

	screen.DrawImage(c.Surface, op)
}

// GetScreenCoords converts world coords into screen coords
func (c *Camera) GetScreenCoords(x, y float64) (float64, float64) {
	w, h := c.WindowSize.X, c.WindowSize.Y
	co := math.Cos(c.Rot)
	si := math.Sin(c.Rot)

	x, y = x-c.X, y-c.Y
	x, y = co*x-si*y, si*x+co*y

	return x*c.Scale + float64(w)/2, y*c.Scale + float64(h)/2
}

// GetWorldCoords converts screen coords into world coords
func (c *Camera) GetWorldCoords(x, y float64) (float64, float64) {
	w, h := c.WindowSize.X, c.WindowSize.Y
	co := math.Cos(-c.Rot)
	si := math.Sin(-c.Rot)

	x, y = (x-float64(w)/2)/c.Scale, (y-float64(h)/2)/c.Scale
	x, y = co*x-si*y, si*x+co*y

	return x + c.X, y + c.Y
}

// GetCursorCoords converts cursor/screen coords into world coords
func (c *Camera) GetCursorCoords() (float64, float64) {
	cx, cy := ebiten.CursorPosition()
	return c.GetWorldCoords(float64(cx), float64(cy))
}

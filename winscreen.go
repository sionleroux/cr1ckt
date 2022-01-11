package cr1ckt

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type WinScreen int

func (w WinScreen) Draw(g *Game, screen *ebiten.Image) {
	txt := "YOU WIN!"
	txtF, _ := font.BoundString(g.fontBig, txt)
	txtW := (txtF.Max.X - txtF.Min.X).Ceil() / 2
	txtH := (txtF.Max.Y - txtF.Min.Y).Ceil() * 2
	text.Draw(screen, txt, g.fontBig, g.Width/2-txtW, txtH, color.White)

	txt = fmt.Sprintf("%d JUMPS", w)
	txtF, _ = font.BoundString(g.fontBig, txt)
	txtW = (txtF.Max.X - txtF.Min.X).Ceil() / 2
	text.Draw(screen, txt, g.fontBig, g.Width/2-txtW, txtH*2, color.White)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.Width/2-g.Cricket.Width/2), float64(txtH*3))
	screen.DrawImage(g.Cricket.Image.SubImage(image.Rect(
		0, 0, g.Cricket.Width, g.Cricket.Width,
	)).(*ebiten.Image), op)

	txt = "Programmer: Siôn le Roux"
	txtF, _ = font.BoundString(g.fontSmall, txt)
	txtW = (txtF.Max.X - txtF.Min.X).Ceil() / 2
	txtH = (txtF.Max.Y - txtF.Min.Y).Ceil() * 2
	text.Draw(screen, txt, g.fontSmall, g.Width/2-txtW, txtH*8, color.White)

	txt = "Designer: Tristan le Roux"
	txtF, _ = font.BoundString(g.fontSmall, txt)
	txtW = (txtF.Max.X - txtF.Min.X).Ceil() / 2
	text.Draw(screen, txt, g.fontSmall, g.Width/2-txtW, txtH*9, color.White)

	txt = "Artist: Rowan Lindeque"
	txtF, _ = font.BoundString(g.fontSmall, txt)
	txtW = (txtF.Max.X - txtF.Min.X).Ceil() / 2
	text.Draw(screen, txt, g.fontSmall, g.Width/2-txtW, txtH*10, color.White)

	txt = "Android build: Levente Jónás"
	txtF, _ = font.BoundString(g.fontSmall, txt)
	txtW = (txtF.Max.X - txtF.Min.X).Ceil() / 2
	text.Draw(screen, txt, g.fontSmall, g.Width/2-txtW, txtH*11, color.White)

	txt = "Music: Sad Night - Robbero ©2021 CC-BY"
	txtF, _ = font.BoundString(g.fontSmall, txt)
	txtW = (txtF.Max.X - txtF.Min.X).Ceil() / 2
	text.Draw(screen, txt, g.fontSmall, g.Width/2-txtW, txtH*13, color.White)
}

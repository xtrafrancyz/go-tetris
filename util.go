package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var (
	ColorWhite = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	ColorBlack = color.RGBA{R: 0, G: 0, B: 0, A: 0}
)

func drawEmptyRect(dst *ebiten.Image, x, y, w, h float64, clr color.Color) {
	x += 1
	h -= 1
	w -= 1
	ebitenutil.DrawLine(dst, x, y, x+w, y, clr)
	ebitenutil.DrawLine(dst, x, y, x, y+h+1, clr)
	ebitenutil.DrawLine(dst, x, y+h, x+w, y+h, clr)
	ebitenutil.DrawLine(dst, x+w, y, x+w, y+h, clr)
}

func isKeyPressed(key ebiten.Key) bool {
	if inpututil.IsKeyJustPressed(key) {
		return true
	}
	d := inpututil.KeyPressDuration(key)
	return d > 15 && d%5 == 0
}

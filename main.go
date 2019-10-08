package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

var colors = []color.Color{
	color.RGBA{R: 0xff, G: 0x00, B: 0xff, A: 0xff},
	color.RGBA{R: 0x00, G: 0xff, B: 0xff, A: 0xff},
	color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff},
	color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff},
	color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
	color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff},
}

const (
	fieldWidth   = 10
	fieldHeight  = 20
	pointSize    = 13
	pointOffset  = 1
	fieldOffsetX = 4
	fieldOffsetY = 3
	afterFieldX  = fieldOffsetX + fieldWidth*(pointSize+pointOffset) + 3
	afterFieldY  = fieldOffsetY + fieldHeight*(pointSize+pointOffset) + 3
)

type mode int

const (
	ModeGame = iota
	ModeGameOver
)

type Game struct {
	mode mode

	ticks  int
	points int
	lines  int

	//          -1 = пустая клетка
	//         0-7 = заполненная статикой
	// 1000 - 1007 = текущий фрагмент
	field [][]int

	nextFigure     *figure
	nextColor      int
	activeFigure   *figure
	activeRotation int
	activeColor    int
	fx             int
	fy             int
}

func (g *Game) init() {
	g.field = make([][]int, fieldWidth)
	for i := 0; i < len(g.field); i++ {
		g.field[i] = make([]int, fieldHeight)
		for k := range g.field[i] {
			g.field[i][k] = -1
		}
	}
	g.spawnFigure() // для генерации следующей фигуры
	g.spawnFigure()
	g.points = 0
	g.ticks = 1
	g.lines = 0

	g.mode = ModeGame
}

func (g *Game) update(screen *ebiten.Image) error {
	g.input()
	switch g.mode {
	case ModeGame:
		g.ticks++
		if g.hasIntersections(0, 1, 0) {
			g.cutLines()
			g.spawnFigure()
			if g.hasIntersections(0, 0, 0) {
				g.mode = ModeGameOver
			} else {
				g.move(0, 0, 0) // рисуем, если можно
			}
		} else {
			if g.ticks%10 == 0 {
				g.move(0, 1, 0)
			}
		}
	case ModeGameOver:

	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.Fill(ColorBlack)
	for x, column := range g.field {
		for y, value := range column {
			if value == -1 {
				continue
			}
			if value >= 1000 {
				value -= 1000
			}
			ebitenutil.DrawRect(screen,
				float64(fieldOffsetX+x*(pointSize+pointOffset)),
				float64(pointOffset+fieldOffsetY+y*(pointSize+pointOffset)),
				pointSize, pointSize,
				colors[value],
			)
		}
	}
	drawEmptyRect(screen,
		fieldOffsetX-2,
		fieldOffsetY-1,
		fieldWidth*(pointSize+pointOffset)+3,
		fieldHeight*(pointSize+pointOffset)+3,
		ColorWhite,
	)

	sidebarY := 6
	for _, line := range []string{
		fmt.Sprintf("Points: %d", g.points),
		fmt.Sprintf("Lines: %d", g.lines),
		"",
		"Controls:",
		" Arrows - move",
		" Del - reset game",
		"",
	} {
		if line != "" {
			text.Draw(screen, line, arcadeFont, afterFieldX+1, fieldOffsetY+sidebarY, ColorWhite)
		}
		sidebarY += 16
	}

	text.Draw(screen, "Next:", arcadeFont, afterFieldX+1, fieldOffsetY+sidebarY, ColorWhite)
	sidebarY += 8

	nm := g.nextFigure.getMatrix(0)
	nox, noy := g.nextFigure.getOffset(0)
	for x := nox; x < 4; x++ {
		for y := noy; y < 4; y++ {
			if nm[x][y] {
				ebitenutil.DrawRect(screen,
					float64(afterFieldX+5+(x-nox)*(pointSize+pointOffset)),
					float64(pointOffset+sidebarY+(y-noy)*(pointSize+pointOffset)),
					pointSize, pointSize,
					colors[g.nextColor],
				)
			}
		}
	}

	if g.mode == ModeGameOver {
		text.Draw(screen, "   Game Over", arcadeFont, afterFieldX+1, fieldOffsetY+sidebarY+60, colors[3])
		text.Draw(screen, " Press Del to", arcadeFont, afterFieldX+1, fieldOffsetY+sidebarY+80, colors[4])
		text.Draw(screen, "  play again", arcadeFont, afterFieldX+1, fieldOffsetY+sidebarY+96, colors[4])
	}

	return nil
}

func (g *Game) input() {
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		g.init()
		return
	}
	if g.mode == ModeGame {
		if isKeyPressed(ebiten.KeyLeft) {
			g.tryMove(-1, 0, 0)
		}
		if isKeyPressed(ebiten.KeyRight) {
			g.tryMove(1, 0, 0)
		}
		if d := inpututil.KeyPressDuration(ebiten.KeyDown); d > 0 && d%6 == 1 {
			g.tryMove(0, 1, 0)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.tryMove(0, 0, 1)
		}
	}
}

func (g *Game) spawnFigure() {
	g.activeFigure = g.nextFigure
	g.activeColor = g.nextColor
	g.nextFigure = figures[rand.Intn(len(figures))]
	g.nextColor = rand.Intn(len(colors))
	g.activeRotation = 0
	if g.activeFigure != nil {
		_, oy := g.activeFigure.getOffset(0)
		g.fx = fieldWidth/2 - 2
		g.fy = -oy
	}
}

func (g *Game) cutLines() {
	// сохранение текущей фигуры
	f := g.activeFigure.getMatrix(g.activeRotation)
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if f[x][y] {
				g.field[x+g.fx][y+g.fy] = g.activeColor
			}
		}
	}

	lines := 0
	for y := 0; y < fieldHeight; y++ {
		valid := true
		for x := 0; x < fieldWidth; x++ {
			if g.field[x][y] == -1 {
				valid = false
				break
			}
		}
		if valid {
			// Сдвигаем все верхние линии вниз
			for y1 := y; y1 > 0; y1-- {
				for x := 0; x < fieldWidth; x++ {
					g.field[x][y1] = g.field[x][y1-1]
				}
			}
			// Сетим верхнюю линию как пустую
			for x := 0; x < fieldWidth; x++ {
				g.field[x][0] = -1
			}
			lines++
		} else if lines > 0 {
			g.saveLines(lines)
			lines = 0
		}
	}
	g.saveLines(lines)
}

func (g *Game) saveLines(lines int) {
	switch lines {
	case 1:
		g.points += 40
	case 2:
		g.points += 100
	case 3:
		g.points += 300
	case 4:
		g.points += 1200
	}
	g.lines += lines
}

func (g *Game) tryMove(dx, dy, rotate int) {
	if !g.hasIntersections(dx, dy, rotate) {
		g.move(dx, dy, rotate)
	}
}

func (g *Game) move(dx, dy, rotate int) {
	f := g.activeFigure.getMatrix(g.activeRotation)
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if f[x][y] {
				g.field[x+g.fx][y+g.fy] = -1
			}
		}
	}
	g.activeRotation += rotate
	g.fx += dx
	g.fy += dy
	f = g.activeFigure.getMatrix(g.activeRotation)
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if f[x][y] {
				g.field[x+g.fx][y+g.fy] = g.activeColor + 1000
			}
		}
	}
}

func (g *Game) hasIntersections(dx, dy, rotate int) bool {
	f := g.activeFigure.getMatrix(g.activeRotation + rotate)
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if f[x][y] {
				fx := x + g.fx + dx
				fy := y + g.fy + dy
				if fx < 0 || fy < 0 || fx >= fieldWidth || fy >= fieldHeight {
					return true
				}
				val := g.field[fx][fy]
				if val != -1 && val < 1000 {
					return true
				}
			}
		}
	}
	return false
}

func main() {
	rand.Seed(time.Now().Unix())

	g := &Game{}
	g.init()
	if runtime.GOARCH == "js" || runtime.GOOS == "js" {
		ebiten.SetFullscreen(true)
	}
	if err := ebiten.Run(g.update, 290, 287, 2, "Tetris"); err != nil {
		log.Fatal(err)
	}
}

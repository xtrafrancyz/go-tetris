package main

var figures = func() []*figure {
	list := make([]*figure, 0)
	create := func(matrix []string) {
		pieces := len(matrix) / 4
		f := figure{}
		f.matrix = make([][][]bool, pieces)
		for i := 0; i < pieces; i++ {
			f.matrix[i] = make([][]bool, 4)
			for y := 0; y < 4; y++ {
				f.matrix[i][y] = make([]bool, 4)
			}
			for y := 0; y < 4; y++ {
				idx := i + y*pieces
				for x := 0; x < 4; x++ {
					f.matrix[i][x][y] = matrix[idx][x] == 'x'
				}
			}
		}
		list = append(list, &f)
	}
	create([]string{ // O
		"    ",
		" xx ",
		" xx ",
		"    ",
	})
	create([]string{ // I
		"    ", "  x ",
		"xxxx", "  x ",
		"    ", "  x ",
		"    ", "  x ",
	})
	create([]string{ // S
		"    ", "  x ",
		"  xx", "  xx",
		" xx ", "   x",
		"    ", "    ",
	})
	create([]string{ // Z
		"    ", "   x",
		" xx ", "  xx",
		"  xx", "  x ",
		"    ", "    ",
	})
	create([]string{ // L
		"    ", "  x ", "   x", " xx ",
		" xxx", "  x ", " xxx", "  x ",
		" x  ", "  xx", "    ", "  x ",
		"    ", "    ", "    ", "    ",
	})
	create([]string{ // J
		"    ", "  xx", " x  ", "  x ",
		" xxx", "  x ", " xxx", "  x ",
		"   x", "  x ", "    ", " xx ",
		"    ", "    ", "    ", "    ",
	})
	create([]string{ // T
		"    ", "  x ", "  x ", "  x ",
		" xxx", "  xx", " xxx", " xx ",
		"  x ", "  x ", "    ", "  x ",
		"    ", "    ", "    ", "    ",
	})
	return list
}()

type figure struct {
	matrix [][][]bool
}

func (f *figure) getMatrix(rotate int) [][]bool {
	idx := rotate % len(f.matrix)
	if rotate < 0 {
		idx = len(f.matrix) - idx
	}
	return f.matrix[idx]
}

func (f *figure) getOffset(rotate int) (int, int) {
	ox := 0
	oy := 0
	m := f.getMatrix(rotate)

o:
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if m[x][y] {
				break o
			}
		}
		ox++
	}
o2:
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if m[x][y] {
				break o2
			}
		}
		oy++
	}
	return ox, oy
}

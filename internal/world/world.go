package world

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"strings"
)

const (
	Me          = '@'
	Exit        = 'Q'
	Wall        = 'w'
	Space       = ' '
	Trace       = '.'
	RouteNode   = '*'
	FrameBorder = '+'
)

type (
	GeoPosition [3]int // x,y,distance
	geo2D       [][]byte
	World       struct {
		geoMap         geo2D
		width, height  int // map size
		startX, startY int // initial position
		exitX, exitY   int // target position
		posX, posY     int // last position
	}
)

func (p GeoPosition) ToArray() [2]int {
	return [2]int{p[0], p[1]}
}

func Construct(text string) (*World, error) {
	m2d := loadFromText(text)
	if len(m2d) < 1 {
		return nil, errors.New("bad data or constructor failed")
	}
	return construct(m2d), nil
}

func construct(m2d geo2D) *World {

	w := World{
		width:  len(m2d[0]),
		height: len(m2d),
		geoMap: m2d,
	}

	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			v := &m2d[y][x]
			switch *v {
			case Exit:
				w.exitX, w.exitY = x, y
			case Me:
				w.startX, w.startY = x, y
				*v = Space // освобождаем место, где мы стоим
			}
		}
	}
	w.posX, w.posY = w.startX, w.startY
	return &w
}

// loadFromText
//
// Example:
//
//	m2d := loadFromText(`
//		wwwwwwwwwwwww
//		w@        w w
//		w wwww wwww Q
//		w           w
//		wwwwwwwwwwwww
//	`)
//	construct(m2d)
func loadFromText(mapAsString string) geo2D {
	var lines []string
	var maxRowLen = 0
	var sc = bufio.NewScanner(strings.NewReader(mapAsString))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			if l := len(line); l > maxRowLen {
				maxRowLen = l
			}
			lines = append(lines, line)
		}
	}

	rows := make([][]byte, len(lines))
	for n, line := range lines {
		row := make([]byte, maxRowLen)
		for i, c := range line {
			row[i] = byte(c)
		}
		rows[n] = row
	}
	return rows
}

func (w *World) GetSizes() (int, int) {
	return w.width, w.height
}

func (w *World) GetStart() GeoPosition {
	return GeoPosition{w.startX, w.startY}
}

func (w *World) GetExit() GeoPosition {
	return GeoPosition{w.exitX, w.exitY}
}

func (w *World) SetStart(start GeoPosition) {
	w.startX, w.startY = start[0], start[1]
	w.posX, w.posY = w.startX, w.startY
}

func (w *World) SetExit(exit GeoPosition) {
	w.exitX, w.exitY = exit[0], exit[1]
	w.SetPoint(w.exitX, w.exitY, Exit)
}

func (w *World) SetPoint(x, y int, value byte) {
	w.geoMap[y][x] = value
}

func (w *World) GetPoint(x, y int) byte {
	return w.geoMap[y][x]
}

func (w *World) getPointAsSymbol(x, y int) string {
	return string(w.GetPoint(x, y))
}

func (w *World) moveablePoint(x, y int) bool {
	if y < w.height && x < w.width && x >= 0 && y >= 0 {
		return w.GetPoint(x, y) != Wall
	}
	return false
}

func (w *World) Move(x, y int, traceEnabled bool) error {

	if traceEnabled {
		if w.posY == y {
			from, to := w.posX, x
			if from > to {
				from, to = to, from
			}
			for i := from; i < to; i++ {
				if w.GetPoint(i, y) == RouteNode {
					continue
				}
				w.SetPoint(i, y, Trace)
			}
		} else if w.posX == x {
			from, to := w.posY, y
			if from > to {
				from, to = to, from
			}
			for i := from; i < to; i++ {
				if w.GetPoint(x, i) == RouteNode {
					continue
				}
				w.SetPoint(x, i, Trace)
			}
		} else {
			msg := "stop diagonal moving [%d,%d] -> [%d,%d]"
			return fmt.Errorf(msg, x, y, w.posX, w.posY)
		}
	}
	w.SetPoint(w.posX, w.posY, RouteNode)
	w.SetPoint(x, y, Me)
	w.posX, w.posY = x, y
	return nil
}

// FindNextMoves возвращает возможные позиции для очередного перемещения
// из точки заданной `[fromX, fromY]`
// TODO: distance and maxFn()
func (w *World) FindNextMoves(fromX, fromY, exitX, exitY int) []GeoPosition {

	//w.posX, w.posY = fromX, fromY

	//var y1, y2, x1, x2 int
	//x1, x2 = w.getHorizontalRange(fromX, fromY)
	//fmt.Printf("Можем перемещаться по горизонтали от [%d, %d] до [%d, %d]\n", x1, fromY, x2, fromY)
	//y1, y2 = w.getVerticalRange(fromX, fromY)
	//fmt.Printf("Можем перемещаться по вертикали от [%d, %d] до [%d, %d]\n", fromX, y1, fromX, y2)
	horizontals := w.findHorizontals(fromX, fromY)
	verticals := w.findVerticals(fromX, fromY)
	//fmt.Printf("Возможные перемещения:\n")
	//fmt.Println(horizontals)
	//fmt.Println(verticals)
	var moves []GeoPosition

	if w.canMoveTo(fromX, fromY, exitX, exitY) {
		// добавим саму точку выхода в очередное возможное перемещение
		distance := func() int {
			var distance int
			if fromX == exitX {
				distance = int(math.Abs(float64(fromY - exitY)))
			} else {
				distance = int(math.Abs(float64(fromX - exitX)))
			}
			return distance
		}()
		mov := GeoPosition{exitX, exitY, distance}
		moves = append(moves, mov)
	}

	for _, h := range horizontals {
		y, xFrom, xTo := h[0], h[1], h[2]

		for _, v := range verticals {
			x, yFrom, yTo := v[0], v[1], v[2]

			if x == fromX && y == fromY { // пропускаем текущую точку
				continue
			}
			if x != fromX && y != fromY { // точки не лежит на одной прямой
				continue
			}
			if yFrom <= y && y <= yTo && xFrom <= x && x <= xTo {
				mov := GeoPosition{x, y, maxFn(h[3], v[3])}
				moves = append(moves, mov)
			}
		}
	}
	return moves
}

// findHorizontals возвращает доступные для перемещения вертикали
// из горизонтали, что задана точкой `[fromX, fromY]`
func (w *World) findVerticals(fromX, fromY int) (axes [][4]int) {
	iMin, iMax := w.getHorizontalRange(fromX, fromY)
	for i := iMin; i <= iMax; i++ {
		from, to := w.getVerticalRange(i, fromY)
		if from < to { // to - from > 0
			distance := int(math.Abs(float64(i - from)))
			axes = append(axes, [4]int{i, from, to, distance})
		}
	}
	//sort.Slice(axes, func(i, j int) bool { return axes[i][3] < axes[j][3] })
	return
}

// findHorizontals возвращает доступные для перемещения горизонтали
// из вертикали, что задана точкой `[fromX, fromY]`
func (w *World) findHorizontals(fromX, fromY int) (axes [][4]int) {
	iMin, iMax := w.getVerticalRange(fromX, fromY)
	for i := iMin; i <= iMax; i++ {
		from, to := w.getHorizontalRange(fromX, i)
		if from < to { // to - from > 0
			distance := int(math.Abs(float64(i - from)))
			axes = append(axes, [4]int{i, from, to, distance})
		}
	}
	//sort.Slice(axes, func(i, j int) bool { return axes[i][3] < axes[j][3] })
	return
}

// getVerticalRange возвращает минимальное и максимальное значения по вертикали
// на которое можно переместиться из точки `[fromX, fromY]`
func (w *World) getVerticalRange(fromX, fromY int) (min, max int) {
	min, max = fromY, fromY
	for i := fromY - 1; w.moveablePoint(fromX, i); i-- {
		min = i
	}
	for i := fromY + 1; w.moveablePoint(fromX, i); i++ {
		max = i
	}
	return
}

// getHorizontalRange возвращает минимальное и максимальное значения по горизонтали
// на которые можно переместиться из точки `[fromX, fromY]`
func (w *World) getHorizontalRange(fromX, fromY int) (min, max int) {
	min, max = fromX, fromX
	for i := fromX - 1; w.moveablePoint(i, fromY); i-- {
		min = i
	}
	for i := fromX + 1; w.moveablePoint(i, fromY); i++ {
		max = i
	}
	return
}

func (w *World) canMoveTo(fromX, fromY, toX, toY int) bool {
	if fromX == toX {
		return w.canMoveToY(fromX, fromY, toY)
	} else if fromY == toY {
		return w.canMoveToX(fromX, fromY, toX)
	}
	return false
}

func (w *World) canMoveToX(fromX, fromY, toX int) bool {
	from, to := fromX, toX
	if from > to {
		from, to = to, from
	}
	for x := from + 1; x < to; x++ {
		if !w.moveablePoint(x, fromY) {
			return false
		}
	}
	return true
}

func (w *World) canMoveToY(fromX, fromY, toY int) bool {
	from, to := fromY, toY
	if from > to {
		from, to = to, from
	}
	for y := from + 1; y < to; y++ {
		if !w.moveablePoint(fromX, y) {
			return false
		}
	}
	return true
}

func maxFn(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func PrintMe(w *World) {
	PrintMap(w, w.posX, w.posY)
}

func PrintMapOnly(w *World) {
	PrintMap(w, -1, -1)
}

func PrintMap(w *World, mePosX, mePosY int) {

	fmt.Printf("Map size: %dx%d\n", w.width, w.height)
	fmt.Printf("Start position: [%d,%d]\n", w.startX, w.startY)
	fmt.Printf("Exit position: [%d,%d]\n", w.exitX, w.exitY)
	fmt.Println()

	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			if x == mePosX && y == mePosY {
				fmt.Print(" ", string(Me))
				continue
			}
			fmt.Print(" ", w.getPointAsSymbol(x, y))
		}
		fmt.Println()
	}
	fmt.Println()
}

func (w *World) Pack() [][]byte {
	src := w.geoMap
	result := make([][]byte, len(src))
	for i := range src {
		srcRow := src[i]
		row := make([]byte, len(srcRow))
		copy(row, srcRow)
		result[i] = row
	}
	return result
}

func (w *World) Unpack(src [][]byte) {
	result := make([][]byte, len(src))
	for i := range src {
		srcRow := src[i]
		row := make([]byte, len(srcRow))
		copy(row, srcRow)
		result[i] = row
	}
	copy(w.geoMap, result)
}

func (w *World) SetRectangle(x1, y1, x2, y2 int) {

	for i := x1; i < x2; i++ {
		w.SetPoint(i, y1, FrameBorder)
	}
	for i := x1; i < x2; i++ {
		w.SetPoint(i, y2, FrameBorder)
	}
	for i := y1; i < y2; i++ {
		w.SetPoint(x1, i, FrameBorder)
	}
	for i := y1; i < y2; i++ {
		w.SetPoint(x2, i, FrameBorder)
	}
	w.SetPoint(x2, y2, FrameBorder)
}

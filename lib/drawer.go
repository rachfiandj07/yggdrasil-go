package lib

import "fmt"

type Drawer struct {
	canvas [][]rune
}

func NewDrawer(w, h int) (*Drawer, error) {
	if w < 0 || h < 0 {
		return nil, fmt.Errorf("width and height must be non-negative, received %d %d", w, h)
	}

	if w == 0 {
		w++
	}
	if h == 0 {
		h++
	}

	d := new(Drawer)
	d.canvas = make([][]rune, h)
	for i := range d.canvas {
		d.canvas[i] = make([]rune, w)
	}
	return d, nil
}

func (d *Drawer) DrawRune(r rune, x, y int) error {
	w, h := d.Dimens()
	if x >= w || y >= h || x < 0 || y < 0 {
		return fmt.Errorf("position (%d, %d) is outside the canvas of dimension (%d, %d)", x, y, w, h)
	}
	d.canvas[y][x] = r
	return nil
}

func (d *Drawer) DrawDrawer(e *Drawer, x, y int) error {
	w, h := d.Dimens()
	eW, eH := e.Dimens()
	if x+eW-1 >= w || y+eH-1 >= h || x < 0 || y < 0 {
		return fmt.Errorf("canvas e of dimension (%d, %d) drawn in position (%d, %d) overflows canvas d of dimension (%d, %d)", eW, eH, x, y, w, h)
	}
	for i, row := range e.canvas {
		for j, r := range row {
			d.canvas[i+y][j+x] = r
		}
	}
	return nil
}

func (d *Drawer) Dimens() (w, h int) {
	h, w = len(d.canvas), len(d.canvas[0])
	return
}

func (d *Drawer) String() string {
	var s string
	for _, row := range d.canvas {
		for _, b := range row {
			if b == 0 {
				s += " "
			} else {
				s += string(b)
			}
		}
		s += "\n"
	}
	return s
}

package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"yggdrasil-go/lib/drawer"
)

type NodeValue interface {
	Draw() *drawer.Drawer
}

type NodeInt64 int64

func (i NodeInt64) Draw() *drawer.Drawer {
	return NodeString(strconv.Itoa(int(i))).Draw()
}

type NodeString string

func (s NodeString) Draw() *drawer.Drawer {
	lines := strings.Split(string(s), "\n")
	var maxLineLength int
	for _, line := range lines {
		realLineLength := utf8.RuneCountInString(line)
		if realLineLength > maxLineLength {
			maxLineLength = realLineLength
		}
	}
	d, err := drawer.NewDrawer(maxLineLength, len(lines))
	if err != nil {
		log.Fatal(fmt.Errorf("error while allocating new drawer in NodeString.Draw: %v", err))
	}
	for y, line := range lines {

		x := 0
		for _, r := range line {
			err := d.DrawRune(r, x, y)
			if err != nil {
				log.Fatal(fmt.Errorf("error while drawing %d th rune of %s line %d in NodeString.Draw() method: %v", x, s, y, err))
			}
			x++
		}
	}
	return d
}

type NodeFloat64 float64

func (f NodeFloat64) Draw() *drawer.Drawer {
	return NodeString(fmt.Sprintf("%v", f)).Draw()
}

type NodeComplex128 complex128

// Draw satisfies the NodeValue interface.
func (z NodeComplex128) Draw() *drawer.Drawer {
	return NodeString(fmt.Sprintf("%v", z)).Draw()
}

package utils

import (
	"fmt"
	"log"
	"math"
	"sort"

	"yggdrasil-go/lib/drawer"
)

func pencil(t *Tree) *drawer.Drawer {
	dVal := t.val.Draw()
	dValW, dValH := dVal.Dimens()

	if len(t.Children()) == 0 {

		d, err := drawer.NewDrawer(dValW+2+1-dValW%2, dValH+2)
		if err != nil {
			log.Fatal(fmt.Errorf("error while allocating new drawer with no children: %v", err))
		}

		err = d.DrawDrawer(dVal, 1, 1)
		if err != nil {
			log.Fatal(fmt.Errorf("error while drawing val with no children: %v", err))
		}

		err = addBoxAround(d, 0, 0, dValW+1, dValH+1)
		if err != nil {
			log.Fatal(fmt.Errorf("error while adding box with no children: %v", err))
		}
		return d
	}

	if len(t.Children()) == 1 {
		var dChild *drawer.Drawer

		tChild, err := t.Child(0)
		if err != nil {
			log.Fatal(fmt.Errorf("error while getting child 0 with one child: %v", err))
		}
		dChild = pencil(tChild)
		dChildW, dChildH := dChild.Dimens()

		w := int(math.Max(float64(dValW+2), float64(dChildW)))
		w += 1 - w%2
		h := dValH + 3 + dChildH

		d, err := drawer.NewDrawer(w, h)
		if err != nil {
			log.Fatal(fmt.Errorf("error while allocating new drawer with one child: %v", err))
		}

		err = d.DrawDrawer(dVal, (w-dValW)/2, 1)
		if err != nil {
			log.Fatal(fmt.Errorf("error while drawing val with one child: %v", err))
		}

		err = addBoxAround(d, (w-dValW)/2-1, 0, (w-dValW)/2+dValW, dValH+1)
		if err != nil {
			log.Fatal(fmt.Errorf("error while adding box with one child: %v", err))
		}

		err = d.DrawRune('┬', w/2, dValH+1)
		if err != nil {
			log.Fatal(fmt.Errorf("error while drawing ┬ with one child: %v", err))
		}

		err = d.DrawRune('│', w/2, dValH+2)
		if err != nil {
			log.Fatal(fmt.Errorf("error while drawing | with one child: %v", err))
		}

		err = d.DrawDrawer(dChild, (w-dChildW)/2, dValH+3)
		if err != nil {
			log.Fatal(fmt.Errorf("error while drawing child drawer with one child: %v", err))
		}

		err = d.DrawRune('┴', w/2, dValH+3)
		if err != nil {
			log.Fatal(fmt.Errorf("error while drawing ┴ with one child: %v", err))
		}

		return d
	}

	nChildren := len(t.Children())
	dChildren := make([]*drawer.Drawer, 0, nChildren)
	childrenLeft := make([]int, 0, nChildren)
	childrenMiddle := make([]int, 0, nChildren)

	childrenW := 0
	maxChildH := 0

	for i, tChild := range t.Children() {
		dChild := pencil(tChild)
		dChildren = append(dChildren, dChild)
		dChildW, dChildH := dChild.Dimens()
		maxChildH = int(math.Max(float64(maxChildH), float64(dChildH)))

		if i == nChildren-1 {
			if (childrenW+dChildW)%2 == 1 {
				childrenLeft = append(childrenLeft, childrenW)
				childrenMiddle = append(childrenMiddle, childrenW+dChildW/2)
				childrenW += dChildW
			} else {
				childrenLeft = append(childrenLeft, childrenW+1)
				childrenMiddle = append(childrenMiddle, childrenW+1+dChildW/2)
				childrenW += dChildW + 1
			}
		} else {
			childrenLeft = append(childrenLeft, childrenW)
			childrenMiddle = append(childrenMiddle, childrenW+dChildW/2)
			childrenW += dChildW + 1
		}
	}

	sorted := sort.SliceIsSorted(childrenLeft, func(i, j int) bool { return childrenLeft[i] < childrenLeft[j] })
	if !sorted {
		log.Fatal(fmt.Errorf("childrenLeft is not sorted"))
	}
	sorted = sort.SliceIsSorted(childrenMiddle, func(i, j int) bool { return childrenMiddle[i] < childrenMiddle[j] })
	if !sorted {
		log.Fatal(fmt.Errorf("childrenMiddle is not sorted"))
	}

	var w int
	if dValW+2 > childrenW {
		w = dValW + 2
		for i := 0; i < nChildren; i++ {
			childrenLeft[i] += (w - childrenW) / 2
			childrenMiddle[i] += (w - childrenW) / 2
		}
	} else {
		w = childrenW
	}
	h := dValH + 3 + maxChildH

	d, err := drawer.NewDrawer(w, h)
	if err != nil {
		log.Fatal(fmt.Errorf("error while allocating new drawer with more children: %v", err))
	}

	err = d.DrawDrawer(dVal, (w-dValW)/2, 1)
	if err != nil {
		log.Fatal(fmt.Errorf("error while drawing val with more children: %v", err))
	}

	err = addBoxAround(d, (w-dValW)/2-1, 0, (w-dValW)/2+dValW, dValH+1)
	if err != nil {
		log.Fatal(fmt.Errorf("error while adding box with more children: %v", err))
	}

	for i := 0; i < nChildren; i++ {
		err = d.DrawDrawer(dChildren[i], childrenLeft[i], dValH+3)
		if err != nil {
			log.Fatal(fmt.Errorf("error while drawing %d child: %v", i, err))
		}
	}

	err = d.DrawRune('┬', w/2, dValH+1)
	if err != nil {
		log.Fatal(fmt.Errorf("error while drawing upper-link ┬ under the parent: %v", err))
	}

	for i, x := range childrenMiddle {
		err = d.DrawRune('┴', x, dValH+3)
		if err != nil {
			log.Fatal(fmt.Errorf("error while drawing lower-link ┴ above the %dth child: %v", i, err))
		}
	}

	err = d.DrawRune('╭', childrenMiddle[0], dValH+2)
	if err != nil {
		log.Fatal(fmt.Errorf("error while drawing left-corner ╭ above the left most child: %v", err))
	}
	err = d.DrawRune('╮', childrenMiddle[len(childrenMiddle)-1], dValH+2)
	if err != nil {
		log.Fatal(fmt.Errorf("error while drawing right-corner ╮ above the right most child: %v", err))
	}

	for x := childrenMiddle[0] + 1; x < childrenMiddle[len(childrenMiddle)-1]; x++ {
		underParent := x == w/2
		shouldBeAt := sort.SearchInts(childrenMiddle, x)
		aboveChild := shouldBeAt < len(childrenMiddle) && childrenMiddle[shouldBeAt] == x
		var connection rune
		switch {
		case underParent && aboveChild:
			connection = '┼'
		case underParent:
			connection = '┴'
		case aboveChild:
			connection = '┬'
		default:
			connection = '─'
		}
		err = d.DrawRune(connection, x, dValH+2)
		if err != nil {
			log.Fatal(fmt.Errorf("error while drawing %c at position %d to finish connection: %v", connection, x, err))
		}
	}

	return d
}

func addBoxAround(d *drawer.Drawer, startX, startY, endX, endY int) error {
	if startX < 0 || startY < 0 || endX < 0 || endY < 0 {
		return fmt.Errorf("can't draw on negative coordinates %d %d %d %d", startX, startY, endX, endY)
	}
	if startX > endX || startY > endY {
		return fmt.Errorf("start should be before end %d %d %d %d", startX, startY, endX, endY)
	}
	dW, dH := d.Dimens()
	if endX >= dW || endY >= dH {
		return fmt.Errorf("end overflows the drawer with dimes %d %d, %d %d %d %d", dW, dH, startX, startY, endX, endY)
	}

	err := d.DrawRune('╭', startX, startY)
	if err != nil {
		return fmt.Errorf("error while drawing ╭: %v", err)
	}
	err = d.DrawRune('╮', endX, startY)
	if err != nil {
		return fmt.Errorf("error while drawing ╮: %v", err)
	}
	err = d.DrawRune('╰', startX, endY)
	if err != nil {
		return fmt.Errorf("error while drawing ╰: %v", err)
	}
	err = d.DrawRune('╯', endX, endY)
	if err != nil {
		return fmt.Errorf("error while drawing ╯: %v", err)
	}

	for x := startX + 1; x < endX; x++ {
		for yMul := 0; yMul <= 1; yMul++ {
			err = d.DrawRune('─', x, yMul*(endY-startY)+startY)
			if err != nil {
				return fmt.Errorf("error while drawing ─: %v", err)
			}
		}
	}
	for y := startY + 1; y < endY; y++ {
		for xMul := 0; xMul <= 1; xMul++ {
			err = d.DrawRune('│', xMul*(endX-startX)+startX, y)
			if err != nil {
				return fmt.Errorf("error while drawing │: %v", err)
			}
		}
	}
	return nil
}

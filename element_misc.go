package pbc

/*
#include <pbc/pbc.h>
*/
import "C"
import "math/big"

func (el *Element) Pairing() *Pairing { return el.pairing }

func (el *Element) NewFieldElement() *Element {
	newElement := makeUncheckedElement(el.pairing, false, G1)
	C.element_init_same_as(newElement.cptr, el.cptr)
	if el.checked {
		newElement.checked = true
		newElement.fieldPtr = el.fieldPtr
		newElement.isInteger = el.isInteger
	}
	return newElement
}

func (el *Element) Len() int {
	return int(C.element_item_count(el.cptr))
}

func (el *Element) Item(i int) *Element {
	if el.checked && i >= el.Len() {
		panic(ErrOutOfRange)
	}
	newElement := &Element{
		pairing: el.pairing,
		cptr:    C.element_item(el.cptr, C.int(i)),
	}
	if el.checked {
		newElement.fieldPtr = newElement.cptr.field
		newElement.isInteger = (newElement.Len() == 0)
	}
	return newElement
}

func (el *Element) X() *big.Int {
	return el.Item(0).BigInt()
}

func (el *Element) Y() *big.Int {
	return el.Item(1).BigInt()
}

func (el *Element) BruteForceDL(g, h *Element) *Element {
	if el.checked {
		el.checkInteger()
		g.ensureChecked()
		g.checkCompatible(h)
	}
	C.element_dlog_brute_force(el.cptr, g.cptr, h.cptr)
	return el
}

func (el *Element) PollardRhoDL(g, h *Element) *Element {
	if el.checked {
		el.checkInteger()
		g.ensureChecked()
		g.checkCompatible(h)
	}
	C.element_dlog_pollard_rho(el.cptr, g.cptr, h.cptr)
	return el
}

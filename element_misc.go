package pbc

/*
#include <pbc/pbc.h>
*/
import "C"

import "math/big"

// Pairing returns the pairing associated with this element.
func (el *Element) Pairing() *Pairing { return el.pairing }

// NewFieldElement creates a new element in the same field as el. The new
// element will be unchecked if and only if el is unchecked.
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

// Len returns the length of this element. For points, this is the number of
// coordinates. For polynomials, it is the number of coefficients. For infinite
// points, it is zero. For all other values, it is zero.
func (el *Element) Len() int {
	return int(C.element_item_count(el.cptr))
}

// Item returns the specified sub-element. For points, this returns a
// coordinate. For polynomials, it returns a coefficient. For other elements,
// this operation is invalid. i must be greater than or equal to 0 and less
// than el.Len(). Bounds checking is only performed for checked elements.
func (el *Element) Item(i int) *Element {
	if el.checked && i >= el.Len() {
		panic(ErrOutOfRange)
	}
	newElement := &Element{
		pairing: el.pairing,
		cptr:    C.element_item(el.cptr, C.int(i)),
	}
	if newElement.cptr == nil {
		panic(ErrOutOfRange)
	}
	if el.checked {
		newElement.fieldPtr = newElement.cptr.field
		newElement.isInteger = (newElement.Len() == 0)
	}
	return newElement
}

// X returns the X coordinate of el. Equivalent to el.Item(0).BigInt().
//
// Requirements:
// el must be a point on an elliptic curve.
func (el *Element) X() *big.Int {
	return el.Item(0).BigInt()
}

// Y returns the Y coordinate of el. Equivalent to el.Item(1).BigInt().
//
// Requirements:
// el must be a point on an elliptic curve.
func (el *Element) Y() *big.Int {
	return el.Item(1).BigInt()
}

// BruteForceDL sets el such that g^el = h using brute force.
//
// Requirements:
// g and h must be from the same algebraic structure; and
// el must be an element of an integer mod ring (e.g., Zn for some n, typically
// the order of the algebraic structure that g lies in).
func (el *Element) BruteForceDL(g, h *Element) *Element {
	if el.checked {
		el.checkInteger()
		g.ensureChecked()
		g.checkCompatible(h)
	}
	C.element_dlog_brute_force(el.cptr, g.cptr, h.cptr)
	return el
}

// PollardRhoDL sets el such that g^el = h using Pollard rho method.
//
// Requirements:
// g and h must be from the same algebraic structure; and
// el must be an element of an integer mod ring (e.g., Zn for some n, typically
// the order of the algebraic structure that g lies in).
func (el *Element) PollardRhoDL(g, h *Element) *Element {
	if el.checked {
		el.checkInteger()
		g.ensureChecked()
		g.checkCompatible(h)
	}
	C.element_dlog_pollard_rho(el.cptr, g.cptr, h.cptr)
	return el
}

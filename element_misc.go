// Copyright © 2018 Nik Unger
//
// This file is part of The PBC Go Wrapper.
//
// The PBC Go Wrapper is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// The PBC Go Wrapper is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public
// License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with The PBC Go Wrapper. If not, see <http://www.gnu.org/licenses/>.
//
// The PBC Go Wrapper makes use of The PBC library. The PBC Library and its use
// are covered under the terms of the GNU Lesser General Public License
// version 3, or (at your option) any later version.

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

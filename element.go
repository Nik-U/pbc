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

struct element_s* newElementStruct() { return malloc(sizeof(struct element_s)); }
void freeElementStruct(struct element_s* x) {
	element_clear(x);
	free(x);
}
*/
import "C"

import "runtime"

// Element represents an element in one of the algebraic structures associated
// with a pairing. Arithmetic operations can be performed on elements to
// complete computations. Elements can also be paired using the associated
// pairing's bilinear map. Elements can be exported or imported in a variety of
// formats.
//
// The arithmetic methods for Elements generally follow the style of big.Int. A
// typical operation has a signature like this:
//
// 	func (el *Element) Add(x, y *Element) *Element
//
// This method stores x + y in el and returns el. Since these arithmetic
// operations return the targets, they can be used in method chaining:
//
// 	x.Add(a, b).Mul(x, c).Square(x)
//
// This assigns x = ((a+b)*c)^2. Whenever possible, the methods defined on
// Element use the same names as those in the math/big package.
//
// This technique is useful because it allows the target of operations to be
// different than the operands. However, several convenience functions have
// been provided to improve the readability of chained calls. These functions
// are of the form Then*, and implicitly specify the target as the first
// operand. The above example can be rewritten as:
//
// 	x.Add(a, b).ThenMul(c).ThenSquare()
//
// For some applications, it is more readable to avoid method chaining:
//
// 	x.Add(a, b)
// 	x.Mul(x, c)
// 	x.Square(x)
//
// The addition and multiplication functions perform addition and
// multiplication operations in rings and fields. For groups of points on an
// elliptic curve, such as the G1 and G2 groups associated with pairings, both
// addition and multiplication represent the group operation (and similarly
// both 0 and 1 represent the identity element). It is recommended that
// programs choose one convention and stick with it to avoid confusion.
//
// In contrast, the GT group is currently implemented as a subgroup of a finite
// field, so only multiplicative operations should be used for GT.
//
// Not all operations are valid for all elements. For example, pairing
// operations require an element from G1, an element from G2, and a target from
// GT. As another example, elements in a ring cannot be inverted in general.
//
// The PBC library does not attempt to detect invalid element operations. If an
// invalid operation is performed, several outcomes are possible. In the best
// case, the operation will be treated as a no-op. The target element might
// be set to a nonsensical value. In the worst case, the program may segfault.
//
// The pbc wrapper provides some protection against invalid operations. When
// elements are initialized by a Pairing, they can either be created as checked
// or unchecked. Unchecked elements do not perform any sanity checks; calls are
// passed directly to the C library, with the possible consequences mentioned
// above. Checked elements attempt to catch a variety of errors, such as when
// invalid operations are performed on elements from mismatched algebraic
// structures or pairings. If an error is detected, the operation will panic
// with ErrIllegalOp, ErrUncheckedOp, ErrIncompatible, or a similar error.
//
// The decision on whether or not to check operations is based solely on
// whether or not the target element is checked. Thus, if an unchecked element
// is passed a checked element as part of an operation, the operation will not
// be checked. Checked elements expect that all arguments to their methods are
// also checked, and will panic with ErrUncheckedOp if they are not.
//
// Note that not all possible errors can be detected by checked elements;
// ultimately, it is the responsibility of the caller to ensure that the
// requested computations make sense.
type Element struct {
	pairing *Pairing // Prevents garbage collection
	cptr    *C.struct_element_s

	checked   bool
	fieldPtr  *C.struct_field_s
	isInteger bool
}

func clearElement(element *Element) {
	C.freeElementStruct(element.cptr)
}

func makeUncheckedElement(pairing *Pairing, initialize bool, field Field) *Element {
	element := &Element{
		cptr:    C.newElementStruct(),
		pairing: pairing,
	}
	if initialize {
		switch field {
		case G1:
			C.element_init_G1(element.cptr, pairing.cptr)
		case G2:
			C.element_init_G2(element.cptr, pairing.cptr)
		case GT:
			C.element_init_GT(element.cptr, pairing.cptr)
		case Zr:
			C.element_init_Zr(element.cptr, pairing.cptr)
		default:
			panic(ErrUnknownField)
		}
	}
	runtime.SetFinalizer(element, clearElement)
	return element
}

func makeCheckedElement(pairing *Pairing, field Field, fieldPtr *C.struct_field_s) *Element {
	element := makeUncheckedElement(pairing, true, field)
	element.checked = true
	element.fieldPtr = fieldPtr
	element.isInteger = (field == Zr)
	element.Set0()
	return element
}

func checkFieldsMatch(f1, f2 *C.struct_field_s) {
	if f1 != f2 {
		panic(ErrIncompatible)
	}
}

func (el *Element) ensureChecked() {
	if !el.checked {
		panic(ErrUncheckedOp)
	}
}

func (el *Element) checkCompatible(other *Element) {
	other.ensureChecked()
	checkFieldsMatch(el.fieldPtr, other.fieldPtr)
}

func (el *Element) checkAllCompatible(elements ...*Element) {
	for _, other := range elements {
		el.checkCompatible(other)
	}
}

func (el *Element) checkInteger() {
	el.ensureChecked()
	if !el.isInteger {
		panic(ErrIllegalOp)
	}
}

func (el *Element) checkPoint() {
	el.ensureChecked()
	if el.isInteger {
		panic(ErrIllegalOp)
	}
}

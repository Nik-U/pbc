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

import "math/big"

// ThenAdd is an alias for el.Add(el, y).
//
// Requirements:
// el and y must be from the same algebraic structure.
func (el *Element) ThenAdd(y *Element) *Element { return el.Add(el, y) }

// ThenSub is an alias for el.Sub(el, y).
//
// Requirements:
// el and y must be from the same algebraic structure.
func (el *Element) ThenSub(y *Element) *Element { return el.Sub(el, y) }

// ThenMul is an alias for el.Mul(el, y).
//
// Requirements:
// el and y must be from the same algebraic structure.
func (el *Element) ThenMul(y *Element) *Element { return el.Mul(el, y) }

// ThenMulBig is an alias for el.MulBig(el, i).
func (el *Element) ThenMulBig(i *big.Int) *Element { return el.MulBig(el, i) }

// ThenMulInt32 is an alias for el.MulInt32(el, i).
func (el *Element) ThenMulInt32(i int32) *Element { return el.MulInt32(el, i) }

// ThenMulZn is an alias for el.MulZn(el, i).
//
// Requirements:
// i must be an element of an integer mod ring (e.g., Zn for some n).
func (el *Element) ThenMulZn(i *Element) *Element { return el.MulZn(el, i) }

// ThenDiv is an alias for el.Div(el, y).
//
// Requirements:
// el and y must be from the same algebraic structure.
func (el *Element) ThenDiv(y *Element) *Element { return el.Div(el, y) }

// ThenDouble is an alias for el.Double(el).
func (el *Element) ThenDouble() *Element { return el.Double(el) }

// ThenHalve is an alias for el.Halve(el).
func (el *Element) ThenHalve() *Element { return el.Halve(el) }

// ThenSquare is an alias for el.Square(el).
func (el *Element) ThenSquare() *Element { return el.Square(el) }

// ThenNeg is an alias for el.Neg(el).
func (el *Element) ThenNeg() *Element { return el.Neg(el) }

// ThenInvert is an alias for el.Invert(el).
func (el *Element) ThenInvert() *Element { return el.Invert(el) }

// ThenPowBig is an alias for el.PowBig(el, i).
func (el *Element) ThenPowBig(i *big.Int) *Element { return el.PowBig(el, i) }

// ThenPowZn is an alias for el.PowZn(el, i).
//
// Requirements:
// i must be an element of an integer mod ring (e.g., Zn for some n, typically
// the order of the algebraic structure that el lies in).
func (el *Element) ThenPowZn(i *Element) *Element { return el.PowZn(el, i) }

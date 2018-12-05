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

struct element_pp_s* newElementPPStruct() { return malloc(sizeof(struct element_pp_s)); }
void freeElementPPStruct(struct element_pp_s* x) {
	element_pp_clear(x);
	free(x);
}

struct pairing_pp_s* newPairingPPStruct() { return malloc(sizeof(struct pairing_pp_s)); }
void freePairingPPStruct(struct pairing_pp_s* x) {
	pairing_pp_clear(x);
	free(x);
}
*/
import "C"

import (
	"math/big"
	"runtime"
	"unsafe"
)

// Set0 sets el to zero and returns el. For curves, this sets the element to
// the infinite point (identity element).
func (el *Element) Set0() *Element {
	C.element_set0(el.cptr)
	return el
}

// Set1 sets el to one and returns el. For curves, this sets the element to the
// infinite point (identity element).
func (el *Element) Set1() *Element {
	C.element_set1(el.cptr)
	return el
}

// Rand sets el to a random value and returns el. For algebraic structures
// where this does not make sense, this is equivalent to Set0.
func (el *Element) Rand() *Element {
	C.element_random(el.cptr)
	return el
}

// Equals returns true if el == x.
//
// Requirements:
// el and x must be from the same algebraic structure.
func (el *Element) Equals(x *Element) bool {
	if el.checked {
		el.checkCompatible(x)
	}
	return int64(C.element_cmp(el.cptr, x.cptr)) == 0
}

// Is0 returns true if el is zero (or the identity element for curves).
func (el *Element) Is0() bool {
	return C.element_is0(el.cptr) != 0
}

// Is1 returns true if el is one (or the identity element for curves).
func (el *Element) Is1() bool {
	return C.element_is1(el.cptr) != 0
}

// IsSquare returns true if el is a perfect square (quadratic residue).
func (el *Element) IsSquare() bool {
	return C.element_is_sqr(el.cptr) != 0
}

// Sign returns 0 if el is 0. If el is not 0, the behavior depends on the
// algebraic structure, but has the property that el.Sign() == -neg.Sign()
// where neg is the negation of el.
func (el *Element) Sign() int {
	sign := int64(C.element_sign(el.cptr))
	if sign > 0 {
		return 1
	}
	if sign < 0 {
		return -1
	}
	return 0
}

// Add sets el = x + y and returns el. For curve points, + denotes the group
// operation.
//
// Requirements:
// el, x, and y must be from the same algebraic structure.
func (el *Element) Add(x, y *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
	}
	C.element_add(el.cptr, x.cptr, y.cptr)
	return el
}

// Sub sets el = x - y and returns el. More precisely, el = x + (-y). For curve
// points, + denotes the group operation.
//
// Requirements:
// el, x, and y must be from the same algebraic structure.
func (el *Element) Sub(x, y *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
	}
	C.element_sub(el.cptr, x.cptr, y.cptr)
	return el
}

// Mul sets el = x * y and returns el. For curve points, * denotes the group
// operation.
//
// Requirements:
// el, x, and y must be from the same algebraic structure.
func (el *Element) Mul(x, y *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
	}
	C.element_mul(el.cptr, x.cptr, y.cptr)
	return el
}

// MulBig sets el = i * x and returns el. More precisely, el = x + x + ... + x
// where there are i x's. For curve points, + denotes the group operation.
//
// Requirements:
// el and x must be from the same algebraic structure.
func (el *Element) MulBig(x *Element, i *big.Int) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_mul_mpz(el.cptr, x.cptr, &big2mpz(i).i[0])
	return el
}

// MulInt32 sets el = i * x and returns el. More precisely,
// el = x + x + ... + x where there are i x's. For curve points, + denotes the
// group operation.
//
// Requirements:
// el and x must be from the same algebraic structure.
func (el *Element) MulInt32(x *Element, i int32) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_mul_si(el.cptr, x.cptr, C.long(i))
	return el
}

// MulZn sets el = i * x and returns el. More precisely,
// el = x + x + ... + x where there are i x's. For curve points, + denotes the
// group operation.
//
// Requirements:
// el and x must be from the same algebraic structure; and
// i must be an element of an integer mod ring (e.g., Zn for some n).
func (el *Element) MulZn(x, i *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
		i.checkInteger()
	}
	C.element_mul_zn(el.cptr, x.cptr, i.cptr)
	return el
}

// Div sets el = x / y and returns el. More precisely, el = x * (1/y). For
// curve points, * denotes the group operation.
//
// Requirements:
// el, x, and y must be from the same algebraic structure.
func (el *Element) Div(x, y *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
	}
	C.element_div(el.cptr, x.cptr, y.cptr)
	return el
}

// Double sets el = x + x and returns el.
//
// Requirements:
// el and x must be from the same algebraic structure.
func (el *Element) Double(x *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_double(el.cptr, x.cptr)
	return el
}

// Halve sets el = x / 2 and returns el.
//
// Requirements:
// el and x must be from the same algebraic structure.
func (el *Element) Halve(x *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_halve(el.cptr, x.cptr)
	return el
}

// Square sets el = x * x and returns el.
//
// Requirements:
// el and x must be from the same algebraic structure.
func (el *Element) Square(x *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_square(el.cptr, x.cptr)
	return el
}

// Neg sets el = -x and returns el.
//
// Requirements:
// el and x must be from the same algebraic structure.
func (el *Element) Neg(x *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_neg(el.cptr, x.cptr)
	return el
}

// Invert sets el = 1/x and returns el.
//
// Requirements:
// el and x must be from the same algebraic structure.
func (el *Element) Invert(x *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_invert(el.cptr, x.cptr)
	return el
}

// PowBig sets el = x^i and returns el. More precisely, el = x * x * ... * x
// where there are i x's. For curve points, * denotes the group operation.
//
// Requirements:
// el and x must be from the same algebraic structure.
func (el *Element) PowBig(x *Element, i *big.Int) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_pow_mpz(el.cptr, x.cptr, &big2mpz(i).i[0])
	return el
}

// PowZn sets el = x^i and returns el. More precisely, el = x * x * ... * x
// where there are i x's. For curve points, * denotes the group operation.
//
// Requirements:
// el and x must be from the same algebraic structure; and
// i must be an element of an integer mod ring (e.g., Zn for some n, typically
// the order of the algebraic structure that x lies in).
func (el *Element) PowZn(x, i *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
		i.checkInteger()
	}
	C.element_pow_zn(el.cptr, x.cptr, i.cptr)
	return el
}

// Pow2Big sets el = x^i * y^j and returns el. This is generally faster than
// performing separate exponentiations.
//
// Requirements:
// el, x, and y must be from the same algebraic structure.
func (el *Element) Pow2Big(x *Element, i *big.Int, y *Element, j *big.Int) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
	}
	C.element_pow2_mpz(el.cptr, x.cptr, &big2mpz(i).i[0], y.cptr, &big2mpz(j).i[0])
	return el
}

// Pow2Zn sets el = x^i * y^j and returns el. This is generally faster than
// performing separate exponentiations.
//
// Requirements:
// el, x, and y must be from the same algebraic structure; and
// i and j must be elements of integer mod rings (e.g., Zn for some n,
// typically the order of the algebraic structures that x and y lie in).
func (el *Element) Pow2Zn(x, i, y, j *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
		i.checkInteger()
		j.checkInteger()
	}
	C.element_pow2_zn(el.cptr, x.cptr, i.cptr, y.cptr, j.cptr)
	return el
}

// Pow3Big sets el = x^i * y^j * z^k and returns el. This is generally faster
// than performing separate exponentiations.
//
// Requirements:
// el, x, y, and z must be from the same algebraic structure.
func (el *Element) Pow3Big(x *Element, i *big.Int, y *Element, j *big.Int, z *Element, k *big.Int) *Element {
	if el.checked {
		el.checkAllCompatible(x, y, z)
	}
	C.element_pow3_mpz(el.cptr, x.cptr, &big2mpz(i).i[0], y.cptr, &big2mpz(j).i[0], z.cptr, &big2mpz(k).i[0])
	return el
}

// Pow3Zn sets el = x^i * y^j * z^k and returns el. This is generally faster
// than performing separate exponentiations.
//
// Requirements:
// el, x, y, and z must be from the same algebraic structure; and
// i, j, and k must be elements of integer mod rings (e.g., Zn for some n,
// typically the order of the algebraic structures that x, y, and z lie in).
func (el *Element) Pow3Zn(x, i, y, j, z, k *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y, z)
		i.checkInteger()
		j.checkInteger()
		k.checkInteger()
	}
	C.element_pow3_zn(el.cptr, x.cptr, i.cptr, y.cptr, j.cptr, z.cptr, k.cptr)
	return el
}

// Power stores pre-processed information to quickly exponentiate an element.
// A Power can be generated for Element x by calling x.PreparePower(). When
// PowBig or PowZn is called with Element target and integer i, the result of
// x^i will be stored in target. Once a Power has been generated, the original
// element can be changed without affecting the pre-processed data.
type Power struct {
	source *Element // Prevents garbage collection
	pp     *C.struct_element_pp_s
}

// Source returns the Element for which the pre-processed data was generated.
func (power *Power) Source() *Element { return power.source }

// PowBig sets target = s^i where s was the source Element for the Power, and
// returns target. It is equivalent to target.PowerBig(power, i).
//
// Requirements:
// target and s must be from the same algebraic structure.
func (power *Power) PowBig(target *Element, i *big.Int) *Element {
	return target.PowerBig(power, i)
}

// PowZn sets target = s^i where s was the source Element for the Power, and
// returns target. It is equivalent to target.PowerZn(power, i).
//
// Requirements:
// target and s must be from the same algebraic structure.
func (power *Power) PowZn(target *Element, i *Element) *Element {
	return target.PowerZn(power, i)
}

func clearPower(power *Power) {
	C.freeElementPPStruct(power.pp)
}

// PreparePower generates pre-processing data for repeatedly exponentiating el.
// The returned Power can be used to raise el to a power several times, and is
// generally faster than repeatedly calling the standard Pow methods on el.
func (el *Element) PreparePower() *Power {
	power := &Power{
		source: el,
		pp:     C.newElementPPStruct(),
	}
	C.element_pp_init(power.pp, el.cptr)
	runtime.SetFinalizer(power, clearPower)
	return power
}

// PowerBig sets el = s^i where s was the source Element for the Power, and
// returns el. It is equivalent to power.PowBig(el, i).
//
// Requirements:
// el and s must be from the same algebraic structure.
func (el *Element) PowerBig(power *Power, i *big.Int) *Element {
	if el.checked {
		el.checkCompatible(power.source)
	}
	C.element_pp_pow(el.cptr, &big2mpz(i).i[0], power.pp)
	return el
}

// PowerZn sets el = s^i where s was the source Element for the Power, and
// returns el. It is equivalent to power.PowZn(el, i).
//
// Requirements:
// el and s must be from the same algebraic structure; and
// i must be an element of an integer mod ring (e.g., Zn for some n, typically
// the order of the algebraic structure that s lies in).
func (el *Element) PowerZn(power *Power, i *Element) *Element {
	if el.checked {
		el.checkCompatible(power.source)
		i.checkInteger()
	}
	C.element_pp_pow_zn(el.cptr, i.cptr, power.pp)
	return el
}

// Pair sets el = e(x,y) where e denotes the pairing operation, and returns el.
//
// Requirements:
// el, x, and y must belong to the same pairing;
// el must belong to the pairing's GT group;
// x must belong to the pairing's G1 group (or G2 for symmetric pairings); and
// y must belong to the pairing's G2 group (or G1 for symmetric pairings).
func (el *Element) Pair(x, y *Element) *Element {
	if el.checked {
		x.ensureChecked()
		y.ensureChecked()
		pairing := el.pairing.cptr
		checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
		checkFieldsMatch(x.fieldPtr, pairing.G1)
		checkFieldsMatch(y.fieldPtr, pairing.G2)
	}
	C.pairing_apply(el.cptr, x.cptr, y.cptr, el.pairing.cptr)
	return el
}

func (el *Element) doProdPair(in1, in2 []C.struct_element_s) *Element {
	x := (*C.element_t)(unsafe.Pointer(&in1[0]))
	y := (*C.element_t)(unsafe.Pointer(&in2[0]))
	C.element_prod_pairing(el.cptr, x, y, C.int(len(in1)))
	return el
}

// ProdPair sets el to the product of several pairings, and returns el. The
// elements are paired in groups of two.
//
// For example:
// 	el.ProdPair(a,b,c,d,e,f)
// will set el = e(a,b) * e(c,d) * e(e,f).
//
// Requirements:
// all elements must belong to the same pairing;
// el must belong to the pairing's GT group;
// there must be an even number of parameters;
// odd numbered parameters must belong to the pairing's G1 group (or G2 for
// symmetric pairings); and
// even numbered parameters must belong to the pairing's G2 group (or G1 for
// symmetric pairings).
func (el *Element) ProdPair(elements ...*Element) *Element {
	n := len(elements)
	if n%2 != 0 {
		panic(ErrBadPairList)
	}
	if el.checked {
		pairing := el.pairing.cptr
		checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
		for i := 1; i < n; i += 2 {
			elements[i-1].ensureChecked()
			elements[i].ensureChecked()
			checkFieldsMatch(elements[i-1].fieldPtr, pairing.G1)
			checkFieldsMatch(elements[i].fieldPtr, pairing.G2)
		}
	}
	half := n / 2
	in1 := make([]C.struct_element_s, half)
	in2 := make([]C.struct_element_s, half)
	for i, j := 0, 0; j < n; i, j = i+1, j+2 {
		in1[i] = *elements[j].cptr
		in2[i] = *elements[j+1].cptr
	}
	return el.doProdPair(in1, in2)
}

// ProdPairSlice sets el to the product of several pairings, and returns el.
// Elements from x will be paired with elements in y having the same index.
//
// For example:
// 	el.ProdPairSlice([]*Element{a,b,c}, []*Element{d,e,f})
// will set el = e(a,d) * e(b,e) * e(c,f).
//
// Requirements:
// all elements must belong to the same pairing;
// el must belong to the pairing's GT group;
// the slices must have the same number of elements;
// elements in x must belong to the pairing's G1 group (or G2 for symmetric
// pairings); and
// elements in y must belong to the pairing's G2 group (or G1 for symmetric
// pairings).
func (el *Element) ProdPairSlice(x, y []*Element) *Element {
	n := len(x)
	if n != len(y) {
		panic(ErrBadPairList)
	}
	if el.checked {
		pairing := el.pairing.cptr
		checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
		for i := 1; i < n; i++ {
			x[i].ensureChecked()
			checkFieldsMatch(x[i].fieldPtr, pairing.G1)
		}
		n = len(y)
		for i := 1; i < n; i++ {
			y[i].ensureChecked()
			checkFieldsMatch(y[i].fieldPtr, pairing.G2)
		}
	}
	in1 := make([]C.struct_element_s, n)
	in2 := make([]C.struct_element_s, n)
	for i := 0; i < n; i++ {
		in1[i] = *x[i].cptr
		in2[i] = *y[i].cptr
	}
	return el.doProdPair(in1, in2)
}

// Pairer stores pre-processed information to quickly pair an element. A Pairer
// can be generated for Element x by calling x.PreparePairer(). When Pair is
// called with Elements target and y, the result of e(x,y) will be stored in
// target. Once a Pairer has been generated, the original element can be
// changed without affecting the pre-processed data.
type Pairer struct {
	source *Element // Prevents garbage collection
	pp     *C.struct_pairing_pp_s
}

// Source returns the Element for which the pre-processed data was generated.
func (pairer *Pairer) Source() *Element { return pairer.source }

// Pair sets target = e(s,y) and returns target, where e denotes the pairing
// operation, and s was the source Element for the Pairer. It is equivalent to
// target.PairerPair(pairer, y).
//
// Requirements:
// target, s, and y must belong to the same pairing;
// target must belong to the pairing's GT group;
// s must belong to the pairing's G1 group (or G2 for symmetric pairings); and
// y must belong to the pairing's G2 group (or G1 for symmetric pairings).
func (pairer *Pairer) Pair(target *Element, y *Element) *Element {
	return target.PairerPair(pairer, y)
}

func clearPairer(pairer *Pairer) {
	C.freePairingPPStruct(pairer.pp)
}

// PreparePairer generates pre-processing data for repeatedly pairing el. The
// returned Pairer can be used to pair el several times, and is generally
// faster than repeatedly calling Pair on el.
//
// Requirements:
// el must belong to the pairing's G1 group (or G2 for symmetric pairings).
func (el *Element) PreparePairer() *Pairer {
	if el.checked {
		checkFieldsMatch(el.fieldPtr, el.pairing.cptr.G1)
	}
	pairer := &Pairer{
		source: el,
		pp:     C.newPairingPPStruct(),
	}
	C.pairing_pp_init(pairer.pp, el.cptr, el.pairing.cptr)
	runtime.SetFinalizer(pairer, clearPairer)
	return pairer
}

// PairerPair sets el = e(s,y) and returns el, where e denotes the pairing
// operation, and s was the source Element for the Pairer. It is equivalent to
// pairer.Pair(el, y).
//
// Requirements:
// el, s, and y must belong to the same pairing;
// el must belong to the pairing's GT group;
// s must belong to the pairing's G1 group (or G2 for symmetric pairings); and
// y must belong to the pairing's G2 group (or G1 for symmetric pairings).
func (el *Element) PairerPair(pairer *Pairer, y *Element) *Element {
	if el.checked {
		pairer.source.ensureChecked()
		y.ensureChecked()
		pairing := el.pairing.cptr
		checkFieldsMatch(pairer.source.fieldPtr, pairing.G1)
		checkFieldsMatch(y.fieldPtr, pairing.G2)
		checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
	}
	C.pairing_pp_apply(el.cptr, y.cptr, pairer.pp)
	return el
}

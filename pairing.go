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

struct pairing_s* newPairingStruct() { return malloc(sizeof(struct pairing_s)); }
void freePairingStruct(struct pairing_s* x) {
	pairing_clear(x);
	free(x);
}
*/
import "C"

import (
	"bytes"
	"io"
	"runtime"
)

// Field denotes the various possible algebraic structures associated with a
// pairing. G1, G2, and GT are the groups involved in the pairing operation. Zr
// is the field of integers with order r, where r is the order of G1, G2, and
// GT.
type Field int

const (
	G1 Field = iota
	G2 Field = iota
	GT Field = iota
	Zr Field = iota
)

// Pairing represents a pairing and its associated groups. The primary use of a
// pairing object is the initialization of group elements. Elements can be
// created in G1, G2, GT, or Zr. Additionally, elements can be checked or
// unchecked. See the Element type for more details.
type Pairing struct {
	params *Params // Prevents garbage collection
	cptr *C.struct_pairing_s
}

// NewPairing instantiates a pairing from a set of parameters.
func NewPairing(params *Params) *Pairing {
	pairing := makePairing(params)
	C.pairing_init_pbc_param(pairing.cptr, params.cptr)

	// pairing.params must point to params during the C call. Otherwise, the
	// garbage collector might free params.cptr (through the Params finalizer)
	// while pairing_init_pbc_param is still running! However, this function no
	// longer requires access to the parameters once control has returned. We
	// set pairing.params to nil to indicate that we no longer care if params
	// is garbage collected. If the caller still needs params (e.g., to
	// initialize more pairings), the garbage collector will know not to free
	// params even though pairing.params does not point to it.
	pairing.params = nil

	return pairing
}

// NewPairingFromReader loads pairing parameters from a Reader and instantiates
// a pairing.
func NewPairingFromReader(params io.Reader) (*Pairing, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(params)
	return NewPairingFromString(buf.String())
}

// NewPairingFromString loads pairing parameters from a string and instantiates
// a pairing.
func NewPairingFromString(params string) (*Pairing, error) {
	p, err := NewParamsFromString(params)
	if err != nil {
		return nil, err
	}
	return NewPairing(p), nil
}

// IsSymmetric returns true if G1 == G2 for this pairing.
func (pairing *Pairing) IsSymmetric() bool {
	return C.pairing_is_symmetric(pairing.cptr) != 0
}

// G1Length returns the size of elements in G1, in bytes.
func (pairing *Pairing) G1Length() uint {
	return uint(C.pairing_length_in_bytes_G1(pairing.cptr))
}

// G1XLength returns the size of X coordinates of elements in G1, in bytes.
func (pairing *Pairing) G1XLength() uint {
	return uint(C.pairing_length_in_bytes_x_only_G1(pairing.cptr))
}

// G1CompressedLength returns the size of compressed elements in G1, in bytes.
func (pairing *Pairing) G1CompressedLength() uint {
	return uint(C.pairing_length_in_bytes_compressed_G1(pairing.cptr))
}

// G2Length returns the size of elements in G2, in bytes.
func (pairing *Pairing) G2Length() uint {
	return uint(C.pairing_length_in_bytes_G2(pairing.cptr))
}

// G2XLength returns the size of X coordinates of elements in G2, in bytes.
func (pairing *Pairing) G2XLength() uint {
	return uint(C.pairing_length_in_bytes_x_only_G2(pairing.cptr))
}

// G2CompressedLength returns the size of compressed elements in G2, in bytes.
func (pairing *Pairing) G2CompressedLength() uint {
	return uint(C.pairing_length_in_bytes_compressed_G2(pairing.cptr))
}

// GTLength returns the size of elements in GT, in bytes.
func (pairing *Pairing) GTLength() uint {
	return uint(C.pairing_length_in_bytes_GT(pairing.cptr))
}

// ZrLength returns the size of elements in Zr, in bytes.
func (pairing *Pairing) ZrLength() uint {
	return uint(C.pairing_length_in_bytes_Zr(pairing.cptr))
}

// NewG1 creates a new checked element in G1.
func (pairing *Pairing) NewG1() *Element {
	return makeCheckedElement(pairing, G1, pairing.cptr.G1)
}

// NewG2 creates a new checked element in G2.
func (pairing *Pairing) NewG2() *Element {
	return makeCheckedElement(pairing, G2, pairing.cptr.G2)
}

// NewGT creates a new checked element in GT.
func (pairing *Pairing) NewGT() *Element {
	return makeCheckedElement(pairing, GT, &pairing.cptr.GT[0])
}

// NewZr creates a new checked element in Zr.
func (pairing *Pairing) NewZr() *Element {
	return makeCheckedElement(pairing, Zr, &pairing.cptr.Zr[0])
}

// NewUncheckedElement creates a new unchecked element in the target field.
// Unchecked elements are dangerous; see the Element documentation before
// deciding to use this method. It is safer to create elements using the NewG1,
// NewG2, NewGT, or NewZr methods.
func (pairing *Pairing) NewUncheckedElement(field Field) *Element {
	return makeUncheckedElement(pairing, true, field)
}

func clearPairing(pairing *Pairing) {
	C.freePairingStruct(pairing.cptr)
}

func makePairing(params *Params) *Pairing {
	pairing := &Pairing{
		params: params,
		cptr: C.newPairingStruct(),
	}
	runtime.SetFinalizer(pairing, clearPairing)
	return pairing
}

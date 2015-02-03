package pbc

/*
#include <pbc/pbc.h>
*/
import "C"

import (
	"bytes"
	"io"
	"runtime"
)

// Field denotes the various possible algebraic structures associated with a
// pairing. G1, G2, and GT are the groups involved in the pairing operation. Zr
// is the field of integers with order r, where r is the order of the groups.
type Field int

const (
	G1 Field = iota
	G2 Field = iota
	GT Field = iota
	Zr Field = iota
)

// Pairing represents a pairing and its associated groups. The primary use of a
// pairing object is the initialization of group elements. Elements can be
// created in G1, G2, GT, or Zr. Additionally, elements can either be checked
// or unchecked. Unchecked elements are slightly faster, but do not check to
// ensure that operations make sense. Checked elements defend against a variety
// of errors. For more details, see the Element documentation.
type Pairing struct {
	cptr *C.struct_pairing_s
}

// NewPairing instantiates a pairing from a set of parameters.
func NewPairing(params *Params) *Pairing {
	pairing := makePairing()
	C.pairing_init_pbc_param(pairing.cptr, params.cptr)
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

func (pairing *Pairing) G1Length() uint {
	return uint(C.pairing_length_in_bytes_G1(pairing.cptr))
}

func (pairing *Pairing) G1XLength() uint {
	return uint(C.pairing_length_in_bytes_x_only_G1(pairing.cptr))
}

func (pairing *Pairing) G1CompressedLength() uint {
	return uint(C.pairing_length_in_bytes_compressed_G1(pairing.cptr))
}

func (pairing *Pairing) G2Length() uint {
	return uint(C.pairing_length_in_bytes_G2(pairing.cptr))
}

func (pairing *Pairing) G2XLength() uint {
	return uint(C.pairing_length_in_bytes_x_only_G2(pairing.cptr))
}

func (pairing *Pairing) G2CompressedLength() uint {
	return uint(C.pairing_length_in_bytes_compressed_G2(pairing.cptr))
}

func (pairing *Pairing) GTLength() uint {
	return uint(C.pairing_length_in_bytes_GT(pairing.cptr))
}

func (pairing *Pairing) ZrLength() uint {
	return uint(C.pairing_length_in_bytes_Zr(pairing.cptr))
}

func (pairing *Pairing) NewG1() *Element {
	return makeCheckedElement(pairing, G1, pairing.cptr.G1)
}

func (pairing *Pairing) NewG2() *Element {
	return makeCheckedElement(pairing, G2, pairing.cptr.G2)
}

func (pairing *Pairing) NewGT() *Element {
	return makeCheckedElement(pairing, GT, &pairing.cptr.GT[0])
}

func (pairing *Pairing) NewZr() *Element {
	return makeCheckedElement(pairing, Zr, &pairing.cptr.Zr[0])
}

func (pairing *Pairing) NewElement(field Field) *Element {
	return makeUncheckedElement(pairing, true, field)
}

func clearPairing(pairing *Pairing) {
	C.pairing_clear(pairing.cptr)
}

func makePairing() *Pairing {
	pairing := &Pairing{cptr: &C.struct_pairing_s{}}
	runtime.SetFinalizer(pairing, clearPairing)
	return pairing
}

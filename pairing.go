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
type Pairing interface {
	// IsSymmetric returns true if G1 == G2 for this pairing.
	IsSymmetric() bool

	// Various methods to return the sizes of group elements.
	// *Length() returns the length of an element in bytes.
	// *XLength() returns the length of an X coordinate only, in bytes.
	// *CompressedLength() returns the length of a compressed element in bytes.
	G1Length() uint
	G1XLength() uint
	G1CompressedLength() uint
	G2Length() uint
	G2XLength() uint
	G2CompressedLength() uint
	GTLength() uint
	ZrLength() uint

	// Initialization methods for group elements.
	NewG1() Element
	NewG2() Element
	NewGT() Element
	NewZr() Element

	// Initializes an element without type checking.
	NewElement(Field) Element
}

type pairingImpl struct {
	data *C.struct_pairing_s
}

// NewPairing instantiates a pairing from a set of parameters.
func NewPairing(params Params) Pairing {
	pairing := makePairing()
	C.pairing_init_pbc_param(pairing.data, params.(*paramsImpl).data)
	return pairing
}

// NewPairingFromReader loads pairing parameters from a Reader and instantiates
// a pairing.
func NewPairingFromReader(params io.Reader) (Pairing, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(params)
	return NewPairingFromString(buf.String())
}

// NewPairingFromString loads pairing parameters from a string and instantiates
// a pairing.
func NewPairingFromString(params string) (Pairing, error) {
	p, err := NewParamsFromString(params)
	if err != nil {
		return nil, err
	}
	return NewPairing(p), nil
}

func (pairing *pairingImpl) IsSymmetric() bool {
	return C.pairing_is_symmetric(pairing.data) != 0
}

func (pairing *pairingImpl) G1Length() uint {
	return uint(C.pairing_length_in_bytes_G1(pairing.data))
}

func (pairing *pairingImpl) G1XLength() uint {
	return uint(C.pairing_length_in_bytes_x_only_G1(pairing.data))
}

func (pairing *pairingImpl) G1CompressedLength() uint {
	return uint(C.pairing_length_in_bytes_compressed_G1(pairing.data))
}

func (pairing *pairingImpl) G2Length() uint {
	return uint(C.pairing_length_in_bytes_G2(pairing.data))
}

func (pairing *pairingImpl) G2XLength() uint {
	return uint(C.pairing_length_in_bytes_x_only_G2(pairing.data))
}

func (pairing *pairingImpl) G2CompressedLength() uint {
	return uint(C.pairing_length_in_bytes_compressed_G2(pairing.data))
}

func (pairing *pairingImpl) GTLength() uint {
	return uint(C.pairing_length_in_bytes_GT(pairing.data))
}

func (pairing *pairingImpl) ZrLength() uint {
	return uint(C.pairing_length_in_bytes_Zr(pairing.data))
}

func (pairing *pairingImpl) NewG1() Element                 { return makeChecked(pairing, G1, pairing.data.G1) }
func (pairing *pairingImpl) NewG2() Element                 { return makeChecked(pairing, G2, pairing.data.G2) }
func (pairing *pairingImpl) NewGT() Element                 { return makeChecked(pairing, GT, &pairing.data.GT[0]) }
func (pairing *pairingImpl) NewZr() Element                 { return makeChecked(pairing, Zr, &pairing.data.Zr[0]) }
func (pairing *pairingImpl) NewElement(field Field) Element { return makeUnchecked(pairing, field) }

func clearPairing(pairing *pairingImpl) {
	C.pairing_clear(pairing.data)
}

func makePairing() *pairingImpl {
	pairing := &pairingImpl{data: &C.struct_pairing_s{}}
	runtime.SetFinalizer(pairing, clearPairing)
	return pairing
}

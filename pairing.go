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

type Field int

const (
	G1 Field = iota
	G2 Field = iota
	GT Field = iota
	Zr Field = iota
)

type Pairing interface {
	IsSymmetric() bool

	G1Length() uint
	G1XLength() uint
	G1CompressedLength() uint
	G2Length() uint
	G2XLength() uint
	G2CompressedLength() uint
	GTLength() uint
	ZrLength() uint

	NewG1() Element
	NewG2() Element
	NewGT() Element
	NewZr() Element
	NewElement(Field) Element
}

type pairingImpl struct {
	data *C.struct_pairing_s
}

func NewPairing(params io.Reader) (Pairing, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(params)
	return NewPairingFromString(buf.String())
}

func NewPairingFromString(params string) (Pairing, error) {
	p, err := NewParamsFromString(params)
	if err != nil {
		return nil, err
	}
	return NewPairingFromParams(p), nil
}

func NewPairingFromParams(params Params) Pairing {
	pairing := makePairing()
	C.pairing_init_pbc_param(pairing.data, params.(*paramsImpl).data)
	return pairing
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

func clearPairing(pairing *pairingImpl) {
	println("clearpairing")
	C.pairing_clear(pairing.data)
}

func makePairing() *pairingImpl {
	pairing := &pairingImpl{data: &C.struct_pairing_s{}}
	runtime.SetFinalizer(pairing, clearPairing)
	return pairing
}

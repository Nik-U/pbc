package pbc

/*
#include <pbc/pbc.h>
*/
import "C"

import (
	"fmt"
	"hash"
	"math/big"
	"runtime"
)

type Element interface {
	// NewFieldElement initializes a new element in the same group as this one.
	// The returned element will be unchecked if and only if this element is.
	NewFieldElement() Element

	// Methods to directly set the value of the element:
	Set0() Element
	Set1() Element
	SetInt32(int32) Element
	SetBig(*big.Int) Element
	Set(Element) Element

	// Methods to hash a value into a group element:
	SetFromHash([]byte) Element
	SetFromStringHash(s string, h hash.Hash) Element

	// SetString sets the value from an exported string representation.
	SetString(s string, base int) (Element, bool)

	// Methods to support the fmt package's Print and Scan functions:
	Format(fmt.State, rune)
	Scan(fmt.ScanState, rune) error

	// Methods to export the element in human-readable format:
	BigInt() *big.Int
	String() string

	// Methods to export the element as a sequence of bytes:
	BytesLen() int
	Bytes() []byte
	XBytesLen() int
	XBytes() []byte
	CompressedBytesLen() int
	CompressedBytes() []byte

	// Methods to import an element from a sequence of bytes:
	SetBytes([]byte) Element
	SetXBytes([]byte) Element
	SetCompressedBytes([]byte) Element

	// Methods to retrieve sub-elements (coordinates for points, coefficients
	// for polynomials):
	Len() int
	Item(int) Element
	X() *big.Int
	Y() *big.Int

	// Methods to determine the mathematical properties of the element:
	Is0() bool
	Is1() bool
	IsSquare() bool
	Sign() int

	// Methods to compare elements:
	Cmp(x Element) int
	Equals(x Element) bool

	// Methods to perform arithmetic operations. Not all operations are valid
	// for all groups.
	Add(x, y Element) Element
	Sub(x, y Element) Element
	Mul(x, y Element) Element
	MulBig(x Element, i *big.Int) Element
	MulInt32(x Element, i int32) Element
	MulZn(x, y Element) Element
	Div(x, y Element) Element
	Double(x Element) Element
	Halve(x Element) Element
	Square(x Element) Element
	Neg(x Element) Element
	Invert(x Element) Element

	// Methods to exponentiate elements:
	PowBig(x Element, i *big.Int) Element
	PowZn(x, i Element) Element
	Pow2Big(x Element, i *big.Int, y Element, j *big.Int) Element
	Pow2Zn(x, i, y, j Element) Element
	Pow3Big(x Element, i *big.Int, y Element, j *big.Int, z Element, k *big.Int) Element
	Pow3Zn(x, i, y, j, z, k Element) Element

	// Methods to perform pre-processed exponentiation:
	PreparePower() Power
	PowerBig(Power, *big.Int) Element
	PowerZn(Power, Element) Element

	// Methods to brute-force discrete logarithms in the group:
	BruteForceDL(g, h Element) Element
	PollardRhoDL(g, h Element) Element

	// Rand sets the element to a random group element.
	Rand() Element

	// Pairing operations:
	Pair(x, y Element) Element
	ProdPair(elements ...Element) Element
	ProdPairSlice(x, y []Element) Element

	// Methods to perform pre-processed pairing operations:
	PreparePairer() Pairer
	PairerPair(Pairer, Element) Element

	// Pairing returns the pairing associated with this element.
	Pairing() Pairing

	data() *C.struct_element_s
}

// Power stores pre-processed information to quickly exponentiate an element.
// A Power can be generated for Element x by calling x.PreparePower(). When
// PowBig or PowZn is called with Element target and integer i, the result of
// x^i will be stored in target.
type Power interface {
	PowBig(target Element, i *big.Int) Element
	PowZn(target Element, i Element) Element
}

type powerImpl struct {
	data *C.struct_element_pp_s
}

func (power *powerImpl) PowBig(target Element, i *big.Int) Element {
	return target.PowerBig(power, i)
}

func (power *powerImpl) PowZn(target Element, i Element) Element {
	return target.PowerZn(power, i)
}

func clearPower(power *powerImpl) {
	C.element_pp_clear(power.data)
}

func initPower(source Element) Power {
	power := &powerImpl{
		data: &C.struct_element_pp_s{},
	}
	C.element_pp_init(power.data, source.data())
	runtime.SetFinalizer(power, clearPower)
	return power
}

// Pairer stores pre-processed information to quickly pair an element. A Pairer
// can be generated for Element x by calling x.PreparePairer(). When Pair is
// called with Elements target and y, the result of e(x,y) will be stored in
// target.
type Pairer interface {
	Pair(target Element, y Element) Element
}

type pairerImpl struct {
	source Element
	data   *C.struct_pairing_pp_s
}

func (pairer *pairerImpl) Pair(target Element, y Element) Element {
	return target.PairerPair(pairer, y)
}

func clearPairer(pairer *pairerImpl) {
	C.pairing_pp_clear(pairer.data)
}

func initPairer(source Element) Pairer {
	pairer := &pairerImpl{
		source: source,
		data:   &C.struct_pairing_pp_s{},
	}
	C.pairing_pp_init(pairer.data, source.data(), source.Pairing().(*pairingImpl).data)
	runtime.SetFinalizer(pairer, clearPairer)
	return pairer
}

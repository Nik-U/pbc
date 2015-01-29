package pbc

/*
#include <pbc/pbc.h>
*/
import "C"

import (
	"errors"
	"math/big"
)

var (
	ErrIllegalOp    = errors.New("operation is illegal for elements of this type")
	ErrUncheckedOp  = errors.New("unchecked element passed to checked operation")
	ErrIncompatible = errors.New("elements are from incompatible fields or pairings")
	ErrOutOfRange   = errors.New("index out of range")
	ErrInternal     = errors.New("a severe internal error has lead to possible memory corruption")
)

func (el *checkedElement) impl() *elementImpl { return &el.unchecked }

func element2Checked(x Element) *checkedElement {
	checked, ok := x.(*checkedElement)
	if !ok {
		panic(ErrUncheckedOp)
	}
	return checked
}

func checkFieldsMatch(f1, f2 *C.struct_field_s) {
	if f1 != f2 {
		panic(ErrIncompatible)
	}
}

func (el *checkedElement) checkCompatible(other Element) {
	otherChecked := element2Checked(other)
	checkFieldsMatch(el.fieldPtr, otherChecked.fieldPtr)
}

func (el *checkedElement) checkAllCompatible(elements ...Element) {
	for _, other := range elements {
		el.checkCompatible(other)
	}
}

func (el *checkedElement) checkInteger() {
	if !el.isInteger {
		panic(ErrIllegalOp)
	}
}

func (el *checkedElement) NewFieldElement() Element {
	newElement := &checkedElement{}
	*newElement = *el
	initElement(&newElement.unchecked, el.unchecked.pairing, false, G1)
	C.element_init_same_as(newElement.unchecked.data, el.unchecked.data)
	return newElement
}

func (el *checkedElement) Set0() Element {
	el.unchecked.Set0()
	return el
}

func (el *checkedElement) Set1() Element {
	el.unchecked.Set1()
	return el
}

func (el *checkedElement) SetInt32(i int32) Element {
	el.checkInteger()
	el.unchecked.SetInt32(i)
	return el
}

func (el *checkedElement) SetBig(i *big.Int) Element {
	el.checkInteger()
	el.unchecked.SetBig(i)
	return el
}

func (el *checkedElement) Set(src Element) Element {
	el.checkCompatible(src)
	el.unchecked.Set(src)
	return el
}

func (el *checkedElement) SetFromHash(hash []byte) Element {
	el.unchecked.SetFromHash(hash)
	return el
}

func (el *checkedElement) SetBytes(buf []byte) Element {
	el.unchecked.SetBytes(buf)
	return el
}

func (el *checkedElement) SetXBytes(buf []byte) Element {
	el.unchecked.SetXBytes(buf)
	return el
}

func (el *checkedElement) SetCompressedBytes(buf []byte) Element {
	el.unchecked.SetCompressedBytes(buf)
	return el
}

func (el *checkedElement) SetString(s string, base int) (Element, bool) {
	_, ok := el.unchecked.SetString(s, base)
	return el, ok
}

func (el *checkedElement) BigInt() *big.Int {
	el.checkInteger()
	return el.unchecked.BigInt()
}

func (el *checkedElement) BytesLen() int           { return el.unchecked.BytesLen() }
func (el *checkedElement) XBytesLen() int          { return el.unchecked.XBytesLen() }
func (el *checkedElement) CompressedBytesLen() int { return el.unchecked.CompressedBytesLen() }

func checkedWrite(bytesWritten C.int, buffer []byte) []byte {
	if int64(bytesWritten) > int64(len(buffer)) {
		panic(ErrInternal)
	}
	return buffer
}

func (el *checkedElement) Bytes() []byte {
	buf := make([]byte, el.BytesLen())
	return checkedWrite(el.unchecked.writeBytes(buf), buf)
}

func (el *checkedElement) XBytes() []byte {
	buf := make([]byte, el.XBytesLen())
	return checkedWrite(el.unchecked.writeXBytes(buf), buf)
}

func (el *checkedElement) CompressedBytes() []byte {
	buf := make([]byte, el.CompressedBytesLen())
	return checkedWrite(el.unchecked.writeCompressedBytes(buf), buf)
}

func (el *checkedElement) Len() int { return el.unchecked.Len() }

func (el *checkedElement) Item(i int) Element {
	if i >= el.Len() {
		panic(ErrOutOfRange)
	}
	uncheckedData := el.unchecked.Item(i).(*elementImpl)
	item := &checkedElement{
		fieldPtr:  uncheckedData.data.field,
		isInteger: uncheckedData.Len() == 0,
	}
	item.unchecked = *uncheckedData
	return item
}

func (el *checkedElement) X() *big.Int    { return el.unchecked.X() }
func (el *checkedElement) Y() *big.Int    { return el.unchecked.Y() }
func (el *checkedElement) Is0() bool      { return el.unchecked.Is0() }
func (el *checkedElement) Is1() bool      { return el.unchecked.Is1() }
func (el *checkedElement) IsSquare() bool { return el.unchecked.IsSquare() }
func (el *checkedElement) Sign() int      { return el.unchecked.Sign() }

func (el *checkedElement) Cmp(x Element) int {
	el.checkCompatible(x)
	return el.unchecked.Cmp(x)
}

func (el *checkedElement) Add(x Element, y Element) Element {
	el.checkAllCompatible(x, y)
	el.unchecked.Add(x, y)
	return el
}

func (el *checkedElement) Sub(x, y Element) Element {
	el.checkAllCompatible(x, y)
	el.unchecked.Sub(x, y)
	return el
}

func (el *checkedElement) Mul(x, y Element) Element {
	el.checkAllCompatible(x, y)
	el.unchecked.Mul(x, y)
	return el
}

func (el *checkedElement) MulBig(x Element, i *big.Int) Element {
	el.checkCompatible(x)
	el.unchecked.MulBig(x, i)
	return el
}

func (el *checkedElement) MulInt32(x Element, i int32) Element {
	el.checkCompatible(x)
	el.unchecked.MulInt32(x, i)
	return el
}

func (el *checkedElement) MulZn(x, y Element) Element {
	el.checkCompatible(x)
	element2Checked(y).checkInteger()
	el.unchecked.MulZn(x, y)
	return el
}

func (el *checkedElement) Div(x, y Element) Element {
	el.checkAllCompatible(x, y)
	el.unchecked.Div(x, y)
	return el
}

func (el *checkedElement) Double(x Element) Element {
	el.checkCompatible(x)
	el.unchecked.Double(x)
	return el
}

func (el *checkedElement) Halve(x Element) Element {
	el.checkCompatible(x)
	el.unchecked.Halve(x)
	return el
}

func (el *checkedElement) Square(x Element) Element {
	el.checkCompatible(x)
	el.unchecked.Square(x)
	return el
}

func (el *checkedElement) Neg(x Element) Element {
	el.checkCompatible(x)
	el.unchecked.Neg(x)
	return el
}

func (el *checkedElement) Invert(x Element) Element {
	el.checkCompatible(x)
	el.unchecked.Invert(x)
	return el
}

func (el *checkedElement) PowBig(x Element, i *big.Int) Element {
	el.checkCompatible(x)
	el.unchecked.PowBig(x, i)
	return el
}

func (el *checkedElement) PowZn(x, i Element) Element {
	el.checkCompatible(x)
	element2Checked(i).checkInteger()
	el.unchecked.PowZn(x, i)
	return el
}

func (el *checkedElement) Pow2Big(x Element, i *big.Int, y Element, j *big.Int) Element {
	el.checkAllCompatible(x, y)
	el.unchecked.Pow2Big(x, i, y, j)
	return el
}

func (el *checkedElement) Pow2Zn(x, i, y, j Element) Element {
	el.checkAllCompatible(x, y)
	element2Checked(i).checkInteger()
	element2Checked(j).checkInteger()
	el.unchecked.Pow2Zn(x, i, y, j)
	return el
}

func (el *checkedElement) Pow3Big(x Element, i *big.Int, y Element, j *big.Int, z Element, k *big.Int) Element {
	el.checkAllCompatible(x, y, z)
	el.unchecked.Pow3Big(x, i, y, j, z, k)
	return el
}

func (el *checkedElement) Pow3Zn(x, i, y, j, z, k Element) Element {
	el.checkAllCompatible(x, y, z)
	element2Checked(i).checkInteger()
	element2Checked(j).checkInteger()
	element2Checked(k).checkInteger()
	el.unchecked.Pow3Zn(x, i, y, j, z, k)
	return el
}

func (el *checkedElement) PreparePower() Power { return initPower(el) }

func (el *checkedElement) PowerBig(power Power, i *big.Int) Element {
	el.unchecked.PowerBig(power, i)
	return el
}

func (el *checkedElement) PowerZn(power Power, i Element) Element {
	element2Checked(i).checkInteger()
	el.unchecked.PowerZn(power, i)
	return el
}

func (el *checkedElement) Pair(x, y Element) Element {
	pairing := el.unchecked.pairing.data
	checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
	checkFieldsMatch(element2Checked(x).fieldPtr, pairing.G1)
	checkFieldsMatch(element2Checked(y).fieldPtr, pairing.G2)
	el.unchecked.Pair(x, y)
	return el
}

func (el *checkedElement) ProdPair(elements ...Element) Element {
	pairing := el.unchecked.pairing.data
	checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
	n := len(elements)
	for i := 1; i < n; i += 2 {
		checkFieldsMatch(element2Checked(elements[i-1]).fieldPtr, pairing.G1)
		checkFieldsMatch(element2Checked(elements[i]).fieldPtr, pairing.G2)
	}
	el.unchecked.ProdPair(elements...)
	return el
}

func (el *checkedElement) ProdPairSlice(x, y []Element) Element {
	pairing := el.unchecked.pairing.data
	checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
	n := len(x)
	for i := 1; i < n; i++ {
		checkFieldsMatch(element2Checked(x[i]).fieldPtr, pairing.G1)
	}
	n = len(y)
	for i := 1; i < n; i++ {
		checkFieldsMatch(element2Checked(y[i]).fieldPtr, pairing.G2)

	}
	el.unchecked.ProdPairSlice(x, y)
	return el
}

func (el *checkedElement) PreparePairer() Pairer { return initPairer(el) }

func (el *checkedElement) PairerPair(pairer Pairer, x Element) Element {
	in1 := element2Checked(pairer.(*pairerImpl).source)
	in2 := element2Checked(x)
	pairing := el.unchecked.pairing.data
	checkFieldsMatch(in1.fieldPtr, pairing.G1)
	checkFieldsMatch(in2.fieldPtr, pairing.G2)
	checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
	el.unchecked.PairerPair(pairer, x)
	return el
}

func (el *checkedElement) BruteForceDL(g, h Element) Element {
	el.checkInteger()
	element2Checked(g).checkCompatible(h)
	el.unchecked.BruteForceDL(g, h)
	return el
}

func (el *checkedElement) PollardRhoDL(g, h Element) Element {
	el.checkInteger()
	element2Checked(g).checkCompatible(h)
	el.unchecked.PollardRhoDL(g, h)
	return el
}

func (el *checkedElement) Rand() Element {
	el.unchecked.Rand()
	return el
}

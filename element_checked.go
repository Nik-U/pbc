package pbc

/*
#cgo LDFLAGS: /usr/local/lib/libpbc.a -lgmp
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

func (el *checkedElement) impl() *elementImpl { return &el.elementImpl }

func element2Checked(x Element) *checkedElement {
	checked, ok := x.(*checkedElement)
	if !ok {
		panic(ErrUncheckedOp)
	}
	return checked
}

func (el *checkedElement) checkCompatible(other Element) {
	otherChecked := element2Checked(other)
	if el.fieldPtr != otherChecked.fieldPtr {
		panic(ErrIncompatible)
	}
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
	initElement(&newElement.elementImpl, el.pairing, false, G1)
	C.element_init_same_as(newElement.elementImpl.data, el.data)
	return newElement
}

func (el *checkedElement) SetInt32(i int32) Element {
	el.checkInteger()
	return el.elementImpl.SetInt32(i)
}

func (el *checkedElement) SetBig(i *big.Int) Element {
	el.checkInteger()
	return el.elementImpl.SetBig(i)
}

func (el *checkedElement) Set(src Element) Element {
	el.checkCompatible(src)
	return el.elementImpl.Set(src)
}

func (el *checkedElement) BigInt() *big.Int {
	el.checkInteger()
	return el.elementImpl.BigInt()
}

func checkedWrite(bytesWritten C.int, buffer []byte) []byte {
	if int64(bytesWritten) > int64(len(buffer)) {
		panic(ErrInternal)
	}
	return buffer
}

func (el *checkedElement) Bytes() []byte {
	buf := make([]byte, el.BytesLen())
	return checkedWrite(el.elementImpl.writeBytes(buf), buf)
}

func (el *checkedElement) XBytes() []byte {
	buf := make([]byte, el.XBytesLen())
	return checkedWrite(el.elementImpl.writeXBytes(buf), buf)
}

func (el *checkedElement) CompressedBytes() []byte {
	buf := make([]byte, el.CompressedBytesLen())
	return checkedWrite(el.elementImpl.writeCompressedBytes(buf), buf)
}

func (el *checkedElement) Item(i int) Element {
	if i >= el.Len() {
		panic(ErrOutOfRange)
	}
	uncheckedData := el.elementImpl.Item(i).(*elementImpl)
	item := &checkedElement{
		fieldPtr:  uncheckedData.data.field,
		isInteger: uncheckedData.Len() == 0,
	}
	item.elementImpl = *uncheckedData
	return item
}

func (el *checkedElement) Cmp(x Element) int {
	el.checkCompatible(x)
	return el.elementImpl.Cmp(x)
}

func (el *checkedElement) Add(x Element, y Element) Element {
	el.checkAllCompatible(x, y)
	return el.elementImpl.Add(x, y)
}

func (el *checkedElement) Sub(x, y Element) Element {
	el.checkAllCompatible(x, y)
	return el.elementImpl.Sub(x, y)
}

func (el *checkedElement) Mul(x, y Element) Element {
	el.checkAllCompatible(x, y)
	return el.elementImpl.Mul(x, y)
}

func (el *checkedElement) MulBig(x Element, i *big.Int) Element {
	el.checkCompatible(x)
	return el.elementImpl.MulBig(x, i)
}

func (el *checkedElement) MulInt32(x Element, i int32) Element {
	el.checkCompatible(x)
	return el.elementImpl.MulInt32(x, i)
}

func (el *checkedElement) MulZn(x, y Element) Element {
	el.checkCompatible(x)
	element2Checked(y).checkInteger()
	return el.elementImpl.MulZn(x, y)
}

func (el *checkedElement) Div(x, y Element) Element {
	el.checkAllCompatible(x, y)
	return el.elementImpl.Div(x, y)
}

func (el *checkedElement) Double(x Element) Element {
	el.checkCompatible(x)
	return el.elementImpl.Double(x)
}

func (el *checkedElement) Halve(x Element) Element {
	el.checkCompatible(x)
	return el.elementImpl.Halve(x)
}

func (el *checkedElement) Square(x Element) Element {
	el.checkCompatible(x)
	return el.elementImpl.Square(x)
}

func (el *checkedElement) Neg(x Element) Element {
	el.checkCompatible(x)
	return el.elementImpl.Neg(x)
}

func (el *checkedElement) Invert(x Element) Element {
	el.checkCompatible(x)
	return el.elementImpl.Invert(x)
}

func (el *checkedElement) PowBig(x Element, i *big.Int) Element {
	el.checkCompatible(x)
	return el.elementImpl.PowBig(x, i)
}

func (el *checkedElement) PowZn(x, i Element) Element {
	el.checkCompatible(x)
	element2Checked(i).checkInteger()
	return el.elementImpl.PowZn(x, i)
}

func (el *checkedElement) Pow2Big(x Element, i *big.Int, y Element, j *big.Int) Element {
	el.checkAllCompatible(x, y)
	return el.elementImpl.Pow2Big(x, i, y, j)
}

func (el *checkedElement) Pow2Zn(x, i, y, j Element) Element {
	el.checkAllCompatible(x, y)
	element2Checked(i).checkInteger()
	element2Checked(j).checkInteger()
	return el.elementImpl.Pow2Zn(x, i, y, j)
}

func (el *checkedElement) Pow3Big(x Element, i *big.Int, y Element, j *big.Int, z Element, k *big.Int) Element {
	el.checkAllCompatible(x, y, z)
	return el.elementImpl.Pow3Big(x, i, y, j, z, k)
}

func (el *checkedElement) Pow3Zn(x, i, y, j, z, k Element) Element {
	el.checkAllCompatible(x, y, z)
	element2Checked(i).checkInteger()
	element2Checked(j).checkInteger()
	element2Checked(k).checkInteger()
	return el.elementImpl.Pow3Zn(x, i, y, j, z, k)
}

func (el *checkedElement) PreparePower() Power {
	power := &checkedPower{}
	initPower(&power.powerImpl, &el.elementImpl)
	return power
}

func (power *checkedPower) PowZn(i Element) Element {
	element2Checked(i).checkInteger()
	return power.powerImpl.PowZn(i)
}

func (el *checkedElement) BruteForceDL(g, h Element) Element {
	el.checkInteger()
	element2Checked(g).checkCompatible(h)
	return el.elementImpl.BruteForceDL(g, h)
}

func (el *checkedElement) PollardRhoDL(g, h Element) Element {
	el.checkInteger()
	element2Checked(g).checkCompatible(h)
	return el.elementImpl.PollardRhoDL(g, h)
}

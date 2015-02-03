package pbc

/*
#include <pbc/pbc.h>
*/
import "C"
import (
	"hash"
	"math/big"
	"unsafe"
)

func (el *Element) BigInt() *big.Int {
	if el.checked {
		el.checkInteger()
	}
	mpz := newMpz()
	C.element_to_mpz(&mpz[0], el.cptr)
	return mpz2big(mpz)
}

func (el *Element) Set(src *Element) *Element {
	if el.checked {
		el.checkCompatible(src)
	}
	C.element_set(el.cptr, src.cptr)
	return el
}

func (el *Element) SetInt32(i int32) *Element {
	if el.checked {
		el.checkInteger()
	}
	C.element_set_si(el.cptr, C.long(i))
	return el
}

func (el *Element) SetBig(i *big.Int) *Element {
	if el.checked {
		el.checkInteger()
	}
	C.element_set_mpz(el.cptr, &big2mpz(i)[0])
	return el
}

func (el *Element) SetFromHash(hash []byte) *Element {
	C.element_from_hash(el.cptr, unsafe.Pointer(&hash[0]), C.int(len(hash)))
	return el
}

func (el *Element) SetFromStringHash(s string, h hash.Hash) *Element {
	h.Reset()
	if _, err := h.Write([]byte(s)); err != nil {
		panic(ErrHashFailure)
	}
	return el.SetFromHash(h.Sum([]byte{}))
}

func (el *Element) BytesLen() int {
	return int(C.element_length_in_bytes(el.cptr))
}

func (el *Element) Bytes() []byte {
	buf := make([]byte, el.BytesLen())
	written := C.element_to_bytes((*C.uchar)(unsafe.Pointer(&buf[0])), el.cptr)
	if int64(written) > int64(len(buf)) {
		panic(ErrInternal)
	}
	return buf
}

func (el *Element) SetBytes(buf []byte) *Element {
	C.element_from_bytes(el.cptr, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

func (el *Element) XBytesLen() int {
	return int(C.element_length_in_bytes_x_only(el.cptr))
}

func (el *Element) XBytes() []byte {
	buf := make([]byte, el.XBytesLen())
	written := C.element_to_bytes_x_only((*C.uchar)(unsafe.Pointer(&buf[0])), el.cptr)
	if int64(written) > int64(len(buf)) {
		panic(ErrInternal)
	}
	return buf
}

func (el *Element) SetXBytes(buf []byte) *Element {
	C.element_from_bytes_x_only(el.cptr, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

func (el *Element) CompressedBytesLen() int {
	return int(C.element_length_in_bytes_compressed(el.cptr))
}

func (el *Element) CompressedBytes() []byte {
	buf := make([]byte, el.CompressedBytesLen())
	written := C.element_to_bytes_compressed((*C.uchar)(unsafe.Pointer(&buf[0])), el.cptr)
	if int64(written) > int64(len(buf)) {
		panic(ErrInternal)
	}
	return buf
}

func (el *Element) SetCompressedBytes(buf []byte) *Element {
	C.element_from_bytes_compressed(el.cptr, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

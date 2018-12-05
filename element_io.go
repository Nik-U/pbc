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
*/
import "C"

import (
	"hash"
	"math/big"
	"unsafe"
)

// BigInt converts the Element to a big.Int if such an operation makes sense.
// Note that elliptic curve points cannot be converted using this method, even
// though this is possible in the original PBC library. If callers wish to
// convert the first coordinate into an integer, they should explicitly call
// X().
//
// Requirements:
// el is expressible as an integer (e.g., an element of Zn, but not a point).
func (el *Element) BigInt() *big.Int {
	if el.checked {
		el.checkInteger()
	}
	m := newMpz()
	C.element_to_mpz(&m.i[0], el.cptr)
	return mpz2big(m)
}

// Set sets the value of el to be the same as src.
//
// Requirements:
// el and src must be from the same algebraic structure.
func (el *Element) Set(src *Element) *Element {
	if el.checked {
		el.checkCompatible(src)
	}
	C.element_set(el.cptr, src.cptr)
	return el
}

// SetInt32 sets the value of el to the integer i. This operation is only valid
// for elements in integer fields (e.g., Zr for a pairing).
//
// Requirements:
// el must be an element of an integer mod ring (e.g., Zn for some n).
func (el *Element) SetInt32(i int32) *Element {
	if el.checked {
		el.checkInteger()
	}
	C.element_set_si(el.cptr, C.long(i))
	return el
}

// SetBig sets the value of el to the integer i. This operation is only valid
// for elements in integer fields (e.g., Zr for a pairing).
//
// Requirements:
// el must be an element of an integer mod ring (e.g., Zn for some n).
func (el *Element) SetBig(i *big.Int) *Element {
	if el.checked {
		el.checkInteger()
	}
	C.element_set_mpz(el.cptr, &big2mpz(i).i[0])
	return el
}

// SetFromHash generates el deterministically from the bytes in hash.
func (el *Element) SetFromHash(hash []byte) *Element {
	C.element_from_hash(el.cptr, unsafe.Pointer(&hash[0]), C.int(len(hash)))
	return el
}

// SetFromStringHash hashes s with h and then calls SetFromHash. h may or may
// not be a cryptographic hash, depending on the higher level protocol.
func (el *Element) SetFromStringHash(s string, h hash.Hash) *Element {
	h.Reset()
	if _, err := h.Write([]byte(s)); err != nil {
		panic(ErrHashFailure)
	}
	return el.SetFromHash(h.Sum([]byte{}))
}

// BytesLen returns the number of bytes needed to represent el.
func (el *Element) BytesLen() int {
	return int(C.element_length_in_bytes(el.cptr))
}

// Bytes exports el as a byte sequence.
func (el *Element) Bytes() []byte {
	buf := make([]byte, el.BytesLen())
	written := C.element_to_bytes((*C.uchar)(unsafe.Pointer(&buf[0])), el.cptr)
	if int64(written) > int64(len(buf)) {
		panic(ErrInternal)
	}
	return buf
}

// SetBytes imports a sequence exported by Bytes() and sets the value of el.
func (el *Element) SetBytes(buf []byte) *Element {
	C.element_from_bytes(el.cptr, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

// XBytesLen returns the number of bytes needed to represent el's X coordinate.
//
// Requirements:
// el must be a point on an elliptic curve.
func (el *Element) XBytesLen() int {
	if el.checked {
		el.checkPoint()
	}
	return int(C.element_length_in_bytes_x_only(el.cptr))
}

// XBytes exports el's X coordinate as a byte sequence.
//
// Requirements:
// el must be a point on an elliptic curve.
func (el *Element) XBytes() []byte {
	if el.checked {
		el.checkPoint()
	}
	buf := make([]byte, el.XBytesLen())
	written := C.element_to_bytes_x_only((*C.uchar)(unsafe.Pointer(&buf[0])), el.cptr)
	if int64(written) > int64(len(buf)) {
		panic(ErrInternal)
	}
	return buf
}

// SetXBytes imports a sequence exported by XBytes() and sets el to be a point
// on the curve with the given X coordinate. In general, this point is not
// unique. For each X coordinate, there exist two different points (for the
// pairings in PBC), and they are inverses of each other. An application can
// deal with this by either exporting the sign of the element along with the X
// coordinate, or by testing the value to see if it makes sense in the higher
// level protocol (and inverting it if it does not).
//
// Requirements:
// el must be a point on an elliptic curve.
func (el *Element) SetXBytes(buf []byte) *Element {
	if el.checked {
		el.checkPoint()
	}
	C.element_from_bytes_x_only(el.cptr, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

// CompressedBytesLen returns the number of bytes needed to represent a
// compressed form of el.
//
// Requirements:
// el must be a point on an elliptic curve.
func (el *Element) CompressedBytesLen() int {
	if el.checked {
		el.checkPoint()
	}
	return int(C.element_length_in_bytes_compressed(el.cptr))
}

// CompressedBytes exports el in a compressed form as a byte sequence.
//
// Requirements:
// el must be a point on an elliptic curve.
func (el *Element) CompressedBytes() []byte {
	if el.checked {
		el.checkPoint()
	}
	buf := make([]byte, el.CompressedBytesLen())
	written := C.element_to_bytes_compressed((*C.uchar)(unsafe.Pointer(&buf[0])), el.cptr)
	if int64(written) > int64(len(buf)) {
		panic(ErrInternal)
	}
	return buf
}

// SetCompressedBytes imports a sequence exported by CompressedBytes() and sets
// the value of el.
//
// Requirements:
// el must be a point on an elliptic curve.
func (el *Element) SetCompressedBytes(buf []byte) *Element {
	if el.checked {
		el.checkPoint()
	}
	C.element_from_bytes_compressed(el.cptr, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

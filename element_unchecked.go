package pbc

/*
#include <pbc/pbc.h>
*/
import "C"

import (
	"errors"
	"math/big"
	"unsafe"
)

var ErrBadPairList = errors.New("pairing product list is in an invalid format")

func (el *elementImpl) impl() *elementImpl { return el }

func (el *elementImpl) NewFieldElement() Element {
	newElement := &elementImpl{}
	initElement(newElement, el.pairing, false, G1)
	C.element_init_same_as(newElement.data, el.data)
	return newElement
}

func (el *elementImpl) Set0() Element {
	C.element_set0(el.data)
	return el
}

func (el *elementImpl) Set1() Element {
	C.element_set1(el.data)
	return el
}

func (el *elementImpl) SetInt32(i int32) Element {
	C.element_set_si(el.data, C.long(i))
	return el
}

func (el *elementImpl) SetBig(i *big.Int) Element {
	C.element_set_mpz(el.data, &big2mpz(i)[0])
	return el
}

func (el *elementImpl) Set(src Element) Element {
	C.element_set(el.data, src.impl().data)
	return el
}

func (el *elementImpl) SetFromHash(hash []byte) Element {
	C.element_from_hash(el.data, unsafe.Pointer(&hash[0]), C.int(len(hash)))
	return el
}

func (el *elementImpl) SetBytes(buf []byte) Element {
	C.element_from_bytes(el.data, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

func (el *elementImpl) SetXBytes(buf []byte) Element {
	C.element_from_bytes_x_only(el.data, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

func (el *elementImpl) SetCompressedBytes(buf []byte) Element {
	C.element_from_bytes_compressed(el.data, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

func (el *elementImpl) SetString(s string, base int) (Element, bool) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))

	if ok := C.element_set_str(el.data, cstr, C.int(base)); ok == 0 {
		return nil, false
	}
	return el, true
}

func (el *elementImpl) BigInt() *big.Int {
	mpz := newmpz()
	C.element_to_mpz(&mpz[0], el.data)
	return mpz2big(mpz)
}

func (el *elementImpl) BytesLen() int {
	return int(C.element_length_in_bytes(el.data))
}

func (el *elementImpl) writeBytes(buf []byte) C.int {
	return C.element_to_bytes((*C.uchar)(unsafe.Pointer(&buf[0])), el.data)
}

func (el *elementImpl) Bytes() []byte {
	buf := make([]byte, el.BytesLen())
	el.writeBytes(buf)
	return buf
}

func (el *elementImpl) XBytesLen() int {
	return int(C.element_length_in_bytes_x_only(el.data))
}

func (el *elementImpl) writeXBytes(buf []byte) C.int {
	return C.element_to_bytes_x_only((*C.uchar)(unsafe.Pointer(&buf[0])), el.data)
}

func (el *elementImpl) XBytes() []byte {
	buf := make([]byte, el.XBytesLen())
	el.writeXBytes(buf)
	return buf
}

func (el *elementImpl) CompressedBytesLen() int {
	return int(C.element_length_in_bytes_compressed(el.data))
}

func (el *elementImpl) writeCompressedBytes(buf []byte) C.int {
	return C.element_to_bytes_compressed((*C.uchar)(unsafe.Pointer(&buf[0])), el.data)
}

func (el *elementImpl) CompressedBytes() []byte {
	buf := make([]byte, el.CompressedBytesLen())
	el.writeCompressedBytes(buf)
	return buf
}

func (el *elementImpl) Len() int {
	return int(C.element_item_count(el.data))
}

func (el *elementImpl) Item(i int) Element {
	return &elementImpl{
		pairing: el.pairing,
		data:    C.element_item(el.data, C.int(i)),
	}
}

func (el *elementImpl) X() *big.Int {
	return el.Item(0).BigInt()
}

func (el *elementImpl) Y() *big.Int {
	return el.Item(1).BigInt()
}

func (el *elementImpl) Is0() bool {
	return C.element_is0(el.data) != 0
}

func (el *elementImpl) Is1() bool {
	return C.element_is1(el.data) != 0
}

func (el *elementImpl) IsSquare() bool {
	return C.element_is_sqr(el.data) != 0
}

func normalizeSign(sign int64) int {
	if sign > 0 {
		return 1
	}
	if sign < 0 {
		return -1
	}
	return 0
}

func (el *elementImpl) Sign() int {
	return normalizeSign(int64(C.element_sign(el.data)))
}

func (el *elementImpl) Cmp(x Element) int {
	return normalizeSign(int64(C.element_cmp(el.data, x.impl().data)))
}

func (el *elementImpl) Add(x, y Element) Element {
	C.element_add(el.data, x.impl().data, y.impl().data)
	return el
}

func (el *elementImpl) Sub(x, y Element) Element {
	C.element_sub(el.data, x.impl().data, y.impl().data)
	return el
}

func (el *elementImpl) Mul(x, y Element) Element {
	C.element_mul(el.data, x.impl().data, y.impl().data)
	return el
}

func (el *elementImpl) MulBig(x Element, i *big.Int) Element {
	C.element_mul_mpz(el.data, x.impl().data, &big2mpz(i)[0])
	return el
}

func (el *elementImpl) MulInt32(x Element, i int32) Element {
	C.element_mul_si(el.data, x.impl().data, C.long(i))
	return el
}

func (el *elementImpl) MulZn(x, y Element) Element {
	C.element_mul_zn(el.data, x.impl().data, y.impl().data)
	return el
}

func (el *elementImpl) Div(x, y Element) Element {
	C.element_div(el.data, x.impl().data, y.impl().data)
	return el
}

func (el *elementImpl) Double(x Element) Element {
	C.element_double(el.data, x.impl().data)
	return el
}

func (el *elementImpl) Halve(x Element) Element {
	C.element_halve(el.data, x.impl().data)
	return el
}

func (el *elementImpl) Square(x Element) Element {
	C.element_square(el.data, x.impl().data)
	return el
}

func (el *elementImpl) Neg(x Element) Element {
	C.element_neg(el.data, x.impl().data)
	return el
}

func (el *elementImpl) Invert(x Element) Element {
	C.element_invert(el.data, x.impl().data)
	return el
}

func (el *elementImpl) PowBig(x Element, i *big.Int) Element {
	C.element_pow_mpz(el.data, x.impl().data, &big2mpz(i)[0])
	return el
}

func (el *elementImpl) PowZn(x, i Element) Element {
	C.element_pow_zn(el.data, x.impl().data, i.impl().data)
	return el
}

func (el *elementImpl) Pow2Big(x Element, i *big.Int, y Element, j *big.Int) Element {
	C.element_pow2_mpz(el.data, x.impl().data, &big2mpz(i)[0], y.impl().data, &big2mpz(j)[0])
	return el
}

func (el *elementImpl) Pow2Zn(x, i, y, j Element) Element {
	C.element_pow2_zn(el.data, x.impl().data, i.impl().data, y.impl().data, j.impl().data)
	return el
}

func (el *elementImpl) Pow3Big(x Element, i *big.Int, y Element, j *big.Int, z Element, k *big.Int) Element {
	C.element_pow3_mpz(el.data, x.impl().data, &big2mpz(i)[0], y.impl().data, &big2mpz(j)[0], z.impl().data, &big2mpz(k)[0])
	return el
}

func (el *elementImpl) Pow3Zn(x, i, y, j, z, k Element) Element {
	C.element_pow3_zn(el.data, x.impl().data, i.impl().data, y.impl().data, j.impl().data, z.impl().data, k.impl().data)
	return el
}

func (el *elementImpl) PreparePower() Power { return initPower(el) }

func (el *elementImpl) PowerBig(power Power, i *big.Int) Element {
	C.element_pp_pow(el.data, &big2mpz(i)[0], power.(*powerImpl).data)
	return el
}

func (el *elementImpl) PowerZn(power Power, i Element) Element {
	C.element_pp_pow_zn(el.data, i.impl().data, power.(*powerImpl).data)
	return el
}

func (el *elementImpl) Pair(x, y Element) Element {
	C.pairing_apply(el.data, x.impl().data, y.impl().data, el.pairing.data)
	return el
}

func (el *elementImpl) doProdPair(in1, in2 []C.struct_element_s) Element {
	x := (*C.element_t)(unsafe.Pointer(&in1[0]))
	y := (*C.element_t)(unsafe.Pointer(&in2[0]))
	C.element_prod_pairing(el.data, x, y, C.int(len(in1)))
	return el
}

func (el *elementImpl) ProdPair(elements ...Element) Element {
	n := len(elements)
	if n%2 != 0 {
		panic(ErrBadPairList)
	}
	half := n / 2
	in1 := make([]C.struct_element_s, half)
	in2 := make([]C.struct_element_s, half)
	for i, j := 0, 0; j < n; i, j = i+1, j+2 {
		in1[i] = *elements[j].impl().data
		in2[i] = *elements[j+1].impl().data
	}
	return el.doProdPair(in1, in2)
}

func (el *elementImpl) ProdPairSlice(x, y []Element) Element {
	n := len(x)
	if n != len(y) {
		panic(ErrBadPairList)
	}
	in1 := make([]C.struct_element_s, n)
	in2 := make([]C.struct_element_s, n)
	for i := 0; i < n; i++ {
		in1[i] = *x[i].impl().data
		in2[i] = *y[i].impl().data
	}
	return el.doProdPair(in1, in2)
}

func (el *elementImpl) PreparePairer() Pairer { return initPairer(el) }

func (el *elementImpl) PairerPair(pairer Pairer, x Element) Element {
	C.pairing_pp_apply(el.data, x.impl().data, pairer.(*pairerImpl).data)
	return el
}

func (el *elementImpl) BruteForceDL(g, h Element) Element {
	C.element_dlog_brute_force(el.data, g.impl().data, h.impl().data)
	return el
}

func (el *elementImpl) PollardRhoDL(g, h Element) Element {
	C.element_dlog_pollard_rho(el.data, g.impl().data, h.impl().data)
	return el
}

func (el *elementImpl) Rand() Element {
	C.element_random(el.data)
	return el
}

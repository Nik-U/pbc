package pbc

/*
#include <pbc/pbc.h>
*/
import "C"

import (
	"hash"
	"math/big"
	"runtime"
	"unsafe"
)

type uncheckedElement struct {
	pairing *pairingImpl
	d       *C.struct_element_s
}

func clearUnchecked(element *uncheckedElement) {
	C.element_clear(element.d)
}

func initUnchecked(element *uncheckedElement, pairing *pairingImpl, initialize bool, field Field) {
	element.d = &C.struct_element_s{}
	element.pairing = pairing
	if initialize {
		switch field {
		case G1:
			C.element_init_G1(element.d, pairing.data)
		case G2:
			C.element_init_G2(element.d, pairing.data)
		case GT:
			C.element_init_GT(element.d, pairing.data)
		case Zr:
			C.element_init_Zr(element.d, pairing.data)
		default:
			panic(ErrUnknownField)
		}
	}
	runtime.SetFinalizer(element, clearUnchecked)
}

func makeUnchecked(pairing *pairingImpl, field Field) *uncheckedElement {
	element := &uncheckedElement{}
	initUnchecked(element, pairing, true, field)
	return element
}

func (el *uncheckedElement) data() *C.struct_element_s { return el.d }

func (el *uncheckedElement) Pairing() Pairing { return el.pairing }

func (el *uncheckedElement) NewFieldElement() Element {
	newElement := &uncheckedElement{}
	initUnchecked(newElement, el.pairing, false, G1)
	C.element_init_same_as(newElement.d, el.d)
	return newElement
}

func (el *uncheckedElement) Set0() Element {
	C.element_set0(el.d)
	return el
}

func (el *uncheckedElement) Set1() Element {
	C.element_set1(el.d)
	return el
}

func (el *uncheckedElement) SetInt32(i int32) Element {
	C.element_set_si(el.d, C.long(i))
	return el
}

func (el *uncheckedElement) SetBig(i *big.Int) Element {
	C.element_set_mpz(el.d, &big2mpz(i)[0])
	return el
}

func (el *uncheckedElement) Set(src Element) Element {
	C.element_set(el.d, src.data())
	return el
}

func (el *uncheckedElement) SetFromHash(hash []byte) Element {
	C.element_from_hash(el.d, unsafe.Pointer(&hash[0]), C.int(len(hash)))
	return el
}

func (el *uncheckedElement) SetFromStringHash(s string, h hash.Hash) Element {
	h.Reset()
	if _, err := h.Write([]byte(s)); err != nil {
		panic(ErrHashFailure)
	}
	return el.SetFromHash(h.Sum([]byte{}))
}

func (el *uncheckedElement) SetBytes(buf []byte) Element {
	C.element_from_bytes(el.d, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

func (el *uncheckedElement) SetXBytes(buf []byte) Element {
	C.element_from_bytes_x_only(el.d, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

func (el *uncheckedElement) SetCompressedBytes(buf []byte) Element {
	C.element_from_bytes_compressed(el.d, (*C.uchar)(unsafe.Pointer(&buf[0])))
	return el
}

func (el *uncheckedElement) SetString(s string, base int) (Element, bool) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))

	if ok := C.element_set_str(el.d, cstr, C.int(base)); ok == 0 {
		return nil, false
	}
	return el, true
}

func (el *uncheckedElement) BigInt() *big.Int {
	mpz := newMpz()
	C.element_to_mpz(&mpz[0], el.d)
	return mpz2big(mpz)
}

func (el *uncheckedElement) BytesLen() int {
	return int(C.element_length_in_bytes(el.d))
}

func (el *uncheckedElement) writeBytes(buf []byte) C.int {
	return C.element_to_bytes((*C.uchar)(unsafe.Pointer(&buf[0])), el.d)
}

func (el *uncheckedElement) Bytes() []byte {
	buf := make([]byte, el.BytesLen())
	el.writeBytes(buf)
	return buf
}

func (el *uncheckedElement) XBytesLen() int {
	return int(C.element_length_in_bytes_x_only(el.d))
}

func (el *uncheckedElement) writeXBytes(buf []byte) C.int {
	return C.element_to_bytes_x_only((*C.uchar)(unsafe.Pointer(&buf[0])), el.d)
}

func (el *uncheckedElement) XBytes() []byte {
	buf := make([]byte, el.XBytesLen())
	el.writeXBytes(buf)
	return buf
}

func (el *uncheckedElement) CompressedBytesLen() int {
	return int(C.element_length_in_bytes_compressed(el.d))
}

func (el *uncheckedElement) writeCompressedBytes(buf []byte) C.int {
	return C.element_to_bytes_compressed((*C.uchar)(unsafe.Pointer(&buf[0])), el.d)
}

func (el *uncheckedElement) CompressedBytes() []byte {
	buf := make([]byte, el.CompressedBytesLen())
	el.writeCompressedBytes(buf)
	return buf
}

func (el *uncheckedElement) Len() int {
	return int(C.element_item_count(el.d))
}

func (el *uncheckedElement) Item(i int) Element {
	return &uncheckedElement{
		pairing: el.pairing,
		d:       C.element_item(el.d, C.int(i)),
	}
}

func (el *uncheckedElement) X() *big.Int {
	return el.Item(0).BigInt()
}

func (el *uncheckedElement) Y() *big.Int {
	return el.Item(1).BigInt()
}

func (el *uncheckedElement) Is0() bool {
	return C.element_is0(el.d) != 0
}

func (el *uncheckedElement) Is1() bool {
	return C.element_is1(el.d) != 0
}

func (el *uncheckedElement) IsSquare() bool {
	return C.element_is_sqr(el.d) != 0
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

func (el *uncheckedElement) Sign() int {
	return normalizeSign(int64(C.element_sign(el.d)))
}

func (el *uncheckedElement) Cmp(x Element) int {
	return normalizeSign(int64(C.element_cmp(el.d, x.data())))
}

func (el *uncheckedElement) Equals(x Element) bool { return el.Cmp(x) == 0 }

func (el *uncheckedElement) Add(x, y Element) Element {
	C.element_add(el.d, x.data(), y.data())
	return el
}

func (el *uncheckedElement) Sub(x, y Element) Element {
	C.element_sub(el.d, x.data(), y.data())
	return el
}

func (el *uncheckedElement) Mul(x, y Element) Element {
	C.element_mul(el.d, x.data(), y.data())
	return el
}

func (el *uncheckedElement) MulBig(x Element, i *big.Int) Element {
	C.element_mul_mpz(el.d, x.data(), &big2mpz(i)[0])
	return el
}

func (el *uncheckedElement) MulInt32(x Element, i int32) Element {
	C.element_mul_si(el.d, x.data(), C.long(i))
	return el
}

func (el *uncheckedElement) MulZn(x, y Element) Element {
	C.element_mul_zn(el.d, x.data(), y.data())
	return el
}

func (el *uncheckedElement) Div(x, y Element) Element {
	C.element_div(el.d, x.data(), y.data())
	return el
}

func (el *uncheckedElement) Double(x Element) Element {
	C.element_double(el.d, x.data())
	return el
}

func (el *uncheckedElement) Halve(x Element) Element {
	C.element_halve(el.d, x.data())
	return el
}

func (el *uncheckedElement) Square(x Element) Element {
	C.element_square(el.d, x.data())
	return el
}

func (el *uncheckedElement) Neg(x Element) Element {
	C.element_neg(el.d, x.data())
	return el
}

func (el *uncheckedElement) Invert(x Element) Element {
	C.element_invert(el.d, x.data())
	return el
}

func (el *uncheckedElement) PowBig(x Element, i *big.Int) Element {
	C.element_pow_mpz(el.d, x.data(), &big2mpz(i)[0])
	return el
}

func (el *uncheckedElement) PowZn(x, i Element) Element {
	C.element_pow_zn(el.d, x.data(), i.data())
	return el
}

func (el *uncheckedElement) Pow2Big(x Element, i *big.Int, y Element, j *big.Int) Element {
	C.element_pow2_mpz(el.d, x.data(), &big2mpz(i)[0], y.data(), &big2mpz(j)[0])
	return el
}

func (el *uncheckedElement) Pow2Zn(x, i, y, j Element) Element {
	C.element_pow2_zn(el.d, x.data(), i.data(), y.data(), j.data())
	return el
}

func (el *uncheckedElement) Pow3Big(x Element, i *big.Int, y Element, j *big.Int, z Element, k *big.Int) Element {
	C.element_pow3_mpz(el.d, x.data(), &big2mpz(i)[0], y.data(), &big2mpz(j)[0], z.data(), &big2mpz(k)[0])
	return el
}

func (el *uncheckedElement) Pow3Zn(x, i, y, j, z, k Element) Element {
	C.element_pow3_zn(el.d, x.data(), i.data(), y.data(), j.data(), z.data(), k.data())
	return el
}

func (el *uncheckedElement) PreparePower() Power { return initPower(el) }

func (el *uncheckedElement) PowerBig(power Power, i *big.Int) Element {
	C.element_pp_pow(el.d, &big2mpz(i)[0], power.(*powerImpl).data)
	return el
}

func (el *uncheckedElement) PowerZn(power Power, i Element) Element {
	C.element_pp_pow_zn(el.d, i.data(), power.(*powerImpl).data)
	return el
}

func (el *uncheckedElement) Pair(x, y Element) Element {
	C.pairing_apply(el.d, x.data(), y.data(), el.pairing.data)
	return el
}

func (el *uncheckedElement) doProdPair(in1, in2 []C.struct_element_s) Element {
	x := (*C.element_t)(unsafe.Pointer(&in1[0]))
	y := (*C.element_t)(unsafe.Pointer(&in2[0]))
	C.element_prod_pairing(el.d, x, y, C.int(len(in1)))
	return el
}

func (el *uncheckedElement) ProdPair(elements ...Element) Element {
	n := len(elements)
	if n%2 != 0 {
		panic(ErrBadPairList)
	}
	half := n / 2
	in1 := make([]C.struct_element_s, half)
	in2 := make([]C.struct_element_s, half)
	for i, j := 0, 0; j < n; i, j = i+1, j+2 {
		in1[i] = *elements[j].data()
		in2[i] = *elements[j+1].data()
	}
	return el.doProdPair(in1, in2)
}

func (el *uncheckedElement) ProdPairSlice(x, y []Element) Element {
	n := len(x)
	if n != len(y) {
		panic(ErrBadPairList)
	}
	in1 := make([]C.struct_element_s, n)
	in2 := make([]C.struct_element_s, n)
	for i := 0; i < n; i++ {
		in1[i] = *x[i].data()
		in2[i] = *y[i].data()
	}
	return el.doProdPair(in1, in2)
}

func (el *uncheckedElement) PreparePairer() Pairer { return initPairer(el) }

func (el *uncheckedElement) PairerPair(pairer Pairer, y Element) Element {
	C.pairing_pp_apply(el.d, y.data(), pairer.(*pairerImpl).data)
	return el
}

func (el *uncheckedElement) BruteForceDL(g, h Element) Element {
	C.element_dlog_brute_force(el.d, g.data(), h.data())
	return el
}

func (el *uncheckedElement) PollardRhoDL(g, h Element) Element {
	C.element_dlog_pollard_rho(el.d, g.data(), h.data())
	return el
}

func (el *uncheckedElement) Rand() Element {
	C.element_random(el.d)
	return el
}

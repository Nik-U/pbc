package pbc

/*
#include <pbc/pbc.h>
*/
import "C"
import (
	"math/big"
	"runtime"
	"unsafe"
)

func (el *Element) Set0() *Element {
	C.element_set0(el.cptr)
	return el
}

func (el *Element) Set1() *Element {
	C.element_set1(el.cptr)
	return el
}

func (el *Element) Rand() *Element {
	C.element_random(el.cptr)
	return el
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

func (el *Element) Cmp(x *Element) int {
	if el.checked {
		el.checkCompatible(x)
	}
	return normalizeSign(int64(C.element_cmp(el.cptr, x.cptr)))
}

func (el *Element) Equals(x *Element) bool { return el.Cmp(x) == 0 }

func (el *Element) Is0() bool {
	return C.element_is0(el.cptr) != 0
}

func (el *Element) Is1() bool {
	return C.element_is1(el.cptr) != 0
}

func (el *Element) IsSquare() bool {
	return C.element_is_sqr(el.cptr) != 0
}

func (el *Element) Sign() int {
	return normalizeSign(int64(C.element_sign(el.cptr)))
}

func (el *Element) Add(x, y *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
	}
	C.element_add(el.cptr, x.cptr, y.cptr)
	return el
}

func (el *Element) Sub(x, y *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
	}
	C.element_sub(el.cptr, x.cptr, y.cptr)
	return el
}

func (el *Element) Mul(x, y *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
	}
	C.element_mul(el.cptr, x.cptr, y.cptr)
	return el
}

func (el *Element) MulBig(x *Element, i *big.Int) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_mul_mpz(el.cptr, x.cptr, &big2mpz(i)[0])
	return el
}

func (el *Element) MulInt32(x *Element, i int32) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_mul_si(el.cptr, x.cptr, C.long(i))
	return el
}

func (el *Element) MulZn(x, y *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
		y.checkInteger()
	}
	C.element_mul_zn(el.cptr, x.cptr, y.cptr)
	return el
}

func (el *Element) Div(x, y *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
	}
	C.element_div(el.cptr, x.cptr, y.cptr)
	return el
}

func (el *Element) Double(x *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_double(el.cptr, x.cptr)
	return el
}

func (el *Element) Halve(x *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_halve(el.cptr, x.cptr)
	return el
}

func (el *Element) Square(x *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_square(el.cptr, x.cptr)
	return el
}

func (el *Element) Neg(x *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_neg(el.cptr, x.cptr)
	return el
}

func (el *Element) Invert(x *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_invert(el.cptr, x.cptr)
	return el
}

func (el *Element) PowBig(x *Element, i *big.Int) *Element {
	if el.checked {
		el.checkCompatible(x)
	}
	C.element_pow_mpz(el.cptr, x.cptr, &big2mpz(i)[0])
	return el
}

func (el *Element) PowZn(x, i *Element) *Element {
	if el.checked {
		el.checkCompatible(x)
		i.checkInteger()
	}
	C.element_pow_zn(el.cptr, x.cptr, i.cptr)
	return el
}

func (el *Element) Pow2Big(x *Element, i *big.Int, y *Element, j *big.Int) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
	}
	C.element_pow2_mpz(el.cptr, x.cptr, &big2mpz(i)[0], y.cptr, &big2mpz(j)[0])
	return el
}

func (el *Element) Pow2Zn(x, i, y, j *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y)
		i.checkInteger()
		j.checkInteger()
	}
	C.element_pow2_zn(el.cptr, x.cptr, i.cptr, y.cptr, j.cptr)
	return el
}

func (el *Element) Pow3Big(x *Element, i *big.Int, y *Element, j *big.Int, z *Element, k *big.Int) *Element {
	if el.checked {
		el.checkAllCompatible(x, y, z)
	}
	C.element_pow3_mpz(el.cptr, x.cptr, &big2mpz(i)[0], y.cptr, &big2mpz(j)[0], z.cptr, &big2mpz(k)[0])
	return el
}

func (el *Element) Pow3Zn(x, i, y, j, z, k *Element) *Element {
	if el.checked {
		el.checkAllCompatible(x, y, z)
		i.checkInteger()
		j.checkInteger()
		k.checkInteger()
	}
	C.element_pow3_zn(el.cptr, x.cptr, i.cptr, y.cptr, j.cptr, z.cptr, k.cptr)
	return el
}

// Power stores pre-processed information to quickly exponentiate an element.
// A Power can be generated for Element x by calling x.PreparePower(). When
// PowBig or PowZn is called with Element target and integer i, the result of
// x^i will be stored in target.
type Power struct {
	source *Element // Prevents garbage collection
	pp     *C.struct_element_pp_s
}

func (power *Power) Source() *Element { return power.source }

func (power *Power) PowBig(target *Element, i *big.Int) *Element {
	return target.PowerBig(power, i)
}

func (power *Power) PowZn(target *Element, i *Element) *Element {
	return target.PowerZn(power, i)
}

func clearPower(power *Power) {
	C.element_pp_clear(power.pp)
}

func (el *Element) PreparePower() *Power {
	power := &Power{
		source: el,
		pp:     &C.struct_element_pp_s{},
	}
	C.element_pp_init(power.pp, el.cptr)
	runtime.SetFinalizer(power, clearPower)
	return power
}

func (el *Element) PowerBig(power *Power, i *big.Int) *Element {
	C.element_pp_pow(el.cptr, &big2mpz(i)[0], power.pp)
	return el
}

func (el *Element) PowerZn(power *Power, i *Element) *Element {
	if el.checked {
		i.checkInteger()
	}
	C.element_pp_pow_zn(el.cptr, i.cptr, power.pp)
	return el
}

func (el *Element) Pair(x, y *Element) *Element {
	if el.checked {
		x.ensureChecked()
		y.ensureChecked()
		pairing := el.pairing.cptr
		checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
		checkFieldsMatch(x.fieldPtr, pairing.G1)
		checkFieldsMatch(y.fieldPtr, pairing.G2)
	}
	C.pairing_apply(el.cptr, x.cptr, y.cptr, el.pairing.cptr)
	return el
}

func (el *Element) doProdPair(in1, in2 []C.struct_element_s) *Element {
	x := (*C.element_t)(unsafe.Pointer(&in1[0]))
	y := (*C.element_t)(unsafe.Pointer(&in2[0]))
	C.element_prod_pairing(el.cptr, x, y, C.int(len(in1)))
	return el
}

func (el *Element) ProdPair(elements ...*Element) *Element {
	n := len(elements)
	if n%2 != 0 {
		panic(ErrBadPairList)
	}
	if el.checked {
		pairing := el.pairing.cptr
		checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
		for i := 1; i < n; i += 2 {
			elements[i-1].ensureChecked()
			elements[i].ensureChecked()
			checkFieldsMatch(elements[i-1].fieldPtr, pairing.G1)
			checkFieldsMatch(elements[i].fieldPtr, pairing.G2)
		}
	}
	half := n / 2
	in1 := make([]C.struct_element_s, half)
	in2 := make([]C.struct_element_s, half)
	for i, j := 0, 0; j < n; i, j = i+1, j+2 {
		in1[i] = *elements[j].cptr
		in2[i] = *elements[j+1].cptr
	}
	return el.doProdPair(in1, in2)
}

func (el *Element) ProdPairSlice(x, y []*Element) *Element {
	n := len(x)
	if n != len(y) {
		panic(ErrBadPairList)
	}
	if el.checked {
		pairing := el.pairing.cptr
		checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
		for i := 1; i < n; i++ {
			x[i].ensureChecked()
			checkFieldsMatch(x[i].fieldPtr, pairing.G1)
		}
		n = len(y)
		for i := 1; i < n; i++ {
			y[i].ensureChecked()
			checkFieldsMatch(y[i].fieldPtr, pairing.G2)
		}
	}
	in1 := make([]C.struct_element_s, n)
	in2 := make([]C.struct_element_s, n)
	for i := 0; i < n; i++ {
		in1[i] = *x[i].cptr
		in2[i] = *y[i].cptr
	}
	return el.doProdPair(in1, in2)
}

// Pairer stores pre-processed information to quickly pair an element. A Pairer
// can be generated for Element x by calling x.PreparePairer(). When Pair is
// called with Elements target and y, the result of e(x,y) will be stored in
// target.
type Pairer struct {
	source *Element // Prevents garbage collection
	pp     *C.struct_pairing_pp_s
}

func (pairer *Pairer) Source() *Element { return pairer.source }

func (pairer *Pairer) Pair(target *Element, y *Element) *Element {
	return target.PairerPair(pairer, y)
}

func clearPairer(pairer *Pairer) {
	C.pairing_pp_clear(pairer.pp)
}

func (el *Element) PreparePairer() *Pairer {
	pairer := &Pairer{
		source: el,
		pp:     &C.struct_pairing_pp_s{},
	}
	C.pairing_pp_init(pairer.pp, el.cptr, el.pairing.cptr)
	runtime.SetFinalizer(pairer, clearPairer)
	return pairer
}

func (el *Element) PairerPair(pairer *Pairer, y *Element) *Element {
	if el.checked {
		pairer.source.ensureChecked()
		y.ensureChecked()
		pairing := el.pairing.cptr
		checkFieldsMatch(pairer.source.fieldPtr, pairing.G1)
		checkFieldsMatch(y.fieldPtr, pairing.G2)
		checkFieldsMatch(el.fieldPtr, &pairing.GT[0])
	}
	C.pairing_pp_apply(el.cptr, y.cptr, pairer.pp)
	return el
}

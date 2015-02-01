package pbc

/*
#include <stdint.h>
#include <pbc/pbc.h>

typedef struct {
	int typeD;
	pbc_param_ptr params;
	uint32_t rbits;
	uint32_t qbits;
} check_pairing_settings_t;

int checkPairing(pbc_cm_t cm, void* p) {
	check_pairing_settings_t* settings = (check_pairing_settings_t*)p;

	unsigned int rbits = (unsigned int)mpz_sizeinbase(cm->r, 2);
	unsigned int qbits = (unsigned int)mpz_sizeinbase(cm->q, 2);
	if (rbits < settings->rbits || qbits < settings->qbits) return 0;

	if (settings->typeD) {
		pbc_param_init_d_gen(settings->params, cm);
	} else {
		pbc_param_init_g_gen(settings->params, cm);
	}
	return 1;
}
*/
import "C"

import (
	"math/big"
	"unsafe"
)

// GenerateA generates a pairing on the curve y^2 = x^3 + x over the field F_q
// for some prime q = 3 mod 4. Type A pairings are symmetric (i.e., G1 == G2).
// Type A pairings are best used when speed of computation is the primary
// concern.
//
// To be secure, generic discrete log algorithms must be infeasible in groups of
// order r, and finite field discrete log algorithms must be infeasible in
// groups of order q^2.
//
// For example:
// 	params := pbc.GenerateA(160, 512)
func GenerateA(rbits uint32, qbits uint32) Params {
	params := makeParams()
	C.pbc_param_init_a_gen(params.data, C.int(rbits), C.int(qbits))
	return params
}

// GenerateA1 generates a type A pairing given a fixed order for G1, G2, and GT.
// This form of pairing can be used to produce groups of composite order, where
// r is the product of two large primes. In this case, r should infeasible to
// factor. Each prime should be at least 512 bits (causing r to be 1024 bits in
// general), but preferably 1024 bits or more.
func GenerateA1(r *big.Int) Params {
	params := makeParams()
	C.pbc_param_init_a1_gen(params.data, &big2mpz(r)[0])
	return params
}

// GenerateD generates a pairing on a curve with embedding degree 6 whose order
// is h * r where r is prime and h is a small constant. Type D pairings are
// asymmetric, but have small group elements. This makes them well-suited for
// applications where message size is the primary concern, but speed is also
// important.
//
// Parameters are generated using the constant multiplication (CM) method for a
// given fundamental discriminant D. It is required that D > 0, no square of an
// odd prime divides D, and D = 0 or 3 mod 4. The bitlimit parameter sets a cap
// on the number of bits in the group order. It is possible that for some values
// of D, no suitable curves can be found. In this case, GenerateD returns nil
// and ErrNoSuitableCurves.
//
// The rbits and qbits parameters sit minimum sizes for group orders. To be
// secure, generic discrete log algorithms must be infeasible in groups of order
// r, and finite field discrete log algorithms must be infeasible in groups of
// order q^6.
//
// For example:
// 	params, err := pbc.GenerateD(9563, 160, 171, 500)
func GenerateD(d uint32, rbits uint32, qbits uint32, bitlimit uint32) (Params, error) {
	return generateWithCM(true, d, rbits, qbits, bitlimit)
}

// GenerateE generates a pairing entirely within a order r subgroup of an order
// q field. These pairings are symmetric, but serve little purpose beyond being
// mathematically interesting. Use of these pairings is not recommended unless
// new algorithms are discovered for solving discrete logs in elliptic curves as
// easily as for finite fields.
//
// For security, generic discrete log algorithms must be infeasible in groups of
// order r, and finite field discrete log algorithms must be infeasible in
// finite fields of order q.
//
// For example:
// 	params, err := pbc.GenerateE(160, 1024)
func GenerateE(rbits uint32, qbits uint32) Params {
	params := makeParams()
	C.pbc_param_init_e_gen(params.data, C.int(rbits), C.int(qbits))
	return params
}

// GenerateF generates an asymmetric pairing with extremely small group
// elements. This is the best pairing to use when space is an overriding
// priority. However, type F pairings are slow compared to the other types. Type
// D pairings provide a more balanced alternative.
//
// The bits parameter specifies the approximate number of bits in the group
// order, r, and the order of the base field, q. For security, generic discrete
// log algorithms must be infeasible in groups of order r, and finite field
// discrete log algorithms must be infeasible in finite fields of order q^12.
//
// For example:
// 	params, err := pbc.GenerateF(160)
func GenerateF(bits uint32) Params {
	params := makeParams()
	C.pbc_param_init_f_gen(params.data, C.int(bits))
	return params
}

// GenerateG generates a pairing on a curve with embedding degree 10 whose order
// is h * r where r is prime and h is a small constant. Type G pairings are
// asymmetric, but have extremely small group elements. However, these pairings
// are even slower than type F pairings, making type F a better choice.
//
// Like type D pairings, parameters are generated using the constant
// multiplication (CM) method. See the GenerateD function for a description of
// the parameters.
//
// For example:
// 	params, err := pbc.GenerateG(9563, 160, 171, 500)
func GenerateG(d uint32, rbits uint32, qbits uint32, bitlimit uint32) (Params, error) {
	return generateWithCM(false, d, qbits, rbits, bitlimit)
}

func generateWithCM(typeD bool, d uint32, rbits uint32, qbits uint32, bitlimit uint32) (Params, error) {
	params := makeParams()
	settings := &C.check_pairing_settings_t{
		params: params.data,
		rbits:  C.uint32_t(rbits),
		qbits:  C.uint32_t(qbits),
	}
	if typeD {
		settings.typeD = C.int(1)
	}
	res := C.pbc_cm_search_d((*[0]byte)(C.checkPairing), unsafe.Pointer(settings), C.uint(d), C.uint(bitlimit))
	if res != 1 {
		return nil, ErrNoSuitableCurves
	}
	return params, nil
}

package pbc

/*
#include <pbc/pbc.h>

int acceptPairingD(pbc_cm_t cm, void* p) {
	pbc_param_init_d_gen((pbc_param_ptr)p, cm);
	return 1;
}

int acceptPairingG(pbc_cm_t cm, void* p) {
	pbc_param_init_g_gen((pbc_param_ptr)p, cm);
	return 1;
}
*/
import "C"

import (
	"math/big"
	"unsafe"
)

func GenerateA(rbits uint32, qbits uint32) Params {
	params := makeParams()
	C.pbc_param_init_a_gen(params, C.int(rbits), C.int(qbits))
	return params
}

func GenerateA1(n *big.Int) Params {
	params := makeParams()
	C.pbc_param_init_a1_gen(params, &big2mpz(n)[0])
	return params
}

func GenerateD(d uint32, bitlimit uint32) Params {
	params := makeParams()
	C.pbc_cm_search_d((*[0]byte)(C.acceptPairingD), unsafe.Pointer(params), C.uint(d), C.uint(bitlimit))
	return params
}

func GenerateE(rbits uint32, qbits uint32) Params {
	params := makeParams()
	C.pbc_param_init_e_gen(params, C.int(rbits), C.int(qbits))
	return params
}

func GenerateF(bits uint32) Params {
	params := makeParams()
	C.pbc_param_init_f_gen(params, C.int(bits))
	return params
}

func GenerateG(d uint32, bitlimit uint32) Params {
	params := makeParams()
	C.pbc_cm_search_d((*[0]byte)(C.acceptPairingG), unsafe.Pointer(params), C.uint(d), C.uint(bitlimit))
	return params
}

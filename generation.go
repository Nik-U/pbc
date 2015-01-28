package pbc

/*
#cgo LDFLAGS: /usr/local/lib/libpbc.a -lgmp
#include <pbc/pbc.h>

int acceptPairingD(pbc_cm_t cm, void* p) {
	pbc_param_init_d_gen((pbc_param_ptr)p, cm);
	return 1;
}

int acceptPairingG(pbc_cm_t cm, void* p) {
	pbc_param_init_g_gen((pbc_param_ptr)p, cm);
	return 1;
}

void genPairingD(pbc_param_ptr p, unsigned int D, unsigned int bitlimit) {
	pbc_cm_search_d(acceptPairingD, p, D, bitlimit);
}

void genPairingG(pbc_param_ptr p, unsigned int D, unsigned int bitlimit) {
	pbc_cm_search_g(acceptPairingG, p, D, bitlimit);
}
*/
import "C"

import "math/big"

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
	C.genPairingD(params, C.uint(d), C.uint(bitlimit))
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
	C.genPairingG(params, C.uint(d), C.uint(bitlimit))
	return params
}

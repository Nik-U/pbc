package pbc

/*
#cgo LDFLAGS: /usr/local/lib/libpbc.a -lgmp
#include <pbc/pbc.h>
*/
import "C"

import "io"
import "runtime"

type Pairing interface{}

type pairingImpl struct {
	data C.pairing_ptr
}

func pairingFinalize(p *pairingImpl) {

}

func NewPairing(params io.Reader) Pairing {
	x := &pairingImpl{}
	runtime.SetFinalizer(x, pairingFinalize)
	return x
}

func NewPairingFromString(params string) Pairing {
	return nil
}

func NewPairingFromParams(params Params) Pairing {
	return nil
}

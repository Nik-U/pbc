package pbc

/*
#cgo LDFLAGS: /usr/local/lib/libpbc.a -lgmp
#include <pbc/pbc.h>

int param_out_str_wrapper(char** bufp, size_t* sizep, pbc_param_t p) {
	FILE* handle = open_memstream(bufp, sizep);
	if (!handle) return 0;
	pbc_param_out_str(handle, p);
	fclose(handle);
	return 1;
}
*/
import "C"

import (
	"errors"
	"io"
	"runtime"
	"unsafe"
)

var ErrInvalidParamString = errors.New("invalid pairing parameters")

type Params interface {
	NewPairing() Pairing
	WriteTo(w io.Writer) (n int64, err error)
	String() string
}

func NewParamsFromString(s string) (Params, error) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))

	params := makeParams()
	if ok := C.pbc_param_init_set_str(params, cstr); ok != 0 {
		return nil, ErrInvalidParamString
	}
	return params, nil
}

func (params *C.struct_pbc_param_s) NewPairing() Pairing {
	return NewPairingFromParams(params)
}

func (params *C.struct_pbc_param_s) WriteTo(w io.Writer) (n int64, err error) {
	count, err := io.WriteString(w, params.String())
	return int64(count), err
}

func (params *C.struct_pbc_param_s) String() string {
	var buf *C.char
	var bufLen C.size_t
	if C.param_out_str_wrapper(&buf, &bufLen, params) == 0 {
		return ""
	}
	str := C.GoStringN(buf, C.int(bufLen))
	C.free(unsafe.Pointer(buf))
	return str
}

func clearParams(params *C.struct_pbc_param_s) {
	println("clearparams")
	C.pbc_param_clear(params)
}

func makeParams() *C.struct_pbc_param_s {
	params := &C.struct_pbc_param_s{}
	runtime.SetFinalizer(params, clearParams)
	return params
}

package pbc

/*
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
	"io"
	"io/ioutil"
	"runtime"
	"unsafe"
)

// Params represents the parameters required for creating a pairing. Parameters
// can be generated using the generation functions or read from a Reader.
// Normally, parameters are generated once using a generation program and then
// distributed with the final application. Parameters can be exported for this
// purpose using the WriteTo or String methods.
//
// For applications requiring fast computation, type A pairings are preferred.
// Applications requiring small message sizes should consider type D pairings.
// If speed is not a concern, type F pairings yield the smallest messages at
// the cost of additional computation. Applications requiring symmetric
// pairings should use type A. If a specific group order must be used (e.g.,
// for composite orders), then type A1 pairings are required.
type Params struct {
	cptr *C.struct_pbc_param_s
}

// NewParams loads pairing parameters from a Reader.
func NewParams(r io.Reader) (*Params, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewParamsFromString(string(b))
}

// NewParamsFromString loads pairing parameters from a string.
func NewParamsFromString(s string) (*Params, error) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))

	params := makeParams()
	if ok := C.pbc_param_init_set_str(params.cptr, cstr); ok != 0 {
		return nil, ErrInvalidParamString
	}
	return params, nil
}

// NewPairing creates a Pairing using these parameters.
func (params *Params) NewPairing() *Pairing {
	return NewPairing(params)
}

// WriteTo writes the pairing parameters to a Writer.
func (params *Params) WriteTo(w io.Writer) (n int64, err error) {
	count, err := io.WriteString(w, params.String())
	return int64(count), err
}

// String returns a string representation of the pairing parameters.
func (params *Params) String() string {
	var buf *C.char
	var bufLen C.size_t
	if C.param_out_str_wrapper(&buf, &bufLen, params.cptr) == 0 {
		return ""
	}
	str := C.GoStringN(buf, C.int(bufLen))
	C.free(unsafe.Pointer(buf))
	return str
}

func clearParams(params *Params) {
	C.pbc_param_clear(params.cptr)
}

func makeParams() *Params {
	params := &Params{cptr: &C.struct_pbc_param_s{}}
	runtime.SetFinalizer(params, clearParams)
	return params
}

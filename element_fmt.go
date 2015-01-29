package pbc

/*
#include <pbc/pbc.h>

int element_out_str_wrapper(char** bufp, size_t* sizep, int base, element_t e) {
	FILE* handle = open_memstream(bufp, sizep);
	if (!handle) return 0;
	element_out_str(handle, base, e);
	fclose(handle);
	return 1;
}
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

var (
	ErrBadInput = errors.New("invalid element format during scan")
	ErrBadVerb  = errors.New("invalid verb specified for scan")
)

func (el *elementImpl) errorFormat(f fmt.State, c rune, err string) {
	fmt.Fprintf(f, "%%!%c(%s pbc.Element)", c, err)
}

func (el *elementImpl) nativeFormat(f fmt.State, c rune) {
	base := 10
	if width, ok := f.Width(); ok {
		if width < 2 || width > 36 {
			el.errorFormat(f, c, "BADBASE")
			return
		}
		base = width
	}
	var buf *C.char
	var bufLen C.size_t
	if C.element_out_str_wrapper(&buf, &bufLen, C.int(base), el.data) == 0 {
		el.errorFormat(f, c, "INTERNALERROR")
		return
	}
	str := C.GoStringN(buf, C.int(bufLen))
	C.free(unsafe.Pointer(buf))
	fmt.Fprintf(f, "%s", str)
}

func (el *elementImpl) customFormat(f fmt.State, c rune) {
	count := el.Len()
	if count == 0 {
		el.BigInt().Format(f, c)
	} else {
		fmt.Fprintf(f, "[")
		for i := 0; i < count; i++ {
			el.Item(i).impl().customFormat(f, c)
			if i+1 < count {
				fmt.Fprintf(f, ", ")
			}
		}
		fmt.Fprintf(f, "]")
	}
}

func (el *elementImpl) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('#') {
			fmt.Fprintf(f, "pbc.Element{Checked: false, Pairing: %p, Addr: %p}", el.pairing, el)
			break
		}
		fallthrough
	case 's':
		el.nativeFormat(f, c)
	case 'd', 'b', 'o', 'x', 'X':
		el.customFormat(f, c)
	default:
		el.errorFormat(f, c, "BADVERB")
	}
}

func (el *checkedElement) Format(f fmt.State, c rune) {
	if c == 'v' && f.Flag('#') {
		fmt.Fprintf(f, "pbc.Element{Checked: true, Integer: %t, Field: %p, Pairing: %p, Addr: %p}", el.isInteger, el.fieldPtr, el.unchecked.pairing, el)
	} else {
		el.unchecked.Format(f, c)
	}
}

func (el *elementImpl) String() string {
	return fmt.Sprintf("%s", el)
}

func (el *checkedElement) String() string { return el.unchecked.String() }

func (el *elementImpl) Scan(state fmt.ScanState, verb rune) error {
	if verb != 's' && verb != 'v' {
		return ErrBadVerb
	}
	base, ok := state.Width()
	if !ok {
		base = 10
	} else if base < 2 || base > 36 {
		return ErrBadVerb
	}
	maxDigit := '9'
	maxAlpha := 'z'
	if base < 10 {
		maxDigit = rune('0' + (base - 1))
	}
	if base < 36 {
		maxAlpha = rune('a' + (base - 11))
	}

	state.SkipSpace()

	tokensFound := make([]uint, 0, 5)
	inToken := false
	justDescended := false
	expectTokenDone := false
	var buf bytes.Buffer

ReadLoop:
	for {
		r, _, err := state.ReadRune()
		if err != nil {
			if err == io.EOF {
				if len(tokensFound) == 0 {
					break ReadLoop
				}
				return ErrBadInput
			}
			return err
		}
		buf.WriteRune(r)

		if expectTokenDone && r != ',' && r != ']' {
			return ErrBadInput
		}
		expectTokenDone = false

		switch r {
		case '[':
			if inToken {
				return ErrBadInput
			}
			tokensFound = append(tokensFound, 0)
		case ']':
			if !inToken || len(tokensFound) == 0 || tokensFound[len(tokensFound)-1] == 0 {
				return ErrBadInput
			}
			tokensFound = tokensFound[:len(tokensFound)-1]
			if len(tokensFound) == 0 {
				break ReadLoop
			}
		case ',':
			if len(tokensFound) == 0 || (!inToken && !justDescended) {
				return ErrBadInput
			}
			tokensFound[len(tokensFound)-1]++
			inToken = false
			state.SkipSpace()
		case 'O':
			if inToken {
				return ErrBadInput
			}
			expectTokenDone = true
			inToken = true
		default:
			if (r < '0' || r > maxDigit) && (r < 'a' || r > maxAlpha) {
				return ErrBadInput
			}
			inToken = true
		}
		justDescended = (r == ']')
	}
	if _, ok := el.SetString(buf.String(), base); !ok {
		return ErrBadInput
	}
	return nil
}

func (el *checkedElement) Scan(state fmt.ScanState, verb rune) error {
	return el.unchecked.Scan(state, verb)
}

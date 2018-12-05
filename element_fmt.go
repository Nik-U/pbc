// Copyright © 2018 Nik Unger
//
// This file is part of The PBC Go Wrapper.
//
// The PBC Go Wrapper is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// The PBC Go Wrapper is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public
// License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with The PBC Go Wrapper. If not, see <http://www.gnu.org/licenses/>.
//
// The PBC Go Wrapper makes use of The PBC library. The PBC Library and its use
// are covered under the terms of the GNU Lesser General Public License
// version 3, or (at your option) any later version.

package pbc

/*
#include <pbc/pbc.h>
#include "memstream.h"

int element_out_str_wrapper(char** bufp, size_t* sizep, int base, element_t e) {
	memstream_t* stream = pbc_open_memstream();
	if (stream == NULL) return 0;
	element_out_str(pbc_memstream_to_fd(stream), base, e);
	return pbc_close_memstream(stream, bufp, sizep);
}
*/
import "C"

import (
	"bytes"
	"fmt"
	"io"
	"unsafe"
)

func (el *Element) errorFormat(f fmt.State, c rune, err string) {
	fmt.Fprintf(f, "%%!%c(%s pbc.Element)", c, err)
}

func (el *Element) nativeFormat(f fmt.State, c rune) {
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
	if C.element_out_str_wrapper(&buf, &bufLen, C.int(base), el.cptr) == 0 {
		el.errorFormat(f, c, "INTERNALERROR")
		return
	}
	str := C.GoStringN(buf, C.int(bufLen))
	C.free(unsafe.Pointer(buf))
	fmt.Fprintf(f, "%s", str)
}

func (el *Element) customFormat(f fmt.State, c rune) {
	count := el.Len()
	if count == 0 {
		el.BigInt().Format(f, c)
	} else {
		fmt.Fprintf(f, "[")
		for i := 0; i < count; i++ {
			el.Item(i).customFormat(f, c)
			if i+1 < count {
				fmt.Fprintf(f, ", ")
			}
		}
		fmt.Fprintf(f, "]")
	}
}

// Format is a support routine for fmt.Formatter. It accepts many formats. The
// 'v' (value) and 's' (string) verbs will format the Element using the PBC
// library's internal formatting routines. These verbs accept variable widths
// to specify the base of the integers. Valid values are 2 to 36, inclusive.
//
// If the 'v' verb is used with the '#' (alternate format) flag, the output is
// metadata about the element in a pseudo-Go syntax. Checked elements will
// print more information than unchecked elements in this mode.
//
// If the 'd', 'b', 'o', 'x', or 'X' verbs are used, then the element is
// formatted within Go. The syntax approximates the PBC library's formatting,
// but integers are converted to big.Int for formatting. All of the verbs and
// flags that can be used in math/big will be used to format the elements.
func (el *Element) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('#') {
			if el.checked {
				fmt.Fprintf(f, "pbc.Element{Checked: true, Integer: %t, Field: %p, Pairing: %p, Addr: %p}", el.isInteger, el.fieldPtr, el.pairing, el)
			} else {
				fmt.Fprintf(f, "pbc.Element{Checked: false, Pairing: %p, Addr: %p}", el.pairing, el)
			}
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

// String converts el to a string using the default PBC library format.
func (el *Element) String() string {
	return fmt.Sprintf("%s", el)
}

// SetString sets el to the value contained in s. Returns (el, true) if
// successful, and (nil, false) if an error occurs. s is expected to be in the
// same format produced by String().
func (el *Element) SetString(s string, base int) (*Element, bool) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))

	if ok := C.element_set_str(el.cptr, cstr, C.int(base)); ok == 0 {
		return nil, false
	}
	return el, true
}

// Scan is a support routine for fmt.Scanner. It accepts the verbs 's' and 'v'
// only; only strings produced in PBC library format can be scanned. The width
// is used to denote the base of integers in the data.
func (el *Element) Scan(state fmt.ScanState, verb rune) error {
	// Verify verbs
	if verb != 's' && verb != 'v' {
		return ErrBadVerb
	}

	// Verify base
	base, ok := state.Width()
	if !ok {
		base = 10
	} else if base < 2 || base > 36 {
		return ErrBadVerb
	}

	// Compute valid integer symbols
	maxDigit := '9'
	maxAlpha := 'z'
	if base < 10 {
		maxDigit = rune('0' + (base - 1))
	}
	if base < 36 {
		maxAlpha = rune('a' + (base - 11))
	}

	state.SkipSpace()

	// Validate the input using a state machine (passing PBC invalid input is
	// likely to lead to bad outcomes)

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

	// The string seems valid; pass it to PBC
	if _, ok := el.SetString(buf.String(), base); !ok {
		return ErrBadInput
	}
	return nil
}

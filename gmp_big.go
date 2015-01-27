package pbc

/*
#cgo CFLAGS: -std=gnu99
#cgo LDFLAGS: -lgmp
#include <gmp.h>
*/
import "C"

import "math/big"
import "runtime"
import "unsafe"

type mpz_t struct {
	data C.mpz_t
}

var wordSize C.size_t
var bitsPerWord C.size_t

func clearmpz(x *mpz_t) {
	C.mpz_clear(&x.data[0])
}

func newmpz() *mpz_t {
	out := &mpz_t{}
	C.mpz_init(&out.data[0])
	runtime.SetFinalizer(out, clearmpz)
	return out
}

func big2mpz(num *big.Int) *mpz_t {
	words := num.Bits()
	out := newmpz()
	if len(words) > 0 {
		C.mpz_import(&out.data[0], C.size_t(len(words)), -1, wordSize, 0, 0, unsafe.Pointer(&words[0]))
	}
	return out
}

func mpz2big(num *mpz_t) (out *big.Int) {
	wordsNeeded := (C.mpz_sizeinbase(&num.data[0], 2) + (bitsPerWord - 1)) / bitsPerWord
	words := make([]big.Word, wordsNeeded)
	var wordsWritten C.size_t
	C.mpz_export(unsafe.Pointer(&words[0]), &wordsWritten, -1, wordSize, 0, 0, &num.data[0])
	out = &big.Int{}
	out.SetBits(words)
	return
}

func init() {
	var oneWord big.Word
	size := unsafe.Sizeof(oneWord)
	wordSize = C.size_t(size)
	bitsPerWord = C.size_t(8 * size)
}

package pbc

/*
#include <gmp.h>
*/
import "C"

import (
	"math/big"
	"runtime"
	"unsafe"
)

var wordSize C.size_t
var bitsPerWord C.size_t

func clearMpz(x *C.mpz_t) {
	C.mpz_clear(&x[0])
}

func newMpz() *C.mpz_t {
	out := &C.mpz_t{}
	C.mpz_init(&out[0])
	runtime.SetFinalizer(out, clearMpz)
	return out
}

// big2thisMpz imports the value of num into out
func big2thisMpz(num *big.Int, out *C.mpz_t) {
	words := num.Bits()
	if len(words) > 0 {
		C.mpz_import(&out[0], C.size_t(len(words)), -1, wordSize, 0, 0, unsafe.Pointer(&words[0]))
	}
}

// big2mpz allocates a new mpz_t and imports a big.Int value
func big2mpz(num *big.Int) *C.mpz_t {
	out := newMpz()
	big2thisMpz(num, out)
	return out
}

// mpz2thisBig imports the value of num into out
func mpz2thisBig(num *C.mpz_t, out *big.Int) {
	wordsNeeded := (C.mpz_sizeinbase(&num[0], 2) + (bitsPerWord - 1)) / bitsPerWord
	words := make([]big.Word, wordsNeeded)
	var wordsWritten C.size_t
	C.mpz_export(unsafe.Pointer(&words[0]), &wordsWritten, -1, wordSize, 0, 0, &num[0])
	out.SetBits(words)
}

// mpz2big allocates a new big.Int and imports an mpz_t value
func mpz2big(num *C.mpz_t) (out *big.Int) {
	out = &big.Int{}
	mpz2thisBig(num, out)
	return
}

func init() {
	var oneWord big.Word
	size := unsafe.Sizeof(oneWord)
	wordSize = C.size_t(size)
	bitsPerWord = C.size_t(8 * size)
}

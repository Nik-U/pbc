package pbc

/*
#include <pbc/pbc.h>

void installRandomHook();
void uninstallRandomHook();
*/
import "C"

import (
	cryptorand "crypto/rand"
	"io"
	"math/big"
	"math/rand"
	"unsafe"
)

var logging bool

func Logging() bool { return logging }

func SetLogging(log bool) {
	logging = log
	if log {
		C.pbc_set_msg_to_stderr(C.int(1))
	} else {
		C.pbc_set_msg_to_stderr(C.int(0))
	}
}

type RandomSource interface {
	Rand(limit *big.Int) *big.Int
}

var randomProvider RandomSource

func RandomProvider() RandomSource { return randomProvider }

func SetRandomProvider(provider RandomSource) {
	randomProvider = provider
	if provider == nil {
		C.uninstallRandomHook()
	} else {
		C.installRandomHook()
	}
}

//export goGenerateRandom
func goGenerateRandom(out, limit unsafe.Pointer) {
	outPtr := (*C.mpz_t)(out)
	limitPtr := (*C.mpz_t)(limit)
	r := randomProvider.Rand(mpz2big(limitPtr))
	big2thisMpz(r, outPtr)
}

type readerProvider struct {
	source io.Reader
}

func (provider readerProvider) Rand(limit *big.Int) (result *big.Int) {
	result, err := cryptorand.Int(provider.source, limit)
	if err != nil {
		panic(ErrEntropyFailure)
	}
	return
}

type randProvider struct {
	source *rand.Rand
}

func (provider randProvider) Rand(limit *big.Int) (result *big.Int) {
	result = &big.Int{}
	result.Rand(provider.source, limit)
	return
}

func SetCryptoRandom() { SetReaderRandom(cryptorand.Reader) }

func SetReaderRandom(reader io.Reader) {
	if reader == nil {
		panic(ErrIllegalNil)
	}
	SetRandomProvider(&readerProvider{reader})
}

func SetRandRandom(rand *rand.Rand) { SetRandomProvider(&randProvider{rand}) }

func SetDefaultRandom() { SetRandomProvider(nil) }

func init() {
	SetLogging(false)
	SetDefaultRandom()
}

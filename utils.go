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

// Logging returns true if PBC will send status messages to stderr.
func Logging() bool { return logging }

// SetLogging enables or disables sending PBC status messages to stderr.
// Messages are hidden by default.
func SetLogging(log bool) {
	logging = log
	if log {
		C.pbc_set_msg_to_stderr(C.int(1))
	} else {
		C.pbc_set_msg_to_stderr(C.int(0))
	}
}

// RandomSource generates random numbers for consumption by PBC. Rand returns a
// random integer in [0,limit).
type RandomSource interface {
	Rand(limit *big.Int) *big.Int
}

var randomProvider RandomSource

// RandomProvider returns the current random number source for use by PBC.
func RandomProvider() RandomSource { return randomProvider }

// SetRandomProvider sets the random number source for use by PBC. If provider
// is nil, then PBC will use its internal random number generator, which is the
// default mode. If provider is non-nil, then requests for random numbers will
// be serviced by Go instead of the internal C functions. This is slower, but
// provides greater control. Several convenience functions are provided to set
// common sources of random numbers.
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
	outMpz := &mpz{i: *(**C.mpz_t)(out)}
	limitMpz := &mpz{i: *(**C.mpz_t)(limit)}
	r := randomProvider.Rand(mpz2big(limitMpz))
	big2thisMpz(r, outMpz)
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

// SetCryptoRandom causes PBC to use the crypto/rand package with the globally
// shared rand.Reader as the source of random numbers.
func SetCryptoRandom() { SetReaderRandom(cryptorand.Reader) }

// SetReaderRandom causes PBC to use the crypto/rand package to generate random
// numbers using the given reader as an entropy source.
func SetReaderRandom(reader io.Reader) {
	if reader == nil {
		panic(ErrIllegalNil)
	}
	SetRandomProvider(&readerProvider{reader})
}

// SetRandRandom causes PBC to use the given source of random numbers.
func SetRandRandom(rand *rand.Rand) { SetRandomProvider(&randProvider{rand}) }

// SetDefaultRandom causes PBC to use its internal source of random numbers.
// This is the default mode of operation. Internally, PBC will attempt to read
// from /dev/urandom if it exists, or from the Microsoft Crypto API on Windows.
// If neither of these sources is available, the library will fall back to an
// insecure PRNG.
func SetDefaultRandom() { SetRandomProvider(nil) }

func init() {
	SetLogging(false)
	SetDefaultRandom()
}

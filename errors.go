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

import "errors"

var (
	ErrInvalidParamString = errors.New("invalid pairing parameters")
	ErrNoSuitableCurves   = errors.New("no suitable curves were found")
	ErrUnknownField       = errors.New("unchecked element initialized in unknown field")
	ErrIllegalOp          = errors.New("operation is illegal for elements of this type")
	ErrUncheckedOp        = errors.New("unchecked element passed to checked operation")
	ErrIncompatible       = errors.New("elements are from incompatible fields or pairings")
	ErrBadPairList        = errors.New("pairing product list is in an invalid format")
	ErrBadInput           = errors.New("invalid element format during scan")
	ErrBadVerb            = errors.New("invalid verb specified for scan")
	ErrIllegalNil         = errors.New("received nil when non-nil was expected")
	ErrOutOfRange         = errors.New("index out of range")
	ErrEntropyFailure     = errors.New("error while reading from entropy source")
	ErrHashFailure        = errors.New("error while hashing data")
	ErrInternal           = errors.New("a severe internal error has lead to possible memory corruption")
)

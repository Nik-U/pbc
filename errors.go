package pbc

import "errors"

var (
	ErrInvalidParamString = errors.New("invalid pairing parameters")
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
	ErrInternal           = errors.New("a severe internal error has lead to possible memory corruption")
)

# The PBC Go Wrapper [![Build Status](https://travis-ci.org/Nik-U/pbc.svg)](https://travis-ci.org/Nik-U/pbc) [![GoDoc](https://godoc.org/github.com/Nik-U/pbc?status.svg)](https://godoc.org/github.com/Nik-U/pbc)

Package pbc provides structures for building pairing-based cryptosystems. It
is a wrapper around the Pairing-Based Cryptography (PBC) Library authored by
Ben Lynn (https://crypto.stanford.edu/pbc/).

This wrapper provides access to all PBC functions. It supports generation of
various types of elliptic curves and pairings, element initialization, I/O,
and arithmetic. These features can be used to quickly build pairing-based or
conventional cryptosystems.

The PBC library is designed to be extremely fast. Internally, it uses GMP
for arbitrary-precision arithmetic. It also includes a wide variety of
optimizations that make pairing-based cryptography highly efficient. To
improve performance, PBC does not perform type checking to ensure that
operations actually make sense. The Go wrapper provides the ability to add
compatibility checks to most operations, or to use unchecked elements to
maximize performance.

Since this library provides low-level access to pairing primitives, it is
very easy to accidentally construct insecure systems. This library is
intended to be used by cryptographers or to implement well-analyzed
cryptosystems.

## Features
* 5 different pairing types
* Pairing generation
* Parameter export and import
* Element type checking
* Fast element arithmetic and pairing
* Element randomization
* Element export and import
* Automatic garbage collection
* Integration with `fmt`
* Integration with `math/big`

## Dependencies
This package must be compiled using cgo. It also requires the installation
of GMP and PBC. During the build process, this package will attempt to
include `gmp.h` and `pbc/pbc.h`, and then dynamically link to GMP and PBC.
Installation on Windows requires the use of MinGW.

## Documentation
For additional installation instructions and documentation, see
https://godoc.org/github.com/Nik-U/pbc

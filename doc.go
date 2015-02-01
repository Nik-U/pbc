/*
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
	very easy to construct insecure systems. This library is intended to be used
	by cryptographers or to implement well-analyzed cryptosystems.

	Pairings

	Cryptographic pairings are defined over three mathematical groups: G1, G2,
	and GT, where each group is typically of the same order r. Additionally, a
	bilinear map e maps a pair of elements — one from G1 and another from G2 —
	to an element in GT. This map e has the following additional property:

		For some generator g in G1, generator h in G2, and x and y in Zr:
		e(gˣ, hʸ) = e(g,h)ˣʸ

	If G1 == G2, then a pairing is said to be symmetric. Otherwise, it is
	asymmetric.	Pairings can be used to construct a variety of efficient
	cryptosystems.

	Supported Pairings

	The PBC library currently supports 5 different types of pairings, each with
	configurable parameters. These types are designated alphabetically, roughly
	in chronological order of introduction. Type A, D, E, F, and G pairings are
	implemented in the library. Each type has different time and space
	requirements. For more information about the types, see the documentation
	for the corresponding generator calls, or the PBC manual page at
	https://crypto.stanford.edu/pbc/manual/ch08s03.html.

	Dependencies

	This package must be compiled using cgo. It also requires the installation
	of GMP and PBC. During the build process, this package will attempt to
	include <gmp.h> and <pbc/pbc.h>, and then dynamically link to GMP and PBC.
	It also expects a POSIX-like environment for several C functions. For this
	reason, this package cannot be used in Windows without a POSIX compatibility
	layer and a gcc compiler.

	Most systems include a package for GMP. To install GMP in Debian / Ubuntu:

		sudo apt-get install libgmp-dev

	For an RPM installation with YUM:

		sudo yum install gmp

	For installation with FINK (http://www.finkproject.org/) on Mac OS X:

		sudo fink install gmp gmp-shlibs

	For more information or to compile from source, visit https://gmplib.org/

	To install the PBC library, download the appropriate files for your system
	from https://crypto.stanford.edu/pbc/download.html. The source can be
	compiled and installed using the usual GNU Build System:

		./configure
		make
		make install
*/
package pbc

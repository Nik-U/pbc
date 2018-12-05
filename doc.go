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
	very easy to accidentally construct insecure systems. This library is
	intended to be used by cryptographers or to implement well-analyzed
	cryptosystems.

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
	https://crypto.stanford.edu/pbc/manual/ch05s01.html.

	Dependencies

	This package must be compiled using cgo. It also requires the installation
	of GMP and PBC. During the build process, this package will attempt to
	include <gmp.h> and <pbc/pbc.h>, and then dynamically link to GMP and PBC.

	Most systems include a package for GMP. To install GMP in Debian / Ubuntu:

		sudo apt-get install libgmp-dev

	For an RPM installation with YUM:

		sudo yum install gmp-devel

	For installation with Fink (http://www.finkproject.org/) on Mac OS X:

		sudo fink install gmp gmp-shlibs

	For more information or to compile from source, visit https://gmplib.org/

	To install the PBC library, download the appropriate files for your system
	from https://crypto.stanford.edu/pbc/download.html. PBC has three
	dependencies: the gcc compiler, flex (http://flex.sourceforge.net/), and
	bison (https://www.gnu.org/software/bison/). See the respective sites for
	installation instructions. Most distributions include packages for these
	libraries. For example, in Debian / Ubuntu:

		sudo apt-get install build-essential flex bison

	The PBC source can be compiled and installed using the usual GNU Build
	System:

		./configure
		make
		sudo make install

	After installing, you may need to rebuild the search path for libraries:

		sudo ldconfig

	It is possible to install the package on Windows through the use of MinGW
	and MSYS. MSYS is required for installing PBC, while GMP can be installed
	through a package. Based on your MinGW installation, you may need to add
	"-I/usr/local/include" to CPPFLAGS and "-L/usr/local/lib" to LDFLAGS when
	building PBC. Likewise, you may need to add these options to CGO_CPPFLAGS
	and CGO_LDFLAGS when installing this package.

	License

	This package is free software: you can redistribute it and/or modify it
	under the terms of the GNU Lesser General Public License as published by
	the Free Software Foundation, either version 3 of the License, or (at your
	option) any later version.

	For additional details, see the COPYING and COPYING.LESSER files.
*/
package pbc

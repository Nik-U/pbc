// Copyright © 2015 Nik Unger
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

package pbc_test

import (
	"fmt"

	"github.com/Nik-U/pbc"
)

// This example generates a pairing and some random group elements, then applies
// the pairing operation.
func Example() {
	// In a real application, generate this once and publish it
	params := pbc.GenerateA(160, 512)

	pairing := params.NewPairing()

	// Initialize group elements. pbc automatically handles garbage collection.
	g := pairing.NewG1()
	h := pairing.NewG2()
	x := pairing.NewGT()

	// Generate random group elements and pair them
	g.Rand()
	h.Rand()
	fmt.Printf("g = %s\n", g)
	fmt.Printf("h = %s\n", h)
	x.Pair(g, h)
	fmt.Printf("e(g,h) = %s\n", x)
}

// This example displays an element in a variety of formats.
func ExampleElement_Format() {
	var element *pbc.Element

	// ...populate element...

	fmt.Printf("%v", element)    // Print in PBC format
	fmt.Printf("%s", element)    // Same as above
	fmt.Printf("%36v", element)  // Print in PBC format, base 36
	fmt.Printf("%#v", element)   // Print metadata about element
	fmt.Printf("%d", element)    // Print with Go
	fmt.Printf("%010o", element) // Print with Go, zero-padded width-10 octal
}

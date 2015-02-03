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

func ExampleElement_Format() {
	fmt.Printf("%v", element)    // Print in PBC format
	fmt.Printf("%s", element)    // Same as above
	fmt.Printf("%36v", element)  // Print in PBC format, base 36
	fmt.Printf("%#v", element)   // Print metadata about element
	fmt.Printf("%d", element)    // Print with Go
	fmt.Printf("%010o", element) // Print with Go, zero-padded width-10 octal
}

package pbc

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
	fmt.Printf("g = %v\nh = %v\n", g, h)
	x.Pair(g, h)
	fmt.Printf("e(g,h) = %v\n", x)
}

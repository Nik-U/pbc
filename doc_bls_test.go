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
	"crypto/sha256"
	"fmt"

	"github.com/Nik-U/pbc"
)

// messageData represents a signed message sent over the network
type messageData struct {
	message   string
	signature []byte
}

// This example computes and verifies a Boneh-Lynn-Shacham signature in a
// simulated conversation between Alice and Bob.
func Example_signBLS() {
	// The authority generates system parameters
	params := pbc.GenerateA(160, 512)
	pairing := params.NewPairing()
	g := pairing.NewG2().Rand()

	// The authority distributes params and g to Alice and Bob
	sharedParams := params.String()
	sharedG := g.Bytes()

	// Channel for messages. Normally this would be a network connection.
	messageChannel := make(chan *messageData)

	// Channel for public key distribution. This might be a secure out-of-band
	// channel or something like a web of trust. The public key only needs to
	// be transmitted and verified once. The best way to do this is beyond the
	// scope of this example.
	keyChannel := make(chan []byte)

	// Channel to wait until both simulations are done
	finished := make(chan bool)

	// Simulate the conversation participants
	go alice(sharedParams, sharedG, messageChannel, keyChannel, finished)
	go bob(sharedParams, sharedG, messageChannel, keyChannel, finished)

	// Wait for the communication to finish
	<-finished
	<-finished

	// Output:
	// Signature verified correctly
}

// Alice generates a keypair and signs a message
func alice(sharedParams string, sharedG []byte, messageChannel chan *messageData, keyChannel chan []byte, finished chan bool) {
	// Alice loads the system parameters
	pairing, _ := pbc.NewPairingFromString(sharedParams)
	g := pairing.NewG2().SetBytes(sharedG)

	// Generate keypair (x, g^x)
	privKey := pairing.NewZr().Rand()
	pubKey := pairing.NewG2().PowZn(g, privKey)

	// Send public key to Bob
	keyChannel <- pubKey.Bytes()

	// Some time later, sign a message, hashed to h, as h^x
	message := "some text to sign"
	h := pairing.NewG1().SetFromStringHash(message, sha256.New())
	signature := pairing.NewG2().PowZn(h, privKey)

	// Send the message and signature to Bob
	messageChannel <- &messageData{message: message, signature: signature.Bytes()}

	finished <- true
}

// Bob verifies a message received from Alice
func bob(sharedParams string, sharedG []byte, messageChannel chan *messageData, keyChannel chan []byte, finished chan bool) {
	// Bob loads the system parameters
	pairing, _ := pbc.NewPairingFromString(sharedParams)
	g := pairing.NewG2().SetBytes(sharedG)

	// Bob receives Alice's public key (and presumably verifies it manually)
	pubKey := pairing.NewG2().SetBytes(<-keyChannel)

	// Some time later, Bob receives a message to verify
	data := <-messageChannel
	signature := pairing.NewG1().SetBytes(data.signature)

	// To verify, Bob checks that e(h,g^x)=e(sig,g)
	h := pairing.NewG1().SetFromStringHash(data.message, sha256.New())
	temp1 := pairing.NewGT().Pair(h, pubKey)
	temp2 := pairing.NewGT().Pair(signature, g)
	if !temp1.Equals(temp2) {
		fmt.Println("*BUG* Signature check failed *BUG*")
	} else {
		fmt.Println("Signature verified correctly")
	}

	finished <- true
}

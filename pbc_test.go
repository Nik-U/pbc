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

package pbc

import (
	"crypto/sha256"
	"math/big"
	"runtime"
	"testing"
)

func testPairing(t *testing.T) *Pairing {
	// Generated with pbc_param_init_a_gen(p, 10, 32);
	pairing, err := NewPairingFromString("type a\nq 4025338979\nh 6279780\nr 641\nexp2 9\nexp1 7\nsign1 1\nsign0 1\n")
	if err != nil {
		t.Fatalf("Could not instantiate test pairing")
	}
	return pairing
}

func logElement(e *Element, name string, t *testing.T) {
	t.Logf("%s = %s\n", name, e)
}

// Boneh-Lynn-Shacham short signatures.
// Based on pbc/example/bls.c (C author: Ben Lynn).
func TestBLS(t *testing.T) {
	pairing := testPairing(t)

	g := pairing.NewG2()
	publicKey := pairing.NewG2()
	h := pairing.NewG1()
	sig := pairing.NewG1()
	temp1 := pairing.NewGT()
	temp2 := pairing.NewGT()
	secretKey := pairing.NewZr()

	// Generate system parameters
	g.Rand()
	logElement(g, "g", t)

	// Generate private key
	secretKey.Rand()
	logElement(secretKey, "secret key", t)

	// Compute corresponding public key
	publicKey.PowZn(g, secretKey)
	logElement(publicKey, "public key", t)

	// Generate element from a hash
	// For toy pairings, should check that pairing(g, h) != 1
	h.SetFromHash([]byte("hashofmessage"))
	logElement(h, "message hash", t)

	// h^secret_key is the signature
	// In real life: only output the first coordinate
	sig.PowZn(h, secretKey)
	logElement(sig, "signature", t)

	{
		sigBefore := sig.NewFieldElement().Set(sig)
		data := sig.CompressedBytes()
		sig.SetCompressedBytes(data)
		logElement(sig, "decompressed signature", t)
		if !sig.Equals(sigBefore) {
			t.Fatal("decompressed signature does not match")
		}
	}

	// Verification part 1
	temp1.Pair(sig, g)
	logElement(temp1, "f(sig,g)", t)

	// Verification part 2
	// Should match above
	temp2.Pair(h, publicKey)
	logElement(temp2, "f(hash,pubkey)", t)

	if !temp1.Equals(temp2) {
		t.Fatal("signature does not verify")
	}

	{
		data := sig.XBytes()
		sig.SetXBytes(data)

		temp1.Pair(sig, g)
		if temp1.Equals(temp2) {
			t.Log("signature verified on first try")
		} else {
			temp1.Invert(temp1)
			if temp1.Equals(temp2) {
				t.Log("signature verified on second try")
			} else {
				t.Fatal("signature does not verify")
			}
		}
	}

	// A random signature shouldn't verify
	sig.Rand()
	temp1.Pair(sig, g)
	if temp1.Equals(temp2) {
		t.Fatal("random signature verifies")
	}
}

// Hess ID-based signatures.
// Based on pbc/example/hess.c (C author: Dmitry Kosolapov).
// Based on paper "F. Hess. Efficient Identity Based Signature Schemes Based on
// Pairings. SAC 2002, LNCS 2595, Springer-Verlag, 2000"
func TestHess(t *testing.T) {
	pairing := testPairing(t)

	p := pairing.NewG1()
	p1 := pairing.NewG1()
	qid := pairing.NewG1()
	did := pairing.NewG1()
	ppub := pairing.NewG1()
	t4 := pairing.NewG1()
	t5 := pairing.NewG1()
	u := pairing.NewG1()

	s := pairing.NewZr()
	k := pairing.NewZr()
	v := pairing.NewZr()
	t8 := pairing.NewZr()

	r := pairing.NewGT()
	t1 := pairing.NewGT()
	t6 := pairing.NewGT()
	t7 := pairing.NewGT()

	// h is defined as sha256(m || r) where r is interpreted as bytes
	h := func(target *Element, message []byte, element *Element) {
		hash := sha256.New()
		hash.Write(message)
		hash.Write(element.Bytes())
		i := &big.Int{}
		target.SetBig(i.SetBytes(hash.Sum([]byte{})))
	}

	// Key generation
	p.Rand()
	s.Rand()
	qid.Rand()
	ppub.MulZn(p, s)
	did.MulZn(qid, s)
	logElement(qid, "Qid", t)
	logElement(p, "P", t)
	logElement(ppub, "Ppub", t)

	// Sign
	p1.Rand()
	k.Rand()
	t1.Pair(p1, p)
	r.PowZn(t1, k)
	h(v, []byte("Message"), r)
	t4.MulZn(did, v)
	t5.MulZn(p1, k)
	u.Add(t4, t5)
	logElement(u, "u", t)
	logElement(v, "v", t)

	// Verify
	t6.Pair(u, p)
	ppub.Neg(ppub)
	t7.Pair(qid, ppub)
	t7.PowZn(t7, v)
	r.Mul(t6, t7)
	h(t8, []byte("Message"), r)
	logElement(t8, "h3(m,r)", t)
	if !t8.Equals(v) {
		t.Fatal("signature does not verify")
	}
}

// Joux one-round protocol for tripartite Diffie-Hellman.
// Based on pbc/example/joux.c (C author: Dmitry Kosolapov).
// Based on paper "A. Joux. A One Round Protocol for Tripartie Diffie-Hellman.
// Proceedings of ANTS 4. LNCS 1838, pp. 385-394, 2000."
func TestJoux(t *testing.T) {
	pairing := testPairing(t)

	p := pairing.NewG1()
	t1 := pairing.NewG1()
	t2 := pairing.NewG1()
	t3 := pairing.NewG1()

	a := pairing.NewZr()
	b := pairing.NewZr()
	c := pairing.NewZr()

	t4 := pairing.NewGT()
	t5 := pairing.NewGT()
	t6 := pairing.NewGT()
	ka := pairing.NewGT()
	kb := pairing.NewGT()
	kc := pairing.NewGT()

	p.Rand()
	a.Rand()
	b.Rand()
	c.Rand()
	t1.MulZn(p, a)
	logElement(t1, "aP", t)
	t2.MulZn(p, b)
	logElement(t2, "bP", t)
	t3.MulZn(p, c)
	logElement(t3, "cP", t)

	t4.Pair(t2, t3)
	ka.PowZn(t4, a)
	logElement(ka, "Ka", t)
	t5.Pair(t1, t3)
	kb.PowZn(t5, b)
	logElement(kb, "Kb", t)
	t6.Pair(t1, t2)
	kc.PowZn(t6, c)
	logElement(kc, "Kc", t)

	if !ka.Equals(kb) || !kb.Equals(kc) {
		t.Fatal("shared key derivation failed")
	}
}

// Paterson ID-based signature.
// Based on pbc/example/paterson.c (C author: Dmitry Kosolapov).
// Based on paper "K. G. Paterson. ID-Based Signatures from Pairings on
// Elliptic Curves. Electron. Lett., Vol. 38". Available at
// http://eprint.iacr.org/2002/004."
func TestPaterson(t *testing.T) {
	pairing := testPairing(t)

	p := pairing.NewG1()
	ppub := pairing.NewG1()
	qid := pairing.NewG1()
	did := pairing.NewG1()
	r := pairing.NewG1()
	s1 := pairing.NewG1()
	t2 := pairing.NewG1()
	t4 := pairing.NewG1()
	t5 := pairing.NewG1()
	t7 := pairing.NewG1()

	s2 := pairing.NewZr()
	k := pairing.NewZr()
	t1 := pairing.NewZr()
	t3 := pairing.NewZr()

	t6 := pairing.NewGT()
	t8 := pairing.NewGT()
	t9 := pairing.NewGT()
	t10 := pairing.NewGT()
	t11 := pairing.NewGT()

	// Key generation
	p.Rand()
	s2.Rand()
	ppub.MulZn(p, s2)
	logElement(p, "P", t)
	logElement(ppub, "Ppub", t)
	qid.SetFromHash([]byte("ID"))
	logElement(qid, "Qid", t)
	did.MulZn(qid, s2)

	// Sign
	k.Rand()
	r.MulZn(p, k)
	t1.SetFromHash([]byte("Message"))
	t2.MulZn(p, t1)
	h := sha256.Sum256(r.Bytes())
	t3.SetFromHash(h[:])
	t4.MulZn(did, t3)
	t5.Add(t4, t2)
	k.Invert(k)
	s1.MulZn(t5, k)
	logElement(r, "R", t)
	logElement(s1, "S", t)

	// Verify
	t1.SetFromHash([]byte("Message"))
	t7.MulZn(p, t1)
	t6.Pair(p, t7)
	t8.Pair(ppub, qid)
	h = sha256.Sum256(r.Bytes())
	t3.SetFromHash(h[:])
	t9.PowZn(t8, t3)
	logElement(t8, "t8", t)
	logElement(t9, "t9", t)
	t10.Mul(t6, t9)
	logElement(t10, "[e(P, P)^H2(M)][e(Ppub, Qid)^H3(R)]", t)
	t11.Pair(r, s1)
	logElement(t11, "e(R, S)", t)
	if !t10.Equals(t11) {
		t.Fatal("signature does not verify")
	}
}

// Yuan-Li protocol ID-based AKE.
// Based on pbc/example/yuanli.c (C author: Dmitry Kosolapov).
// Based on paper "A New Efficient ID-Based Authenticated Key Agreement
// Protocol, Cryptology ePrint Archive, Report 2005/309"
func TestYuanLi(t *testing.T) {
	// This protocol has 2 stages: Setup and Extract. We represent them inside
	// one code block.
	pairing := testPairing(t)

	s := pairing.NewZr()
	a := pairing.NewZr()
	b := pairing.NewZr()

	p := pairing.NewG1()
	ppub := pairing.NewG1()
	qa := pairing.NewG1()
	qb := pairing.NewG1()
	sa := pairing.NewG1()
	sb := pairing.NewG1()
	ta := pairing.NewG1()
	tb := pairing.NewG1()
	temp1 := pairing.NewG1()
	temp2 := pairing.NewG1()
	temp3 := pairing.NewG1()
	h := pairing.NewG1()

	kab := pairing.NewGT()
	kba := pairing.NewGT()
	k := pairing.NewGT()
	temp4 := pairing.NewGT()
	temp5 := pairing.NewGT()

	// SETUP:
	//   KGS chooses G1, G2, e: G1*G1 -> G2, P, H: {0, 1}* -> G1, s,
	//     H - some function for key calculation.
	//   KGS calculates Ppub = s*P, publishes {G1, G2, e, P, Ppub, H1, H} and
	//     saves s as master key.
	p.Rand()
	logElement(p, "P", t)
	s.Rand()
	ppub.MulZn(p, s)
	logElement(ppub, "Ppub", t)

	// EXTRACT:
	//   For the user with ID public key can be calculated with Qid = H1(ID).
	//     KGS generates bound public key Sid = s*Qid.
	//   1. A chooses random a from Z_p*, calculates Ta = a*P.
	//     A -> B: Ta
	//   2. B chooses random b from Z_p*, calculates Tb = b*P.
	//     B -> A: Tb
	//   3. A calculates h = a*Tb = a*b*P and shared secret key
	//     Kab = e(a*Ppub + Sa, Tb + Qb)
	//   4. B calculates h = b*Ta = a*b*P and shared secret key
	//     Kba = e(Ta + Qa, b*Ppub + Sb)
	//   Session key is K = H(A, B, h, Kab).
	//   H was not defined in the original article.
	//   It is defined here as H(A, B, h, Kab)=e(h,H1(A)+H1(B))+Kab.
	qa.SetFromHash([]byte("A"))
	qb.SetFromHash([]byte("B"))
	sa.MulZn(qa, s)
	sb.MulZn(qb, s)
	logElement(sa, "Sa", t)
	logElement(sb, "Sb", t)

	// Step 1
	a.Rand()
	ta.MulZn(p, a)
	logElement(ta, "A->B Ta", t)

	// Step 2
	b.Rand()
	tb.MulZn(p, b)
	logElement(tb, "B->A Tb", t)

	// Step 3
	h.MulZn(tb, a)
	logElement(h, "h", t)
	temp1.MulZn(ppub, a)
	temp1.Add(temp1, sa)
	temp2.Add(tb, qb)
	kab.Pair(temp1, temp2)
	logElement(kab, "Kab", t)

	// Step 4
	h.MulZn(ta, b)
	logElement(h, "h", t)
	temp1.Add(ta, qa)
	temp2.MulZn(ppub, b)
	temp2.Add(temp2, sb)
	kba.Pair(temp1, temp2)
	logElement(kba, "Kba", t)

	// Conclusion
	temp3.Add(qa, qb)
	temp4.Pair(h, temp3)

	k.Add(temp4, kab)
	logElement(k, "A's key K", t)
	temp5.Set(k)

	k.Add(temp4, kba)
	logElement(k, "B's key K", t)

	if !temp5.Equals(k) {
		t.Fatalf("derived keys did not match")
	}
}

// Zhang-Kim ID-based Blind Signature scheme.
// Based on pbc/example/zhangkim.c (C author: Dmitry Kosolapov).
// Based on paper "F. Zang, K. Kim. ID-based Blind Signature and Ring Signature
// from Pairings. Advances in Cryptology - Asiacrypt 2002, LNCS Vol. 2510,
// Springer-Verlag, 2002"
func TestZhangKim(t *testing.T) {
	pairing := testPairing(t)

	p := pairing.NewG1()
	ppub := pairing.NewG1()
	qid := pairing.NewG1()
	sid := pairing.NewG1()
	r1 := pairing.NewG1()
	s1 := pairing.NewG1()
	t1 := pairing.NewG1()
	t2 := pairing.NewG1()
	t7 := pairing.NewG1()
	t8 := pairing.NewG1()
	t9 := pairing.NewG1()

	rr := pairing.NewZr()
	sr := pairing.NewZr()
	c := pairing.NewZr()
	a := pairing.NewZr()
	b := pairing.NewZr()
	negc := pairing.NewZr()
	t6 := pairing.NewZr()
	t14 := pairing.NewZr()

	t3 := pairing.NewGT()
	t10 := pairing.NewGT()
	t11 := pairing.NewGT()
	t12 := pairing.NewGT()

	// h is defined as sha256(m || r) where r is interpreted as bytes
	h := func(target *Element, message []byte, element *Element) {
		hash := sha256.New()
		hash.Write(message)
		hash.Write(element.Bytes())
		i := &big.Int{}
		target.SetBig(i.SetBytes(hash.Sum([]byte{})))
	}

	// Setup
	p.Rand()
	sr.Rand()
	ppub.MulZn(p, sr)
	logElement(p, "P", t)
	logElement(ppub, "Ppub", t)

	// Extract
	qid.SetFromHash([]byte("ID"))
	sid.MulZn(qid, sr)
	logElement(qid, "Public key Qid", t)
	logElement(sid, "Private key Sid", t)

	// Issue blind signature
	rr.Rand()
	r1.MulZn(p, rr)
	// Signer sends r1 = rr*P to user
	logElement(r1, "R", t)
	// Now we blind
	a.Rand()
	b.Rand()
	t1.MulZn(p, a)
	t1.Add(r1, t1)
	t2.MulZn(qid, b)
	t2.Add(t2, t1)
	t3.Pair(t2, ppub)
	h(t6, []byte("Message"), t3)
	c.Add(t6, b)
	// User sends c to signer
	logElement(c, "c", t)
	// Now we sign
	t7.MulZn(ppub, rr)
	t8.MulZn(sid, c)
	s1.Add(t8, t7)
	// Signer sends s1
	logElement(s1, "S", t)
	// Now we unblind
	t9.MulZn(ppub, a)
	s1.Add(s1, t9)
	c.Sub(c, b)
	// Blind signature is now (S, c)
	logElement(s1, "S1", t)
	logElement(c, "c1", t)

	// Verification
	t10.Pair(qid, ppub)
	negc.Neg(c)
	t10.PowZn(t10, negc)
	t11.Pair(s1, p)
	t12.Mul(t11, t10)
	h(t14, []byte("Message"), t12)
	logElement(c, "c1", t)
	logElement(t14, "H(m, [e(S1, P)][e(Qid, Ppub)^(-c1)])", t)
	if !t14.Equals(c) {
		t.Fatal("signature does not verify")
	}
}

// TestGC ensures that there are no errors when running struct finalizers.
func TestGC(t *testing.T) {
	TestBLS(t)
	for i := 0; i < 5; i++ { // Multiple rounds to resolve dependencies
		runtime.GC()
	}
}

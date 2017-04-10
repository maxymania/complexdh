/*
MIT License

Copyright (c) 2017 Simon Schmidt

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
This Package implements operations for large complex numbers modulo P.
The functions can be used for a Diffie-Hellman key exchange, where P is usually a prime number.
*/
package complexdh

import "math/big"

func multiply(dr, di, ar, ai, br, bi, modulus *big.Int) {
	var t1,t2,t3,t1b big.Int
	
	/*
	dr = ar*br - ai*bi   mod m
	di = ar*bi + ai*br   mod m
	*/
	r := t1.Sub(t2.Mul(ar,br),t3.Mul(ai,bi))
	i := t1b.Add(t2.Mul(ar,bi),t3.Mul(ai,br))
	dr.Mod(r, modulus)
	di.Mod(i, modulus)
}

/* In-place-self-multiply */
func ips_multiply(ar, ai, modulus *big.Int) {
	var t big.Int
	
	/*
	dr = ar*ar - ai*ai   mod m
	di = ar*ai * 2       mod m
	
	t = ar*ai
	ar *= ar
	ai *= ai
	ar -= ai
	
	ai = t*2 = t<<1
	*/
	
	T := t.Mul(ar,ai)
	
	ar.Mul(ar,ar)
	ai.Mul(ai,ai)
	ar.Sub(ar,ai)
	
	ai.Lsh(T,1)
	
	ar.Mod(ar,modulus)
	ai.Mod(ai,modulus)
}

type ModulusGroup struct{
	Modulus *big.Int
	Gr, Gi *big.Int
}

func (m *ModulusGroup) Multiply(dr, di, ar, ai, br, bi *big.Int) {
	multiply(dr,di,ar,ai,br,bi,m.Modulus)
}
func (m *ModulusGroup) Exponent(dr,di, sr,si *big.Int, scalar []byte) {
	dr.SetUint64(1)
	di.SetUint64(0)
	for _,k := range scalar {
		ips_multiply(dr,di,m.Modulus)
		if (k&0x80)==0x80 {
			multiply(dr,di,dr,di,sr,si,m.Modulus)
		}
	}
}
func (m *ModulusGroup) BaseExp(dr,di *big.Int, scalar []byte) {
	m.Exponent(dr,di,m.Gr,m.Gi,scalar)
}


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

package complexdh

import "math/big"
import "crypto/elliptic"

func sqrt(n uint) uint {
	if n<0 { return -sqrt(-n) }
	if n<2 { return n }
	sc := sqrt(n>>2)<<1
	lc := sc+1
	if (lc*lc)>n { return sc }
	return lc
}
func dropBits(bl uint) uint {
	//sq := sqrt(bl)*45
	sq := sqrt(sqrt(bl)*45)*10
	//sq := uint(2025)
	if sq>=bl { return 0 }
	return bl-sq
}

type asCurve struct{
	*ModulusGroup
}

func (a *asCurve) Params() *elliptic.CurveParams {
	n := new(elliptic.CurveParams)
	n.P = a.Modulus
	n.N = new(big.Int).Rsh(n.P,dropBits(uint(n.P.BitLen())))
	n.B = new(big.Int)
	n.Gx = a.Gr
	n.Gy = a.Gi
	n.BitSize = n.N.BitLen()
	n.Name = "N/A(complex)"
	return  n
}
func (a *asCurve) IsOnCurve(x, y *big.Int) bool { return true }

func (a *asCurve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	x = new(big.Int)
	y = new(big.Int)
	a.Multiply(x,y,x1,y1,x2,y2)
	return
}

func (a *asCurve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	x = new(big.Int).SetBits(x1.Bits())
	y = new(big.Int).SetBits(y1.Bits())
	ips_multiply(x,y,a.Modulus)
	return
}

func (a *asCurve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	x = new(big.Int)
	y = new(big.Int)
	a.Exponent(x,y,x1,y1,k)
	return
}
func (a *asCurve) ScalarBaseMult(k []byte) (x, y *big.Int) {
	x = new(big.Int)
	y = new(big.Int)
	a.BaseExp(x,y,k)
	return
}

/*
Returns an elliptic.Curve-Wrapper around this group. The *elliptic.CurveParams returned
by the .Params() method the Curve should not be used for doing ScalarMult, etc.!
*/
func (m *ModulusGroup) AsCurve() elliptic.Curve {
	return &asCurve{m}
}


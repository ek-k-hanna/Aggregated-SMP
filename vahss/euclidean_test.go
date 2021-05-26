/*
 * Copyright (C) 2021 Hanna Ek
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package vahss

import(
  "testing"
  "fmt"
  "math/big"
)

func TestGcd(t *testing.T){

  gcd := Gcd(big.NewInt(7),big.NewInt(9))
  isSameGcd(t,big.NewInt(1),gcd)

  gcd = Gcd(big.NewInt(8),big.NewInt(18))
  isSameGcd(t,big.NewInt(2),gcd)

  gcd = Gcd(big.NewInt(-12),big.NewInt(24))
  isSameGcd(t,big.NewInt(-12),gcd)

  gcd = Gcd(big.NewInt(12),big.NewInt(24))
  isSameGcd(t,big.NewInt(12),gcd)

  gcd = Gcd(big.NewInt(4864),big.NewInt(3458))
  isSameGcd(t,big.NewInt(38),gcd)
}

func TestEuclidean(t *testing.T){
  a,b,c := big.NewInt(32), big.NewInt(-45), big.NewInt(38)
  x,y,z := Extended_euclidean_algorithm(big.NewInt(4864), big.NewInt(3458))
  ok := a.Cmp(x)==0 && b.Cmp(y)==0 && c.Cmp(z)==0
  if ok != true {
      t.Errorf("Assert failure: expected true, actual: %t", ok)
  }

  a,b,c = big.NewInt(-45), big.NewInt(32), big.NewInt(38)
  x,y,z = Extended_euclidean_algorithm(big.NewInt(3458), big.NewInt(4864))
  ok = a.Cmp(x)==0 && b.Cmp(y)==0 && c.Cmp(z)==0
  if ok != true {
      t.Errorf("Assert failure: expected true, actual: %t", ok)
  }

}

func TestGcdAndMod(t *testing.T){
   testGcdAndMod(1,1,0,2,t)
    testGcdAndMod(1,1,1,2,t)
    testGcdAndMod(0,2,2,2,t)

    testGcdAndMod(6,6,14,7,t)
    testGcdAndMod(2,6,9,7,t)

    testGcdAndMod(38,4864,3458,9923,t)

}

func   testGcdAndMod(ai,bi,ci, modi int,t *testing.T){
  a,b,c,mod := int64(ai), int64(bi), int64(ci),int64(modi)
  modP := InitModular(big.NewInt(mod))
  LHS := InitIntegerModP(big.NewInt(a),modP).num
  RHS := Gcd(InitIntegerModP(big.NewInt(b),modP).num,InitIntegerModP(big.NewInt(c),modP).num)
  isSameGcd(t,LHS,RHS)
}

func TestEuclideanAndMod(t *testing.T){
   testEuclideanAndMod(int64(32),int64(-45),int64(38),int64(4864),int64(3458),int64(9923),t)

}

func   testEuclideanAndMod(a,b,c,x,y,mod int64,t *testing.T){
  modP := InitModular(big.NewInt(mod))
  LHSa := InitIntegerModP(big.NewInt(a),modP).num
  LHSb := InitIntegerModP(big.NewInt(b),modP).num
  LHSc := InitIntegerModP(big.NewInt(c),modP).num
  RHSx,RHSy,RHSz := Extended_euclidean_algorithm(InitIntegerModP(big.NewInt(x),modP).num,InitIntegerModP(big.NewInt(y),modP).num)
  ok := LHSa.Cmp(RHSx)==0 && LHSb.Cmp(InitIntegerModP(RHSy,modP).num)==0 && LHSc.Cmp(RHSz)==0
  //ok := LHSa.Cmp(RHSx)==0 && LHSb.Cmp(RHSy)==0 && LHSc.Cmp(RHSz)==0

  if ok != true {
      t.Errorf("Assert failure: expected true, actual: %t", ok)
      fmt.Println("Test a==x?:",LHSa,RHSx)
      fmt.Println("Test b==y?:",LHSb,InitIntegerModP(RHSy,modP).num)
      fmt.Println("Test c==z?:",LHSc,RHSz)
  }
}
/*
*/
func isSameGcd(t *testing.T, x, y *big.Int){
  ok := x.Cmp(y)==0
  if ok != true {
      t.Errorf("Assert failure: expected true, actual: %t", ok)
  }

}

func isSameEuclidean(t *testing.T, x, y ,z *big.Int){
  ok := x.Cmp(y)==0
  if ok != true {
      t.Errorf("Assert failure: expected true, actual: %t", ok)
  }
}

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

import (
  "testing"
  "math/big"
)

func TestModP(t *testing.T) {
  mod7 := InitModular(big.NewInt(7))
  num1 := InitIntegerModP(big.NewInt(5),mod7)
  num2 := InitIntegerModP(big.NewInt(5),mod7)
  isSame(t,num1, num2) // Sanity check


  num1 = InitIntegerModP(big.NewInt(5),mod7)
  divisor := InitIntegerModP(big.NewInt(3),mod7)
  num2 = ModInverse(divisor)
  isSame(t,num1, num2) // Sanity check


  num1 = InitIntegerModP(big.NewInt(1),mod7)
  fak1 := InitIntegerModP(big.NewInt(3),mod7)
  fak2 := InitIntegerModP(big.NewInt(5),mod7)
  num2 = ModMul(fak1,fak2)
  isSame(t,num1, num2) // Sanity check


  num1 = InitIntegerModP(big.NewInt(3),mod7)
  num2 = InitIntegerModP(big.NewInt(3),mod7)
  test1 := num1.num.Cmp(new(big.Int).Mul(num2.num,big.NewInt(1))) == 0
  if test1 != true {
      t.Errorf("Assert failure: expected true, actual: %t", test1)
  }


  num1 = InitIntegerModP(big.NewInt(2),mod7)
  term1 := InitIntegerModP(big.NewInt(5),mod7)
  term2 := InitIntegerModP(big.NewInt(4),mod7)
  num2 = ModAdd(term1,term2)
  isSame(t,num1, num2) // Sanity check


  term1 = InitIntegerModP(big.NewInt(3),mod7)
  term2 = InitIntegerModP(big.NewInt(4),mod7)
  num1 = InitIntegerModP(big.NewInt(0),mod7)
  num2 = ModAdd(term1,term2)
  isSame(t, num1, num2)
}

func isSame(t *testing.T, x *IntegerModP, y *IntegerModP) {
  ok := ModEq(x,y)
  if ok != true {
      t.Errorf("Assert failure: expected true, actual: %t", ok)
  }
}

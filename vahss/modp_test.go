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

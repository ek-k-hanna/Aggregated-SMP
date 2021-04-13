package vahss

import(
  "math/big"
)

type Modular struct{
  P *big.Int
}

type IntegerModP struct {
  Num     *big.Int
  mod   *Modular
}

func GetNumber(field_object *IntegerModP) *big.Int{
  return field_object.Num
}
func GetMod(field_object *IntegerModP) *big.Int{
  return field_object.mod.P
}
func GetField(mod *Modular) *big.Int{
  return mod.P
}
func InitModular(mod *big.Int) (*Modular){
  field := new(Modular)
  field.P = mod
  return field
}

func InitIntegerModP(Number *big.Int, modular *Modular) (*IntegerModP){
    var field_object = new(IntegerModP)
    field_object.Num = Modolus(Number, modular.P)
    field_object.mod = modular
    return field_object
}

func Modolus(Number *big.Int, mod *big.Int) *big.Int {
  var res *big.Int = new(big.Int).Mod(Number,mod) //OBS!
  //if ((res < int64(0) && mod > int64(0) ) || (res > int64(0) && mod < int64(0)) ){
  //  return res + mod
  //}
  return res
}

func ModAdd(self, other *IntegerModP) (*IntegerModP){
  var field_object = InitIntegerModP( new(big.Int).Add(self.Num,other.Num), self.mod)
  return field_object
}

func ModSub(self, other *IntegerModP) (*IntegerModP){
  var field_object = InitIntegerModP( new(big.Int).Sub(self.Num,other.Num), self.mod)
  return field_object
}

func ModMul(self, other *IntegerModP) (*IntegerModP){
  var mul *big.Int = new(big.Int).Mul(self.Num,other.Num)
  var field_object = InitIntegerModP( mul , self.mod)
  return field_object
}

func ModPow(self, other *IntegerModP) (*IntegerModP){
  var powNum *big.Int =  new(big.Int).Exp(self.Num,other.Num,nil)
  var field_object = InitIntegerModP(powNum, self.mod)
  return field_object
}


func ModNeg(self *IntegerModP) (*IntegerModP){
  var field_object = InitIntegerModP( new(big.Int).Neg(self.Num) , self.mod)
  return field_object
}

func ModEq(self, other *IntegerModP) (bool){
  var sameNum bool = ( self.Num.Cmp(other.Num) == 0)
  var sameMod bool = ( self.mod.P.Cmp(other.mod.P) == 0)
  return (sameNum && sameMod)
}

func ModNe(self, other *IntegerModP) (bool){
  return !(ModEq(self,other))
}

func ModDiv(self, divisor *IntegerModP ) (*IntegerModP){
  field_object,_ := ModDivmod(self,divisor)
  return field_object
}

func ModDivmod(self, divisor *IntegerModP ) (*IntegerModP, *IntegerModP){
  var q,r *big.Int = new(big.Int).Div(self.Num,divisor.Num), new(big.Int).Mod(self.Num,divisor.Num)
  field_object_div := InitIntegerModP(q, self.mod)
  field_object_reminder := InitIntegerModP(r, self.mod)
  return field_object_div, field_object_reminder
}

func ModInverse(self *IntegerModP) (*IntegerModP){
  x,_,_ := Extended_euclidean_algorithm(self.Num, self.mod.P)
  var field_object = InitIntegerModP(x, self.mod)
  return field_object
}

func ModAbs (self *IntegerModP) (*big.Int){
  return Modolus(new(big.Int).Abs(self.Num),self.mod.P)
}

func ModInt (self *IntegerModP) (*big.Int){
  return self.Num
}

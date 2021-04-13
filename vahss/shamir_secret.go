package vahss

import(
  _"math"
  "math/big"
  "crypto/rand"
  _"fmt"
)
/*
func Egcd(a,b int64) (int64, int64, int64){
  if a == 0 {
    return b, 0, 1
  }else{
    g,y,x := Egcd(Modolus(b,a), a )
    return g, x - ( ( b / a ) * y ), y // floor division?
  }
}

func Mod_inverse(k, prime int64) (int64) {
  k = Modolus(k,prime)
  var r int64
  if k < 0 {
    _,_,r = Egcd(prime,-k)
  }else{
    _,_,r = Egcd(prime,k)
  }
  return Modolus( (prime + r ), prime )
}
*/
func Split_points(points [][2]*IntegerModP )([]*IntegerModP,[]*IntegerModP ){
  var x_values  []*IntegerModP
  var y_values  []*IntegerModP
  for _,p := range points {
    x_values = append(x_values , p[0])
    y_values = append(y_values , p[1])
  }
  return x_values, y_values

}

func shamir_random_polynomial(degree int64, secret, modolus *big.Int) ([]*IntegerModP){
 var coefficients []*IntegerModP
 Zp := InitModular(modolus)
 var max *big.Int = modolus

 coefficients = append(coefficients, InitIntegerModP(secret,Zp) )
 for i:=int64(1); i < degree; i++ {
   randInt,err := rand.Int(rand.Reader, max)
   if err != nil{
     panic(err)
   }
   coefficient := InitIntegerModP(randInt,Zp)
   coefficients = append(coefficients,coefficient)
 }
 randomMonicPolynomial := append(coefficients, InitIntegerModP(big.NewInt(1),Zp ) )
 return randomMonicPolynomial
}

func get_polynomial_points(coefficients []*IntegerModP, degree int64, num_points int64) ([][2]*IntegerModP){
   var points [][2]*IntegerModP
   mod := coefficients[0].mod
   for x:=int64(1); x <= num_points; x++{
     x_Zp := InitIntegerModP(big.NewInt(x),mod)
     y_Zp := coefficients[0]

      for i := int64(1) ; i < degree ; i++{
         i_Zp := InitIntegerModP(big.NewInt(i),mod)
         exponentiation := ModPow(x_Zp,i_Zp)
         term := ModMul(coefficients[i],exponentiation)
         y_Zp = ModAdd(y_Zp,term)
       }
      var point = [2]*IntegerModP{ x_Zp, y_Zp }
      points = append(points, point)
    }
   return points
 }

 func Lagrange_coeffs(secret *big.Int, points [][2]*IntegerModP, prime *big.Int) ([]*IntegerModP){
   x_values, y_values := Split_points(points)

   var coeffs []*IntegerModP
   //var lambdas_ijs []*IntegerModP
   nr_servers := int64(len(x_values))
   Zp := x_values[0].mod
   sum := InitIntegerModP(big.NewInt(0),Zp)
   for i := int64(1); i < nr_servers + 1 ; i++ {
     //denominator := InitIntegerModP(1,mod)
     lambda_ij := InitIntegerModP(big.NewInt(1),Zp)
     tmp_i := InitIntegerModP(big.NewInt(i),Zp)
     for j := int64(1) ; j< nr_servers + 1 ; j++{
       if j != i {
         tmp_j := InitIntegerModP(big.NewInt(j),Zp)
         denominator := InitIntegerModP( ModSub(tmp_j,tmp_i).Num ,Zp)
         denominator = ModInverse(denominator)//Mod_inverse(denominator,prime)
         tmp_j = ModMul(tmp_j,denominator)
         lambda_ij = ModMul(tmp_j,lambda_ij)// InitIntegerModP(lambda_ij.num*j*denominator_tmp.num,Zp)
       }
     }
     coeff := ModMul(lambda_ij, y_values[i-1])
     coeffs = append(coeffs, coeff)
     sum = ModAdd(sum, coeff)
     //lambdas_ijs = append(lambdas_ijs, lambda_ij)
   }
   //fmt.Println(sum.num)
   //for i := int64(0); i < nr_servers ; i++{
    // coeff := ModMul(lambdas_ijs[i],y_values[i])
    // coeffs = append(coeffs, coeff)
   //}
   return coeffs
 }


func Generate_input_shamir_secret(secret_int *big.Int, degree int64, num_points int64, prime *big.Int) ([]*IntegerModP) {
  coefficients := shamir_random_polynomial(degree-1, secret_int, prime)
  points := get_polynomial_points(coefficients,degree,num_points) // [][]IntegerModP
  lagrange_pol := Lagrange_coeffs(secret_int, points, prime)
  return lagrange_pol
}

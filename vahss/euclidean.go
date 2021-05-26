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
  "math/big"
)

func Gcd(a,b *big.Int) (*big.Int) {
  if new(big.Int).Abs(a).Cmp(new(big.Int).Abs(b)) == -1 {
    return Gcd(b,a)
  }
  for new(big.Int).Abs(b).Cmp(big.NewInt(0)) == 1 {
     var r *big.Int = Modolus(a,b)
     a, b = b, r
   }
  return a
}

func Extended_euclidean_algorithm(a, b *big.Int) (*big.Int,*big.Int,*big.Int){

  if ( new(big.Int).Abs(b).Cmp( new(big.Int).Abs(a) ) >= 0 ){
    var x,y,d *big.Int = Extended_euclidean_algorithm(b, a)
    return y,x,d
  }

  if new(big.Int).Abs(b).Cmp(big.NewInt(0)) == 0 {
    return big.NewInt(1), big.NewInt(0), a
  }

  var x1, x2, y1, y2 *big.Int
  x1, x2, y1, y2 = big.NewInt(0), big.NewInt(1), big.NewInt(1),big.NewInt(0)
  for new(big.Int).Abs(b).Cmp(big.NewInt(0)) == 1 {
     var q, r *big.Int = new(big.Int).Div(a,b), new(big.Int).Mod(a,b)
     x := new(big.Int).Sub(x2, new(big.Int).Mul(q,x1) )
     y := new(big.Int).Sub(y2, new(big.Int).Mul(q,y1) )
     a, b, x2, x1, y2, y1 = b, r, x1, x, y1, y
   }
  return x2, y2, a
}

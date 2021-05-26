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
  "math/big"
  _"sort"
)

func calculate_mersenne_primes() ([]*big.Int){
  mersenne_prime_exponents := []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(5), big.NewInt(7), big.NewInt(13), big.NewInt(17), big.NewInt(19), big.NewInt(31), big.NewInt(61), big.NewInt(89), big.NewInt(107), big.NewInt(127), big.NewInt(521), big.NewInt(607), big.NewInt(1279)}
    var primes []*big.Int
    for _, s := range mersenne_prime_exponents{
      prime :=  big.NewInt(1)
      for j:= int64(0); j < s.Int64(); j++ {
        prime.Mul(prime,big.NewInt(2))
      }
      prime.Sub(prime,big.NewInt(1))
      primes = append(primes,prime)
    }
  return primes
}

func Get_large_enough_prime(batch []int64) (*big.Int){
  var SMALLEST_257BIT_PRIME, SMALLEST_321BIT_PRIME,SMALLEST_385BIT_PRIME *big.Int
  SMALLEST_257BIT_PRIME =  new(big.Int).Exp(big.NewInt(2),big.NewInt(256),nil)
  SMALLEST_257BIT_PRIME.Add(SMALLEST_257BIT_PRIME,big.NewInt(297))
  SMALLEST_321BIT_PRIME =  new(big.Int).Exp(big.NewInt(2),big.NewInt(320),nil)
  SMALLEST_321BIT_PRIME.Add(SMALLEST_321BIT_PRIME,big.NewInt(27))
  SMALLEST_385BIT_PRIME =  new(big.Int).Exp(big.NewInt(2),big.NewInt(384),nil)
  SMALLEST_385BIT_PRIME.Add(SMALLEST_385BIT_PRIME,big.NewInt(231))

  large_primes := []*big.Int{SMALLEST_257BIT_PRIME, SMALLEST_321BIT_PRIME, SMALLEST_385BIT_PRIME}
  STANDARD_PRIMES := append(calculate_mersenne_primes(),large_primes...)
  primes := STANDARD_PRIMES

  for _,p := range primes{
    var numbers_greater_than_prime []*big.Int
    for _, b := range(batch){
      b_big := big.NewInt(b)
      if (b_big.Cmp(p) == 1) { // if b > p
        numbers_greater_than_prime = append(numbers_greater_than_prime, b_big)
      }
    }

    if ( len(numbers_greater_than_prime) == 0 ){
      return p
    }
  }
  return nil
}

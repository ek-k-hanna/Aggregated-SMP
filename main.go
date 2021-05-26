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
package main

import(

  // Bulletproofs
  /*
  ."hannaekthesis/vahss_bprp"
  _"hannaekthesis/bulletproof"
  */

  // Set membership proofs

  _"hannaekthesis/ccs08"
  ."hannaekthesis/vahss_SM"
  

  //Signature based range proofs
  /*
  //
  ."hannaekthesis/vahss_UL"
    _"hannaekthesis/vahss"
  */
)
func main(){
  // Bulletproofs
  // Main_bprp()

  // Set membership proofs
   Main_SM()

  //Signature based range proofs
  // Main_UL()

}

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

type Server struct {
    Id     int
    Shares []*big.Int
    mod *Modular
}

 func InitServer(j int, mod *Modular) (*Server){
   s_j := new(Server)
   s_j.Id = j
   s_j.mod = mod
   s_j.Shares = nil
   return s_j
 }

 func Set_share(s_j *Server, i int, share *big.Int){
   s_j.Shares = append(s_j.Shares,share)
  }
func GetServerId(s_j *Server) (int){
  return s_j.Id
}

func GetShares(s_j *Server) ([]*big.Int){
  return s_j.Shares
}

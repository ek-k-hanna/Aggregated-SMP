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
type Client struct {
    Id     int
    t     int64
    Xi   *big.Int
    Ri   *big.Int
    //tau_i *p256.P256
    mod  *Modular
}

func GetClientId(c_i *Client) int{
  return c_i.Id
}

func InitClient(id_client int, secret_input *big.Int, t int64, r_i *big.Int, mod *Modular) (*Client){
    c_i := new(Client)
    c_i.Id = id_client
    c_i.Xi = secret_input
    c_i.t = t
    c_i.Ri = r_i
    c_i.mod = mod
    return c_i
  }

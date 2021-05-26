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
package vahss_SM

import(
  "math/big"
  "fmt"
  "github.com/ing-bank/zkrp/crypto/bn256"
  "hannaekthesis/ccs08"
  ."hannaekthesis/vahss"
  "crypto/rand"

)
func Generate_set() (ccs08.ParamsSet){
  n:=10
  max := big.NewInt(2500)
  set := make([]int64,n)
  set[0] = 12
  set[1] = 60
  for i:=2; i<n; i++{
      value,_ := rand.Int(rand.Reader, max)
      set[i] = value.Int64()
  }
  p, _ := ccs08.SetupSet(set)
  return p
}

func Main_SM(){

  const nr_clients = int64(100)
  const nr_servers = int64(5)
  var nr_aggregators int64 = int64(2)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t int64
  var clients []*Client
  var servers []*Server
  var set ccs08.ParamsSet

  // Determines if use agrgeated set membership proof or non aggrgeated
  var aggregate bool = true

  seedGroup := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(seedGroup, 10)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2
  R_is = big.NewInt(0)
  set = Generate_set()

  // Initiate cients
  for i:=int64(1) ; i<= nr_clients ; i++{
    id := int(1)
    if (i!= nr_clients){
      secret := big.NewInt(12)
      R_i,_ := rand.Int(rand.Reader, phiN)
      c_i := InitClient(id,secret,t,R_i,Zp)
      R_is = ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)).Num
      clients = append(clients, c_i)
    }else{
      secret := big.NewInt(12)
      R_n := InitIntegerModP(new(big.Int).Sub( new(big.Int).Mul(phiN, new(big.Int).Div( R_is, phiN ) ) ,  R_is),Zp).Num
      c_n := InitClient(id, secret, t, R_n, Zp)
      clients = append(clients, c_n)
    }
  }


  // Initiate Servers
  for j:= int64(1); j<= nr_servers ; j++{
    id := int(j)
    s_j := InitServer(id,Zp)
    servers = append(servers,s_j)
  }

  var tau_is []*bn256.G2
  var range_proofs []ccs08.ProofSet

  for i:=int64(0); i< nr_clients ; i++{
    id := int(i)
    c_i := clients[id]
    var shares []*big.Int = Generate_shares(c_i.Id, c_i.Xi, t, nr_servers, Zp)
    //var rp ccs08.ProofSet
    rp , _ := ccs08.ProveSet(c_i.Xi.Int64(),c_i.Ri, set)
    range_proofs = append(range_proofs,rp)
    var tau_i *bn256.G2
    tau_i = rp.C
    tau_is = append(tau_is, tau_i)


    for j:= int64(1); j<=nr_servers ; j++{
      Set_share(servers[j-1], c_i.Id, shares[j-1] )
    }
  }


  var y_js []*big.Int
  var sigma_js []*bn256.G2

  for j:= int64(0); j < nr_servers ; j++{
    s_j := servers[j]
    y_js = append(y_js, Partial_eval(s_j, s_j.Id, s_j.Shares, nr_clients, Zp))
    sigma_js = append(sigma_js,Partial_proof(s_j, s_j.Shares,nr_clients, Zp))
  }

  y := Final_eval(y_js, nr_servers, Zp)
  sigma := Final_proof(sigma_js, nr_servers)
  sum := InitIntegerModP(y,Zp).Num

  var ok bool
  if aggregate {
    ok = PartialVerify(tau_is, nr_clients, sigma, y, range_proofs, set, nr_aggregators)
  }else{
    ok = Verify(tau_is, nr_clients, sigma, y, range_proofs, set)
  }

  if ok {
    fmt.Println("Verify ok")
      fmt.Println("Sum is", sum)
  }else{
    fmt.Println("Verification ERROR")
  }
}

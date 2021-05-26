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
package vahss_UL
import(
  "math/big"
  "fmt"
  "github.com/ing-bank/zkrp/crypto/bn256"
  ."hannaekthesis/ccs08"
  ."hannaekthesis/vahss"
  "crypto/rand"

)
func Main_UL(){
  const nr_clients = int64(100)
  const nr_servers = int64(5)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t, min,max int64
  var clients []*Client
  var servers []*Server
  var zkrp Ccs08

  //prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value, 10)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2
  min, max = 18, 200
  length_interval := max - min
  R_is = big.NewInt(0)

  zkrp.Setup(min,max)

  // Initiate cients
  for i:=int64(1) ; i<= nr_clients ; i++{
    id := int(1)
    if (i!= nr_clients){
      secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
      x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
      R_i,_ := rand.Int(rand.Reader, phiN)
      c_i := InitClient(id,x_i,t,R_i,Zp)
      R_is = ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)).Num
      clients = append(clients, c_i)
    }else{
      secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
      x_n := new(big.Int).Add(secretSeed, big.NewInt(min))
      R_n := InitIntegerModP(new(big.Int).Sub( new(big.Int).Mul(phiN, new(big.Int).Div( R_is, phiN ) ) ,  R_is),Zp).Num
      c_n := InitClient(id, x_n, t, R_n, Zp)
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
  var range_proofs []*Proof

  for i:=int64(0); i< nr_clients ; i++{
    id := int(i)
    c_i := clients[id]
    var shares []*big.Int = Generate_shares(c_i.Id, c_i.Xi, t, nr_servers, Zp)
    rp,_ := Prove(c_i.Xi, c_i.Ri, &zkrp)
    range_proofs = append(range_proofs,rp)
    var tau *bn256.G2 = new(bn256.G2).ScalarBaseMult(big.NewInt(1))
    ga := new(bn256.G2).ScalarBaseMult(big.NewInt(zkrp.P.A))
    tau.Add(rp.P2.C,ga)
    tau_is = append(tau_is,tau)
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

  ok := Verify(tau_is, nr_clients, sigma, y, range_proofs, zkrp)
  sum := InitIntegerModP(y,Zp).Num
  if ok {
    fmt.Println("Verify ok")
      fmt.Println("Sum is", sum)
    }else{
      fmt.Println("BPRP error")
    }

}

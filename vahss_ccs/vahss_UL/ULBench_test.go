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
  "testing"
  ."hannaekthesis/ccs08"
  ."hannaekthesis/vahss"
  "crypto/rand"
  "github.com/ing-bank/zkrp/crypto/bn256"
)

func BenchmarkGenerateShares(b *testing.B){
  const nr_servers = int64(5)
  var prime, phiN *big.Int
  var Zp *Modular
  var t, min, max, length_interval int64
  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value,0)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2 // must be less than nr_servers-1
  min, max = 18, 200
  length_interval = max - min

  // Initiate cient
  secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
  x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
  R_i,_ := rand.Int(rand.Reader, phiN )
  c_i := InitClient(1,x_i,t,R_i,Zp) /// include vahss?

  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _ = Generate_shares(c_i.Id,c_i.Xi,t,nr_servers,Zp)
  }
}

func BenchmarkGenerateRangeProof(b *testing.B){
  var min, max,length_interval int64
  var prime, phiN *big.Int
  var zkrp Ccs08

  min, max = 18, 200
  length_interval = max - min
  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value,0)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))

  zkrp.Setup(min,max)


  secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
  x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
  R_i,_ := rand.Int(rand.Reader, phiN )

  b.ResetTimer()
  for n:=0; n<b.N ; n++{
  _ , _ = Prove(x_i, R_i, &zkrp)
  }
}

func BenchmarkPartialEval(b *testing.B){

  const nr_clients = int64(100)
  const nr_servers = int64(5)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t int64
  var clients []*Client
  var servers []*Server

  //prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)

  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value, 10)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2
  R_is = big.NewInt(0)

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



  for i:=int64(0); i< nr_clients ; i++{
    id := int(i)
    c_i := clients[id]
    var shares []*big.Int = Generate_shares(c_i.Id, c_i.Xi, t, nr_servers, Zp)
    //var rp ccs08.ProofSet
    for j:= int64(1); j<=nr_servers ; j++{
      Set_share(servers[j-1], c_i.Id, shares[j-1] )
    }
  }
  s_j := servers[1]
  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _ = Partial_eval(s_j, s_j.Id, s_j.Shares, nr_clients, Zp)
  }
}

func BenchmarkPartialProof(b *testing.B){

  const nr_clients = int64(100)
  const nr_servers = int64(5)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t int64
  var clients []*Client
  var servers []*Server

  //prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)

  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value, 10)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2
  R_is = big.NewInt(0)

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

  for i:=int64(0); i< nr_clients ; i++{
    id := int(i)
    c_i := clients[id]
    var shares []*big.Int = Generate_shares(c_i.Id, c_i.Xi, t, nr_servers, Zp)

    for j:= int64(1); j<=nr_servers ; j++{
      Set_share(servers[j-1], c_i.Id, shares[j-1] )
    }
  }

  s_j := servers[1]
  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _ = Partial_proof(s_j, s_j.Shares,nr_clients, Zp)
  }

}


func BenchmarkFinalEval(b *testing.B){
  const nr_clients = int64(100)
  const nr_servers = int64(5)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t int64
  var clients []*Client
  var servers []*Server

  //prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)

  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value, 10)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2
  R_is = big.NewInt(0)

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

  for i:=int64(0); i< nr_clients ; i++{
    id := int(i)
    c_i := clients[id]
    var shares []*big.Int = Generate_shares(c_i.Id, c_i.Xi, t, nr_servers, Zp)

    for j:= int64(1); j<=nr_servers ; j++{
      Set_share(servers[j-1], c_i.Id, shares[j-1] )
    }
  }


  var y_js []*big.Int

  for j:= int64(0); j < nr_servers ; j++{
    s_j := servers[j]
    y_js = append(y_js, Partial_eval(s_j, s_j.Id, s_j.Shares, nr_clients, Zp))
  }
  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _ = Final_eval(y_js, nr_servers, Zp)
  }
}

func BenchmarkFinalProof(b *testing.B){

  const nr_clients = int64(100)
  const nr_servers = int64(5)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t int64
  var clients []*Client
  var servers []*Server

  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value, 10)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2
  R_is = big.NewInt(0)

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


  for i:=int64(0); i< nr_clients ; i++{
    id := int(i)
    c_i := clients[id]
    var shares []*big.Int = Generate_shares(c_i.Id, c_i.Xi, t, nr_servers, Zp)



    for j:= int64(1); j<=nr_servers ; j++{
      Set_share(servers[j-1], c_i.Id, shares[j-1] )
    }
  }

  var sigma_js []*bn256.G2

  for j:= int64(0); j < nr_servers ; j++{
    s_j := servers[j]
    sigma_js = append(sigma_js,Partial_proof(s_j, s_j.Shares,nr_clients, Zp))
  }
  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _= Final_proof(sigma_js, nr_servers)
  }

}
func BenchmarkVerify_Servers(b *testing.B){
  const nr_clients = int64(100)
  const nr_servers = int64(5)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t,min ,max int64
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
  R_is = big.NewInt(0)

  zkrp.Setup(min,max)
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
  var range_proofs []*Proof

  for i:=int64(0); i< nr_clients ; i++{
    id := int(i)
    c_i := clients[id]
    var shares []*big.Int = Generate_shares(c_i.Id, c_i.Xi, t, nr_servers, Zp)
    //var rp ccs08.ProofSet
    rp , _ := Prove(c_i.Xi,c_i.Ri, &zkrp)
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
  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _ = Verify_Servers(tau_is, nr_clients, sigma, y)
  }
}

func BenchmarkVerify_RP(b *testing.B){
  var prime, phiN *big.Int
  var Zp *Modular
  var t, min, max int64
  var zkrp Ccs08

  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value, 10)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2
  min, max = 18,200
  zkrp.Setup(min,max)

  // Initiate cients
  id := int(1)
  secret := big.NewInt(12)
  R_i,_ := rand.Int(rand.Reader, phiN)
  c_i := InitClient(id,secret,t,R_i,Zp)
  rp ,_:= Prove(c_i.Xi, c_i.Ri, &zkrp)


  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _,_ = Verify_range(rp,&zkrp)
  }

}

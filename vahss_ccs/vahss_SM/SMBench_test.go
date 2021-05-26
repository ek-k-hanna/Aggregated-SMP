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
  "testing"
  "hannaekthesis/ccs08"
  ."hannaekthesis/vahss"
  "crypto/rand"
  "github.com/ing-bank/zkrp/crypto/bn256"
  "fmt"
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
  var x_i int64
  var prime, phiN *big.Int
  var set ccs08.ParamsSet

  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value,0)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  set = Generate_set()
//  secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
  x_i = 12
  R_i,_ := rand.Int(rand.Reader, phiN )

  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _ , _ = ccs08.ProveSet(x_i,R_i, set)
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
  var t int64
  var clients []*Client
  var servers []*Server
  var set ccs08.ParamsSet

  //prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)

  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value, 10)
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
  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _,_ = Verify_Servers(tau_is, nr_clients, sigma, y)
  }
}
func BenchmarkDifferentClients(b *testing.B) {
	var dif_nr_clients [1]int64 = [1]int64{100}//, 25, 50, 75, 100, 125, 150}
  var set ccs08.ParamsSet = Generate_set()
	for _, nr_clints := range dif_nr_clients {
		b.Run(fmt.Sprintf("New loop"), func(b *testing.B) {
			VBenchmarkVerify_RP(nr_clints,set,b)
		})
	}
}
func VBenchmarkVerify_RP(nr_clients int64, set ccs08.ParamsSet,b *testing.B){
  /*
  var prime, phiN *big.Int
  var Zp *Modular
  var t int64
  var set ccs08.ParamsSet

  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value, 10)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2
  set = Generate_set()

  // Initiate cients
  id := int(1)
  secret := big.NewInt(12)
  R_i,_ := rand.Int(rand.Reader, phiN)
  c_i := InitClient(id,secret,t,R_i,Zp)
  rp ,_:= ccs08.ProveSet(c_i.Xi.Int64(),c_i.Ri, set)


  b.ResetTimer()
    for n:=0; n<b.N ; n++{
      _,_ = ccs08.VerifySet(&rp,&set)
  }
*/





  //const nr_clients = int64(100)
  const nr_servers = int64(5)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t int64
  var clients []*Client
  var servers []*Server
  //var set ccs08.ParamsSet

  //prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)

  value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
  prime, _ = new(big.Int).SetString(value, 10)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2
  R_is = big.NewInt(0)
  //fmt.Println("Starting set")
  //set = Generate_set()
  //fmt.Println("Finished set")

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
  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _ = Verify_RP(range_proofs,set)
  }

}

func BenchmarkAggregation(b *testing.B){
    const nr_clients = int64(100)
    const nr_servers = int64(5)

    var prime, R_is, phiN *big.Int
    var Zp *Modular
    var t int64
    var clients []*Client
    var servers []*Server
    var set ccs08.ParamsSet

    value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
    prime, _ = new(big.Int).SetString(value, 10)
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
      rp , _ := ccs08.ProveSet(c_i.Xi.Int64(),c_i.Ri, set)
      range_proofs = append(range_proofs,rp)
      var tau_i *bn256.G2
      tau_i = rp.C
      tau_is = append(tau_is, tau_i)
    }

    challanges, prodChallanges := GetAllChallanges(range_proofs)
    b.ResetTimer()
    for n:=0; n<b.N ; n++{
      _ = ccs08.Aggregation(range_proofs, challanges, prodChallanges, &set)
    }

}

func BenchmarkVerifyAggregated(b *testing.B){

    const nr_clients = int64(100)
    const nr_servers = int64(5)

    var nr_aggregators int64 = int64(1)

    var prime, R_is, phiN *big.Int
    var Zp *Modular
    var t int64
    var clients []*Client
    var servers []*Server
    var set ccs08.ParamsSet

    //prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)

    value := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
    prime, _ = new(big.Int).SetString(value, 10)
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
      rp , _ := ccs08.ProveSet(c_i.Xi.Int64(),c_i.Ri, set)
      range_proofs = append(range_proofs,rp)
      var tau_i *bn256.G2
      tau_i = rp.C
      tau_is = append(tau_is, tau_i)
    }



  var tau_agg []*bn256.G2
  var agg []ccs08.ProofSet
  var chall_agg []*big.Int
  k := nr_clients/nr_aggregators

  for i:=int64(0); i < nr_aggregators; i++{
      low := i*k
      high := (i+1)*k
      RP_subset := range_proofs[low:high]
      tau_subset := tau_is[low:high]
      var prodCommits *bn256.G2 = new(bn256.G2)
      prodCommits.SetInfinity()
      for j:=0; j< len(tau_subset); j++{
        prodCommits.Add(prodCommits,tau_subset[j])
      }
      tau_agg =append(tau_agg,prodCommits)

      challanges, prodChallanges := GetAllChallanges(RP_subset)
      chall_agg = append(chall_agg,prodChallanges)
      aggregatedProof := ccs08.Aggregation(RP_subset, challanges, prodChallanges, &set)
      agg = append(agg,aggregatedProof)

    }

    fmt.Println(len(range_proofs), len(agg), len(chall_agg),len(tau_agg))
    b.ResetTimer()
    for n:=0; n<b.N ; n++{
      for i:= int64(0); i< nr_aggregators ; i++{
        low := i*k
        high := (i+1)*k
        _ = ccs08.VerifyAggregatedSet(range_proofs[low:high], &set, &agg[i], chall_agg[i], tau_agg[i])
      }
    }


}

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

package vahss_bprp

import (
	"crypto/rand"
	"hannaekthesis/bulletproof"
	. "hannaekthesis/vahss"
	"math/big"
	"testing"
	//"github.com/ing-bank/zkrp/crypto/p256"
	"fmt"
  "hannaekthesis/p256"
)


func BenchmarkGenerateShares(b *testing.B){
  const nr_servers = int64(5)
  var prime, phiN *big.Int
  var Zp *Modular
  var t, min, max, length_interval int64

  prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
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

  min, max = 18, 200
  length_interval = max - min
  prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))

  params, _ := bulletproofs.SetupGeneric(min,max)

  secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
  x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
  R_i,_ := rand.Int(rand.Reader, phiN )
  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _ = Generate_Bulletproof(x_i,R_i, params)
  }
}

func BenchmarkPartialEval(b *testing.B){
  const nr_clients = int64(100)
  const nr_servers = int64(5)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t, min, max, length_interval int64
  var clients []*Client
  var servers []*Server
	prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
	phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2 // must be less than nr_servers-1
	min, max = 18, 200
	length_interval = max - min
  R_is = big.NewInt(0)

  // Initiate cients
  for i:=int64(1) ; i<= nr_clients ; i++{
    id := int(1)
    if (i!= nr_clients){
			secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
      x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
      R_i,_ := rand.Int(rand.Reader, phiN )
      c_i := InitClient(id,x_i,t,R_i,Zp) /// include vahss?
			clients = append(clients, c_i)
      R_is = GetNumber(ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)))
    }else{
			secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
      x_n := new(big.Int).Add(secretSeed, big.NewInt(min))
      R_n := new(big.Int).Sub( new(big.Int).Mul(phiN, new(big.Int).Div( R_is, phiN ) ) ,  R_is)
      c_n := InitClient(id,x_n,t,R_n,Zp) /// include vahss?
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
    var shares []*big.Int = Generate_shares(c_i.Id,c_i.Xi,t,nr_servers,Zp)
    for j:= int64(1); j<=nr_servers ; j++{
      Set_share(servers[j-1], GetClientId(c_i), shares[j-1] )
    }
  }

  s_j := servers[1]

  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _ = Partial_eval(s_j, GetServerId(s_j), GetShares(s_j), nr_clients, Zp)
  }
}

func BenchmarkPartialProof(b *testing.B){

  const nr_clients = int64(100)
  const nr_servers = int64(5)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t, min, max, length_interval int64
  var clients []*Client
  var servers []*Server
	prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
	phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2 // must be less than nr_servers-1
	min, max = 18, 200
	length_interval = max - min
  R_is = big.NewInt(0)

	//bulletproofs setup params
	params, _ := bulletproofs.SetupGeneric(min,max)

  // Initiate cients
  for i:=int64(1) ; i<= nr_clients ; i++{
    id := int(1)
    if (i!= nr_clients){
			secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
      x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
      R_i,_ := rand.Int(rand.Reader, phiN )
      c_i := InitClient(id,x_i,t,R_i,Zp) /// include vahss?
			clients = append(clients, c_i)
      R_is = GetNumber(ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)))
    }else{
			secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
      x_n := new(big.Int).Add(secretSeed, big.NewInt(min))
      R_n := new(big.Int).Sub( new(big.Int).Mul(phiN, new(big.Int).Div( R_is, phiN ) ) ,  R_is)
      c_n := InitClient(id,x_n,t,R_n,Zp) /// include vahss?
      clients = append(clients, c_n)
    }
  }

  // Initiate Servers
  for j:= int64(1); j<= nr_servers ; j++{
    id := int(j)
    s_j := InitServer(id,Zp)
    servers = append(servers,s_j)
  }
  var tau_is []*p256.P256
	var range_proofs []bulletproofs.ProofBPRP

  for i:=int64(0); i< nr_clients ; i++{
    id := int(i)
    c_i := clients[id]
    var shares []*big.Int = Generate_shares(c_i.Id,c_i.Xi,t,nr_servers,Zp)
		var rp bulletproofs.ProofBPRP = Generate_Bulletproof(c_i.Xi,c_i.Ri, params)
		range_proofs = append(range_proofs,rp)

		// tau_i  = g^x*h^R, but commitment is g^(x-A)*h^R fix and store tau
		var tau_i *p256.P256 = new(p256.P256).ScalarBaseMult(big.NewInt(1))
		gA := new(p256.P256).ScalarBaseMult( new(big.Int).SetInt64(params.A) )
		tau_i.Multiply(rp.P2.V,gA)
		tau_is = append(tau_is,tau_i)

    for j:= int64(1); j<=nr_servers ; j++{
      Set_share(servers[j-1], GetClientId(c_i), shares[j-1] )
    }
  }
    s_j := servers[1]

  b.ResetTimer()
  for n:=0; n<b.N ; n++{
    _ = Partial_proof(s_j, GetShares(s_j), nr_clients, Zp)
  }

}


func BenchmarkFinalEval(b *testing.B){
  const nr_clients = int64(100)
  const nr_servers = int64(5)

  var prime, R_is, phiN *big.Int
  var Zp *Modular
  var t, min, max, length_interval int64
  var clients []*Client
  var servers []*Server
  prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
  phiN = new(big.Int).Sub(prime,big.NewInt(1))
  Zp = InitModular(prime)
  t = 2 // must be less than nr_servers-1
  min, max = 18, 200
  length_interval = max - min
  R_is = big.NewInt(0)

  // Initiate cients
  for i:=int64(1) ; i<= nr_clients ; i++{
    id := int(1)
    if (i!= nr_clients){
      secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
      x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
      R_i,_ := rand.Int(rand.Reader, phiN )
      c_i := InitClient(id,x_i,t,R_i,Zp) /// include vahss?
      clients = append(clients, c_i)
      R_is = GetNumber(ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)))
    }else{
      secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
      x_n := new(big.Int).Add(secretSeed, big.NewInt(min))
      R_n := new(big.Int).Sub( new(big.Int).Mul(phiN, new(big.Int).Div( R_is, phiN ) ) ,  R_is)
      c_n := InitClient(id,x_n,t,R_n,Zp) /// include vahss?
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
    var shares []*big.Int = Generate_shares(c_i.Id,c_i.Xi,t,nr_servers,Zp)

    for j:= int64(1); j<=nr_servers ; j++{
      Set_share(servers[j-1], GetClientId(c_i), shares[j-1] )
    }
  }
  var y_js []*big.Int

  for j:= int64(0); j < nr_servers ; j++{
    s_j := servers[j]
    y_js = append(y_js, Partial_eval(s_j, GetServerId(s_j), GetShares(s_j), nr_clients, Zp))
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
	  var t, min, max, length_interval int64
	  var clients []*Client
	  var servers []*Server
		prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
		phiN = new(big.Int).Sub(prime,big.NewInt(1))
	  Zp = InitModular(prime)
	  t = 2 // must be less than nr_servers-1
		min, max = 18, 200
		length_interval = max - min
	  R_is = big.NewInt(0)

		//bulletproofs setup params
		params, _ := bulletproofs.SetupGeneric(min,max)

	  // Initiate cients
	  for i:=int64(1) ; i<= nr_clients ; i++{
	    id := int(1)
	    if (i!= nr_clients){
				secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
	      x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
	      R_i,_ := rand.Int(rand.Reader, phiN )
	      c_i := InitClient(id,x_i,t,R_i,Zp) /// include vahss?
				clients = append(clients, c_i)
	      R_is = GetNumber(ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)))
	    }else{
				secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
	      x_n := new(big.Int).Add(secretSeed, big.NewInt(min))
	      R_n := new(big.Int).Sub( new(big.Int).Mul(phiN, new(big.Int).Div( R_is, phiN ) ) ,  R_is)
	      c_n := InitClient(id,x_n,t,R_n,Zp) /// include vahss?
	      clients = append(clients, c_n)
	    }
	  }

	  // Initiate Servers
	  for j:= int64(1); j<= nr_servers ; j++{
	    id := int(j)
	    s_j := InitServer(id,Zp)
	    servers = append(servers,s_j)
	  }

	  var tau_is []*p256.P256
		var range_proofs []bulletproofs.ProofBPRP

	  for i:=int64(0); i< nr_clients ; i++{
	    id := int(i)
	    c_i := clients[id]
	    var shares []*big.Int = Generate_shares(c_i.Id,c_i.Xi,t,nr_servers,Zp)
			var rp bulletproofs.ProofBPRP = Generate_Bulletproof(c_i.Xi,c_i.Ri, params)
			range_proofs = append(range_proofs,rp)

			// tau_i  = g^x*h^R, but commitment is g^(x-A)*h^R fix and store tau
			var tau_i *p256.P256 = new(p256.P256).ScalarBaseMult(big.NewInt(1))
			gA := new(p256.P256).ScalarBaseMult( new(big.Int).SetInt64(params.A) )
			tau_i.Multiply(rp.P2.V,gA)
			tau_is = append(tau_is,tau_i)

	    for j:= int64(1); j<=nr_servers ; j++{
	      Set_share(servers[j-1], GetClientId(c_i), shares[j-1] )
	    }
	  }
	  var sigma_js []*p256.P256

	  for j:= int64(0); j < nr_servers ; j++{
	    s_j := servers[j]
	    sigma_js = append(sigma_js,Partial_proof(s_j, GetShares(s_j), nr_clients, Zp))
	  }

    b.ResetTimer()
    for n:=0; n<b.N ; n++{
      _= Final_proof(sigma_js,nr_servers)
    }
  }


func BenchmarkDifferentClients(b *testing.B) {
	var dif_nr_clients [1]int64 = [1]int64{100}//, 25, 50, 75, 100, 125, 150}
	for _, nr_clints := range dif_nr_clients {
		b.Run(fmt.Sprintf("New loop"), func(b *testing.B) {
			VerifyBench_RP(nr_clints,b)
		})
	}
}
func VerifyBench_RP(nr_clients int64,b *testing.B) {
	//const nr_clients = int64(100)
	const nr_servers = int64(5)
  fmt.Print(nr_clients)

	var prime, R_is, phiN *big.Int
	var Zp *Modular
	var t, min, max, length_interval int64
	var clients []*Client
	var servers []*Server
	prime, _ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
	phiN = new(big.Int).Sub(prime, big.NewInt(1))
	Zp = InitModular(prime)
	t = 2 // must be less than nr_servers-1
	min, max = 18, 200
	length_interval = max - min
	R_is = big.NewInt(0)

	//bulletproofs setup params
	params, _ := bulletproofs.SetupGeneric(min, max)

	// Initiate cients
	for i := int64(1); i <= nr_clients; i++ {
		id := int(1)
		if i != nr_clients {
			secretSeed, _ := rand.Int(rand.Reader, big.NewInt(length_interval))
			x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
			R_i, _ := rand.Int(rand.Reader, phiN)
			c_i := InitClient(id, x_i, t, R_i, Zp) /// include vahss?
			clients = append(clients, c_i)
			R_is = GetNumber(ModAdd(InitIntegerModP(R_is, Zp), InitIntegerModP(R_i, Zp)))
		} else {
			secretSeed, _ := rand.Int(rand.Reader, big.NewInt(length_interval))
			x_n := new(big.Int).Add(secretSeed, big.NewInt(min))
			R_n := new(big.Int).Sub(new(big.Int).Mul(phiN, new(big.Int).Div(R_is, phiN)), R_is)
			c_n := InitClient(id, x_n, t, R_n, Zp) /// include vahss?
			clients = append(clients, c_n)
		}
	}

	// Initiate Servers
	for j := int64(1); j <= nr_servers; j++ {
		id := int(j)
		s_j := InitServer(id, Zp)
		servers = append(servers, s_j)
	}

	var range_proofs []bulletproofs.ProofBPRP

	for i := int64(0); i < nr_clients; i++ {
		id := int(i)
		c_i := clients[id]
		var rp bulletproofs.ProofBPRP = Generate_Bulletproof(c_i.Xi, c_i.Ri, params)
		range_proofs = append(range_proofs, rp)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = Verify_RP(range_proofs)
	}
}


  func BenchmarkVerify_RP(b *testing.B){
    const nr_clients = int64(100)
    const nr_servers = int64(5)

    var prime, R_is, phiN *big.Int
    var Zp *Modular
    var t, min, max, length_interval int64
    var clients []*Client
    var servers []*Server
    prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
    phiN = new(big.Int).Sub(prime,big.NewInt(1))
    Zp = InitModular(prime)
    t = 2 // must be less than nr_servers-1
    min, max = 18, 200
    length_interval = max - min
    R_is = big.NewInt(0)

    //bulletproofs setup params
    params, _ := bulletproofs.SetupGeneric(min,max)

    // Initiate cients
    for i:=int64(1) ; i<= nr_clients ; i++{
      id := int(1)
      if (i!= nr_clients){
        secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
        x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
        R_i,_ := rand.Int(rand.Reader, phiN )
        c_i := InitClient(id,x_i,t,R_i,Zp) /// include vahss?
        clients = append(clients, c_i)
        R_is = GetNumber(ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)))
      }else{
        secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
        x_n := new(big.Int).Add(secretSeed, big.NewInt(min))
        R_n := new(big.Int).Sub( new(big.Int).Mul(phiN, new(big.Int).Div( R_is, phiN ) ) ,  R_is)
        c_n := InitClient(id,x_n,t,R_n,Zp) /// include vahss?
        clients = append(clients, c_n)
      }
    }

    // Initiate Servers
    for j:= int64(1); j<= nr_servers ; j++{
      id := int(j)
      s_j := InitServer(id,Zp)
      servers = append(servers,s_j)
    }


    var range_proofs []bulletproofs.ProofBPRP

    for i:=int64(0); i< nr_clients ; i++{
      id := int(i)
      c_i := clients[id]
      var rp bulletproofs.ProofBPRP = Generate_Bulletproof(c_i.Xi,c_i.Ri, params)
      range_proofs = append(range_proofs,rp)
    }

    b.ResetTimer()
    for n:=0; n<b.N ; n++{
      _ = Verify_RP(range_proofs)
    }
  }


  func BenchmarkVerify_Servers(b *testing.B){
    const nr_clients = int64(100)
    const nr_servers = int64(5)

    var prime, R_is, phiN *big.Int
    var Zp *Modular
    var t, min, max, length_interval int64
    var clients []*Client
    var servers []*Server
    prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
    phiN = new(big.Int).Sub(prime,big.NewInt(1))
    Zp = InitModular(prime)
    t = 2 // must be less than nr_servers-1
    min, max = 18, 200
    length_interval = max - min
    R_is = big.NewInt(0)

    //bulletproofs setup params
    params, _ := bulletproofs.SetupGeneric(min,max)

    // Initiate cients
    for i:=int64(1) ; i<= nr_clients ; i++{
      id := int(1)
      if (i!= nr_clients){
        secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
        x_i := new(big.Int).Add(secretSeed, big.NewInt(min))
        R_i,_ := rand.Int(rand.Reader, phiN )
        c_i := InitClient(id,x_i,t,R_i,Zp) /// include vahss?
        clients = append(clients, c_i)
        R_is = GetNumber(ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)))
      }else{
        secretSeed,_ := rand.Int(rand.Reader, big.NewInt(length_interval) )
        x_n := new(big.Int).Add(secretSeed, big.NewInt(min))
        R_n := new(big.Int).Sub( new(big.Int).Mul(phiN, new(big.Int).Div( R_is, phiN ) ) ,  R_is)
        c_n := InitClient(id,x_n,t,R_n,Zp) /// include vahss?
        clients = append(clients, c_n)
      }
    }

    // Initiate Servers
    for j:= int64(1); j<= nr_servers ; j++{
      id := int(j)
      s_j := InitServer(id,Zp)
      servers = append(servers,s_j)
    }


    var tau_is []*p256.P256
    var range_proofs []bulletproofs.ProofBPRP

    for i:=int64(0); i< nr_clients ; i++{
      id := int(i)
      c_i := clients[id]
      var shares []*big.Int = Generate_shares(c_i.Id,c_i.Xi,t,nr_servers,Zp)
      var rp bulletproofs.ProofBPRP = Generate_Bulletproof(c_i.Xi,c_i.Ri, params)
      range_proofs = append(range_proofs,rp)

      // tau_i  = g^x*h^R, but commitment is g^(x-A)*h^R fix and store tau
      var tau_i *p256.P256 = new(p256.P256).ScalarBaseMult(big.NewInt(1))
      gA := new(p256.P256).ScalarBaseMult( new(big.Int).SetInt64(params.A) )
      tau_i.Multiply(rp.P2.V,gA)
      tau_is = append(tau_is,tau_i)

      for j:= int64(1); j<=nr_servers ; j++{
        Set_share(servers[j-1], GetClientId(c_i), shares[j-1] )
      }
    }
    var y_js []*big.Int
    var sigma_js []*p256.P256

    for j:= int64(0); j < nr_servers ; j++{
      s_j := servers[j]
      y_js = append(y_js, Partial_eval(s_j, GetServerId(s_j), GetShares(s_j), nr_clients, Zp))
      sigma_js = append(sigma_js,Partial_proof(s_j, GetShares(s_j), nr_clients, Zp))
    }

    y := Final_eval(y_js, nr_servers, Zp)
    sigma := Final_proof(sigma_js,nr_servers)

    b.ResetTimer()
    for n:=0; n<b.N ; n++{
      _ = Verify_Servers(tau_is, nr_clients, sigma, y)
    }
  }

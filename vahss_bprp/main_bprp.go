package vahss_bprp

import (
		"fmt"
		"hannaekthesis/bulletproof"
		"github.com/ing-bank/zkrp/crypto/p256"
		"math/big"
		"crypto/rand"
		."hannaekthesis/vahss"
)

func Generate_Bulletproof(secret,randomValue *big.Int, params *bulletproofs.Bprp) (bulletproofs.ProofBPRP){
  proof,_:= bulletproofs.ProveGeneric(secret,randomValue, params)
  return proof
}

func Verify_Bulletproofs(proofs []bulletproofs.ProofBPRP)(bool){
	var ok bool = true
	for _,proof := range proofs{
		ok, _ = proof.Verify()
	}
  return ok
}

func Main_bprp() {

	  const nr_clients = int64(100)
	  const nr_servers = int64(5)

	  var prime, R_is, phiN *big.Int
		var g *p256.P256
	  var Zp *Modular
	  var t int64
	  var clients []*Client
	  var servers []*Server
		prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
		phiN = new(big.Int).Sub(prime,big.NewInt(1))
	  Zp = InitModular(prime)
	  t = 2 // must be less than nr_servers-1
	  R_is = big.NewInt(0)

		//bulletproofs setup params
		params, _ := bulletproofs.SetupGeneric(0, 200)
		g = params.BP1.G

	  // Initiate cients
	  for i:=int64(1) ; i<= nr_clients ; i++{
	    id := int(1)
	    if (i!= nr_clients){
				secret,_ := rand.Int(rand.Reader, big.NewInt(200) )//new(big.Int).SetInt64(int64(40))
	      R_i,_ := rand.Int(rand.Reader, phiN )
	      c_i := InitClient(id,secret,t,R_i,Zp) /// include vahss?
	      R_is = GetNumber(ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)))
	      clients = append(clients, c_i)
	    }else{
				secret,_ := rand.Int(rand.Reader, big.NewInt(200) )
	      R_n := new(big.Int).Sub( new(big.Int).Mul(phiN, new(big.Int).Div( R_is, phiN ) ) ,  R_is)
	      c_n := InitClient(id,secret,t,R_n,Zp) /// include vahss?
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
	    sigma_js = append(sigma_js,Partial_proof(s_j, GetShares(s_j), g, nr_clients, Zp))
	  }

	  y := Final_eval(y_js, nr_servers, Zp)
	  sigma := Final_proof(sigma_js,nr_servers, g)

	  result_vahss := Verify(tau_is, nr_clients, sigma, y, g, Zp)
		result_RP := Verify_Bulletproofs(range_proofs)
	  sum := GetNumber(InitIntegerModP(y,Zp))


	  if result_vahss {
			if result_RP {
	    fmt.Println("Verify ok")
			}else{
				fmt.Println("BPRP error")
			}
	  }else{
	    fmt.Println("Verification ERROR")
	  }
	  fmt.Println("Sum is", sum)

}



/*
func main() {

	  const nr_clients = int64(2)
	  const nr_servers = int64(4)

	  var prime, g, R_is *big.Int
	  var Zp *vahss.Modular
	  var t int64
	  var clients []*vahss.Client
	  var servers []*vahss.Server

	  prime = big.NewInt(3911)//Get_large_enough_prime([]int64{int64(200000000),int64(40478437826842)})
	  fmt.Println("Found prime",prime)
	  Zp = vahss.InitModular(prime)
	  g = big.NewInt(3)
	  t = 4
	  R_is = big.NewInt(0)

		//bulletproofs setup params
		params, _ := bulletproofs.SetupGeneric(2, 200)

		secret := new(big.Int).SetInt64(int64(40))
	  // Initiate cients
	  for i:=int64(1) ; i<= nr_clients ; i++{
	    id := int(1)
	    if (i!= nr_clients){
	      R_i := big.NewInt(0) //rand.Int(rand.Reader, big.NewInt(1) )
	      c_i := vahss.InitClient(id,secret,t,g,R_i,Zp) /// include vahss?
	      R_is = vahss.GetNumber(vahss.ModAdd(vahss.InitIntegerModP(R_is,Zp),vahss.InitIntegerModP(R_i,Zp)))
	      clients = append(clients, c_i)
	    }else{
	       R_i := big.NewInt(0)
	      //R_i := math.Ceil( R_is / (q-1) ) * (q-1) - R_is
	      c_i := vahss.InitClient(id,secret,t,g,R_i,Zp) /// include vahss?
	      clients = append(clients, c_i)
	    }
	  }


	  // Initiate Servers
	  for j:= int64(1); j<= nr_servers ; j++{
	    id := int(j)
	    s_j := vahss.InitServer(id,Zp)
	    servers = append(servers,s_j)
	  }



	  var tau_is []*big.Int
		var range_proofs []bulletproofs.ProofBPRP

	  for i:=int64(0); i< nr_clients ; i++{
	    id := int(i)
	    c_i := clients[id]
	    var shares []*big.Int = vahss.Generate_shares(c_i,nr_servers)
	    tau_is = append(tau_is,vahss.Get_tau(c_i))
			Ri_RP,_ := rand.Int(rand.Reader,prime)
			var rp bulletproofs.ProofBPRP = Generate_Bulletproof(secret, Ri_RP ,params)
			range_proofs = append(range_proofs,rp)
	    for j:= int64(1); j<=nr_servers ; j++{
	      vahss.Set_share(servers[j-1], vahss.GetClientId(c_i), shares[j-1] )
	    }
	  }
	  var y_js []*big.Int
	  var sigma_js []*big.Int

	  for j:= int64(0); j < nr_servers ; j++{
	    s_j := servers[j]
	    y_js = append(y_js, vahss.Partial_eval(s_j, vahss.GetServerId(s_j), vahss.GetShares(s_j), nr_clients, Zp))
	    sigma_js = append(sigma_js,vahss.Partial_proof(s_j, vahss.GetShares(s_j), g, nr_clients, Zp))
	  }

	  y := vahss.Final_eval(y_js, nr_servers, Zp)
	  sigma := vahss.Final_proof(y,sigma_js,nr_servers, g, Zp)

	  result_vahss := vahss.Verify(tau_is, nr_clients, sigma, y, g, Zp)
		result_RP := Verify_Bulletproofs(range_proofs)
	  sum := vahss.GetNumber(vahss.InitIntegerModP(y,Zp))
	  if result_vahss {
			if result_RP {
	    fmt.Println("Verify ok")
			}else{
				fmt.Println("BPRP error")
			}
	  }else{
	    fmt.Println("Verification ERROR")
	  }
	  fmt.Println("Sum is", sum)

}
*/

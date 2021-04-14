package vahss_bprp

import (
  "math/big"
  "github.com/ing-bank/zkrp/crypto/p256"
  ."hannaekthesis/vahss"
  "hannaekthesis/bulletproof"
  "fmt"
)

func Generate_shares(i int, x_i *big.Int, degree int64 , nr_servers int64, mod *Modular) ([]*big.Int){
  shares_Zp := gen_secret_share_additive_with_hash_functions(i, x_i, degree,nr_servers,mod )
  var shares []*big.Int
  for _,share := range shares_Zp {
    shares = append(shares,share.Num)
  }
  return shares
}
func gen_secret_share_additive_with_hash_functions(i int, x_i *big.Int, degree int64 , nr_servers int64, mod *Modular) ([]*IntegerModP){
  //i: index of the client
  //x_i: secret input of the client i
  //t: threshold (t + 1 reconstruct)
  //nr_servers: Number of servers
  var shares []*IntegerModP = Generate_input_shamir_secret(x_i, degree, nr_servers, mod.P) //[]*IntegerModP
  return shares//, tau_i
}

func Generate_Bulletproof(secret,randomValue *big.Int, params *bulletproofs.Bprp) (bulletproofs.ProofBPRP){
  proof,_:= bulletproofs.ProveGeneric(secret,randomValue, params)
  return proof
}


func Partial_eval(s_j *Server, j int, shares []*big.Int, nr_clients int64, mod *Modular) (*big.Int){
  var partialeval *IntegerModP =  InitIntegerModP(big.NewInt(0),mod)
  for i := int64(0); i < nr_clients; i++ {
      partialeval = ModAdd(partialeval, InitIntegerModP(shares[i],mod))
    }
  return partialeval.Num
}

func Final_eval(partialEvals []*big.Int, nr_servers int64, mod *Modular) (*big.Int){
  finaleval := big.NewInt(0)
  for j := int64(0); j < nr_servers; j++ {
      finaleval = new(big.Int).Add(finaleval,partialEvals[j])
    }
  return finaleval
}
func Partial_proof(s_j *Server, shares []*big.Int, g *p256.P256, nr_clients int64, mod *Modular) *p256.P256{
  var y_j *big.Int = Partial_eval(s_j,s_j.Id, shares, nr_clients, mod)
  var sigma_j *p256.P256 = new(p256.P256).ScalarBaseMult(y_j)
  return sigma_j

}

func Final_proof(partialProofs []*p256.P256, nr_servers int64, g *p256.P256)(*p256.P256){
  var sigma *p256.P256 = partialProofs[0]
  for i:= int64(1) ; i< nr_servers ; i++{
    sigma.Multiply(sigma, partialProofs[i])
  }
  return sigma
}

func Verify_Servers(tau_is []*p256.P256, nr_clients int64, sigma *p256.P256, y *big.Int) (bool){
  var tau *p256.P256 = tau_is[0]
  for i:= int64(1) ; i< nr_clients ; i++{
    tau.Multiply(tau, tau_is[i])
  }
  sigma.Neg(sigma)
  tau.Multiply(tau,sigma)
  ok_1 := tau.IsZero()

  hash := new(p256.P256).ScalarBaseMult(y)
  hash.Multiply(hash,sigma)
  ok_2 :=  hash.IsZero()
  return ok_1 && ok_2
}

func Verify_RP(proofs []bulletproofs.ProofBPRP)(bool){
	var ok bool = true
	for _,proof := range proofs{
		ok, _ = proof.Verify()
	}
  return ok
}

func Verify(tau_is []*p256.P256, nr_clients int64, sigma *p256.P256, y *big.Int, RPs []bulletproofs.ProofBPRP) (bool){
  ok_servers := Verify_Servers(tau_is, nr_clients, sigma, y)
  fmt.Println("Servers honest: ", ok_servers)
  ok_clients :=  Verify_RP(RPs)
  fmt.Println("Cients honest: ", ok_clients)
  return ok_servers && ok_clients


}

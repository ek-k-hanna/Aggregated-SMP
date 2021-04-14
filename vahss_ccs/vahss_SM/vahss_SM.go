package vahss_SM

import (
  "math/big"
    "github.com/ing-bank/zkrp/crypto/bn256"
    ."hannaekthesis/vahss"
    ."hannaekthesis/ccs08"
    "bytes"
    "fmt"
)

func Generate_shares(i int, x_i *big.Int, degree int64 , nr_servers int64, mod *Modular) ([]*big.Int){
  var shares_Zp []*IntegerModP = Generate_input_shamir_secret(x_i, degree, nr_servers, mod.P) //[]*IntegerModP
  var shares []*big.Int
  for _,share := range shares_Zp {
    shares = append(shares,share.Num)
  }
  return shares
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

func Partial_proof(s_j *Server, shares []*big.Int, nr_clients int64, mod *Modular) *bn256.G2{
  var y_j *big.Int = Partial_eval(s_j,s_j.Id, shares, nr_clients, mod)
  var sigma_j *bn256.G2 = new(bn256.G2).ScalarBaseMult(y_j)
  return sigma_j

}

func Final_proof(partialProofs []*bn256.G2, nr_servers int64)(*bn256.G2){
  var sigma *bn256.G2 = partialProofs[0]
  for i:= int64(1) ; i< nr_servers ; i++{
    sigma.Add(sigma, partialProofs[i])
  }
  return sigma
}

func Verify_RP(proofs []ProofSet, set ParamsSet) (bool){
  var ok bool = true
  for _,proof := range proofs{
    ok, _ = VerifySet(&proof,&set)
  }
  return ok
}

func Verify_Servers(tau_is []*bn256.G2, nr_clients int64, sigma *bn256.G2, y *big.Int) (bool){
  var tau *bn256.G2 = tau_is[0]
  for i:= int64(1) ; i< nr_clients ; i++{
    tau.Add(tau,tau_is[i])
  }

  sigmaBytes := sigma.Marshal()
  tauBytes := tau.Marshal()
  ok_1 := bytes.Equal(sigmaBytes, tauBytes)

  hash := new(bn256.G2).ScalarBaseMult(y)
  hashBytes := hash.Marshal()
  ok_2 := bytes.Equal(sigmaBytes,hashBytes)

  return ok_1 && ok_2
}

func Verify(tau_is []*bn256.G2, nr_clients int64, sigma *bn256.G2, y *big.Int, RPs []ProofSet, set ParamsSet ) (bool){
  ok_servers := Verify_Servers(tau_is, nr_clients, sigma, y)
  fmt.Println("Servers honest: ", ok_servers)
  ok_clients :=  Verify_RP(RPs, set)
  fmt.Println("Cients honest: ", ok_clients)
  return ok_servers && ok_clients


}

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

import (
  "math/big"
    "github.com/ing-bank/zkrp/crypto/bn256"
    _"github.com/ing-bank/zkrp/util/bn"

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

func GetAllChallanges(proofs []ProofSet) ([]*big.Int, *big.Int){
  nbrProofs := len(proofs)
  var challanges []*big.Int
  var prod *big.Int = big.NewInt(1)
  for i:=0; i< nbrProofs; i++{
    challange := GetChallange(&proofs[i])
    challanges = append(challanges, challange)
    prod = new(big.Int).Mul(prod,challange)
  }
  return challanges,prod
}

func Verify_RP(proofs []ProofSet, set ParamsSet) (bool){
  var ok bool = true
  for _,proof := range proofs{
    ok,_ = VerifySet(&proof,&set)
  }
  return ok
}
func Verify_AggregatedRP(proofs []ProofSet, set ParamsSet, prodCommits *bn256.G2) (bool){
  challanges, prodChallanges := GetAllChallanges(proofs)
  aggregatedProof := Aggregation(proofs,challanges, prodChallanges, &set)
  ok := VerifyAggregatedSet(proofs,&set, &aggregatedProof,prodChallanges, prodCommits)
  return ok
}

func Verify_Servers(tau_is []*bn256.G2, nr_clients int64, sigma *bn256.G2, y *big.Int) (bool, *bn256.G2){
  var tau *bn256.G2 = new(bn256.G2)
  tau.SetInfinity()
  for i:= int64(0) ; i< nr_clients ; i++{
    tau.Add(tau,tau_is[i])
  }

  sigmaBytes := sigma.Marshal()
  tauBytes := tau.Marshal()
  ok_1 := bytes.Equal(sigmaBytes, tauBytes)
  hash := new(bn256.G2).ScalarBaseMult(y)
  hashBytes := hash.Marshal()
  ok_2 := bytes.Equal(sigmaBytes,hashBytes)

  return ok_1 && ok_2 , tau
}


func PartialVerify(tau_is []*bn256.G2, nr_clients int64, sigma *bn256.G2, y *big.Int, RPs []ProofSet, set ParamsSet ,nbr_aggreatingParties int64) (bool){
  var ok_servers bool
  var ok_clients bool = true
  ok_servers,_ = Verify_Servers(tau_is, nr_clients, sigma, y)
  fmt.Println("Servers honest: ", ok_servers)

  for i:=int64(0); i < nbr_aggreatingParties; i++{
    k := nr_clients/nbr_aggreatingParties
    low := i*k
    high := (i+1)*k
    RP_subset := RPs[low:high]
    tau_subset := tau_is[low:high]
    var prodCommits *bn256.G2 = new(bn256.G2)
    prodCommits.SetInfinity()
    for j:=0; j< len(tau_subset); j++{
      prodCommits.Add(prodCommits,tau_subset[j])
    }
    fmt.Println(len(RP_subset))
    ok_clients = (ok_clients && Verify_AggregatedRP(RP_subset,set,prodCommits) )
  }
  fmt.Println("Aggregated Clients honest", ok_clients)
  return ok_servers && ok_clients
}

func Verify(tau_is []*bn256.G2, nr_clients int64, sigma *bn256.G2, y *big.Int, RPs []ProofSet, set ParamsSet) (bool){
    ok_servers,_ := Verify_Servers(tau_is, nr_clients, sigma, y)
    fmt.Println("Servers honest: ", ok_servers)
    ok_clients :=  Verify_RP(RPs, set)
    fmt.Println("Cients honest: ", ok_clients)
    return (ok_servers && ok_clients)

}

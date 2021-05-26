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
  "math/big"
  "hannaekthesis/p256"//"github.com/ing-bank/zkrp/crypto/p256"
  ."hannaekthesis/vahss"
  ."hannaekthesis/bulletproof"
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

func Generate_Bulletproof(secret,randomValue *big.Int, params *Bprp) (ProofBPRP){
  proof,_:= ProveGeneric(secret,randomValue, params)
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
func Partial_proof(s_j *Server, shares []*big.Int, nr_clients int64, mod *Modular) *p256.P256{
  var y_j *big.Int = Partial_eval(s_j,s_j.Id, shares, nr_clients, mod)
  var sigma_j *p256.P256 = new(p256.P256).ScalarBaseMult(y_j)
  return sigma_j

}

func Final_proof(partialProofs []*p256.P256, nr_servers int64)(*p256.P256){
  var sigma *p256.P256 = partialProofs[0]
  for i:= int64(1) ; i< nr_servers ; i++{
    sigma.Multiply(sigma, partialProofs[i])
  }
  return sigma
}
/*
func AccumulateIP(acc InnerProductProof,ip InnerProductProof)(InnerProductProof)  {
  n := len(acc.Ls)
  for i:= 0; i< n ;i++{
    acc.Ls[i].Multiply(acc.Ls[i],ip.Ls[i])
    acc.Rs[i].Multiply(acc.Rs[i],ip.Rs[i])
  }
  acc.U.Multiply(acc.U,ip.U)
  acc.P.Multiply(acc.P,ip.P)

  acc.Gg.Multiply(acc.Gg,ip.Gg)
  acc.Hh.Multiply(acc.Hh,ip.Hh)

  acc.A = new(big.Int).Add(acc.A,ip.A)
  acc.B = new(big.Int).Add(acc.B,ip.B)

  //acc.Params.P.Multiply(acc.Params.P,bp.Params.P)

  return acc
}

func AccumulateRP(acc,bp BulletProof) (BulletProof) {
  acc.V.Multiply(acc.V,bp.V)
  acc.A.Multiply(acc.A,bp.A)
  acc.S.Multiply(acc.S,bp.S)
  acc.T1.Multiply(acc.T1,bp.T1)
  acc.T2.Multiply(acc.T2,bp.T2)

  acc.Taux = new(big.Int).Add(acc.Taux,bp.Taux)
  acc.Mu = new(big.Int).Add(acc.Mu,bp.Mu)
  acc.Tprime = new(big.Int).Add(acc.Tprime,bp.Tprime)

  acc.Commit.Multiply(acc.Commit,bp.Commit)

  acc.InnerProductProof = AccumulateIP(acc.InnerProductProof, bp.InnerProductProof)

  return acc
}
func Accumulate(proofs []ProofBPRP) (proof ProofBPRP){
  var accProof ProofBPRP
  accProof.P1 = proofs[0].P1
  accProof.P2 = proofs[0].P2

  for i := 1 ; i< len(proofs) ; i++{
    AccumulateRP(accProof.P1 ,proofs[i].P1)
    AccumulateRP(accProof.P2,proofs[i].P2)
  }
  return accProof
}
*/
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

func Verify_RP(proofs []ProofBPRP)(bool){
	var ok bool = true
	for _,proof := range proofs{
		ok, _ = proof.Verify()
	}
  return ok
}

func Verify(tau_is []*p256.P256, nr_clients int64, sigma *p256.P256, y *big.Int, RPs []ProofBPRP) (bool){
  ok_servers := Verify_Servers(tau_is, nr_clients, sigma, y)
  fmt.Println("Servers honest: ", ok_servers)
  ok_clients :=  Verify_RP(RPs)
  fmt.Println("Cients honest: ", ok_clients)
  return ok_servers && ok_clients


}

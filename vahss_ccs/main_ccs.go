package vahss_ccs
/*
import (
	"fmt"
  "crypto/rand"
  "github.com/ing-bank/zkrp/crypto/bn256"
  ."github.com/ing-bank/zkrp/util"
  "github.com/ing-bank/zkrp/util/bn"
  "github.com/ing-bank/zkrp/util/intconversion"
	//"hannaekthesis/bulletproof"
	//	"hannaekthesis/vahss"
	"hannaekthesis/ccs08"
		//. "github.com/ing-bank/zkrp/util"
		"math/big"
		"C"
		//"crypto/rand"
)
func Generate_set() (*paramSet){
    set := make([]int64,4)
    set[0] = 12
    set[1] = 42
    set[2] = 61
    set[3] = 71
    p, _ := SetupSet(set)
    return set
}
func main() {

	  const nr_clients = int64(10)
	  const nr_servers = int64(4)

	  var prime, R_is *big.Int
		var g *p256.P256
	  var Zp *Modular
	  var t int64
	  var clients []*Client
	  var servers []*Server
    var set []*int64 = Generate_set()
		prime,_ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
	  //prime = big.NewInt(115792089237316195423570985008687907852837564279074904382605163141518161494337)//Get_large_enough_prime([]int64{int64(200000000),int64(40478437826842)})
	  fmt.Println("Found prime",prime)
	  Zp = InitModular(prime)
	  //g = big.NewInt(3)
	  t = 4
	  R_is = big.NewInt(0)

		//bulletproofs setup params
		params, _ := bulletproofs.SetupGeneric(2, 200)
		g = params.BP1.G
		secret := new(big.Int).SetInt64(int64(40))
	  // Initiate cients
	  for i:=int64(1) ; i<= nr_clients ; i++{
	    id := int(1)
	    if (i!= nr_clients){
	      R_i := big.NewInt(0) //rand.Int(rand.Reader, big.NewInt(1) )
	      c_i := InitClient(id,secret,t,g,R_i,Zp) /// include vahss?
	      R_is = GetNumber(ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)))
	      clients = append(clients, c_i)
	    }else{
	       R_i := big.NewInt(0)
	      //R_i := math.Ceil( R_is / (q-1) ) * (q-1) - R_is
	      c_i := InitClient(id,secret,t,g,R_i,Zp) /// include vahss?
	      clients = append(clients, c_i)
	    }
	  }


	  // Initiate Servers
	  for j:= int64(1); j<= nr_servers ; j++{
	    id := int(j)
	    s_j := InitServer(id,Zp)
	    servers = append(servers,s_j)
	  }



	  var tau_is []*p256.P256
		//var range_proofs []bulletproofs.ProofBPRP

	  for i:=int64(0); i< nr_clients ; i++{
	    id := int(i)
	    c_i := clients[id]
	    var shares []*big.Int = Generate_shares(c_i,nr_servers)
	    tau_is = append(tau_is,Get_tau(c_i))
			//Ri_RP,_ := rand.Int(rand.Reader,prime)
			//var rp bulletproofs.ProofBPRP = Generate_Bulletproof(secret, Ri_RP ,params)
			//range_proofs = append(range_proofs,rp)
	    for j:= int64(1); j<=nr_servers ; j++{
	      Set_share(servers[j-1], GetClientId(c_i), shares[j-1] )
	    }
	  }
	  var y_js []*big.Int
	  var sigma_js []*p256.P256

	  for j:= int64(0); j < nr_servers ; j++{
	    s_j := servers[j]
	    y_js = append(y_js, Partial_eval(s_j, GetServerId(s_j), GetShares(s_j), nr_clients, Zp))
	    sigma_js = append(sigma_js,Partial_proof_test(s_j, GetShares(s_j), g, nr_clients, Zp))
	  }

	  y := Final_eval(y_js, nr_servers, Zp)
		sum := GetNumber(InitIntegerModP(y,Zp))
	  sigma := Final_proof_test(sigma_js,nr_servers, g)

	  result_vahss := Verify_test(tau_is, nr_clients, sigma, y, g, Zp)
		//result_RP := Verify_Bulletproofs(range_proofs)
	  //sum := GetNumber(InitIntegerModP(y,Zp))
	  if result_vahss {
			//if result_RP {
	    fmt.Println("Verify ok")
			//}else{
		//		fmt.Println("BPRP error")
		//	}
	  }else{
	    fmt.Println("Verification ERROR")
	  }
	  fmt.Println("Sum is", sum)

}

r, _ = rand.Int(rand.Reader, bn256.Order)
*/


import(
  "math/big"
  _"encoding/json"
  "fmt"
  "github.com/ing-bank/zkrp/crypto/bn256"
  "hannaekthesis/ccs08"
  ."hannaekthesis/vahss"
  "crypto/rand"

)
func Generate_set() (ccs08.ParamsSet){
    set := make([]int64,4)
    set[0] = 12
    set[1] = 42
    set[2] = 61
    set[3] = 71
    p, _ := ccs08.SetupSet(set)
    return p
}

func Verify_RP(proofs []ccs08.ProofSet, set ccs08.ParamsSet) (bool){
  var ok bool = true
  for _,proof := range proofs{
    ok, _ = ccs08.VerifySet(&proof,&set)
  }
  return ok
}
func Main_ccs(){
  const nr_clients = int64(4)
  const nr_servers = int64(4)

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
      secret := big.NewInt(12)//rand.Int(rand.Reader, big.NewInt(200) )//new(big.Int).SetInt64(int64(40))
      R_i,_ := rand.Int(rand.Reader, big.NewInt(400))
      //R_i := big.NewInt(0)
      c_i := InitClient(id,secret,t,R_i,Zp) /// include vahss?
      R_is = ModAdd(InitIntegerModP(R_is,Zp),InitIntegerModP(R_i,Zp)).Num
      clients = append(clients, c_i)
    }else{
      secret := big.NewInt(12)//rand.Int(rand.Reader, big.NewInt(200) )//new(big.Int).SetInt64(int64(40))
      R_n := InitIntegerModP(new(big.Int).Sub( new(big.Int).Mul(phiN, new(big.Int).Div( R_is, phiN ) ) ,  R_is),Zp).Num
      //R_i := int64( math.Ceil(float64(R_is) / float64((q-1) )) * float64((q-1)) )- R_is
      c_n := InitClient(id, secret, t, R_n, Zp) /// include vahss?
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

  result_vahss := Verify(tau_is, nr_clients, sigma, y,Zp)
  result_RP := Verify_RP(range_proofs,set)
  fmt.Println(result_RP)
  sum := InitIntegerModP(y,Zp).Num
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

package main

import(
    . "github.com/ing-bank/zkrp/bulletproofs"
    "fmt"
    "C"
)
//export GoSetupGeneric
func GoSetupGeneric(x, y int64) {
  //params, errSetup := SetupGeneric(18, 200)

	fmt.Printf("Go says: adding %v and %v\n", x, y)

  SetupGeneric(x,y)
}

//export GoProveGeneric
func GoProveGeneric(a int64, b int64, x int64, proof C.Proof){
  params, errSetup := SetupGeneric(a, b)
  if errSetup != nil {
    //  t.Errorf(errSetup.Error())
    //  t.FailNow()
  }

  // Create the proof
  bigSecret := new(big.Int).SetInt64(int64(x))
  Go_proof, errProve := ProveGeneric(bigSecret, params)
  if errProve != nil {
      //t.Errorf(errProve.Error())
      //t.FailNow()
  }
  transer Go_proof -> proof
  return proof
}

//export GoVerifyGeneric
func GoVerifyGeneric(proof C.Proof){

  // convert to Go
  // Verify the proof
  ok, errVerify := *proof.Verify()

  if errVerify != nil {
      //t.Errorf(errVerify.Error())
      //t.FailNow()
  }
  //assert.True(t, ok, "should verify")
}

func main() {} // Required but ignored

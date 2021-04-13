package vahss

import(
  "math/big"
)

type Server struct {
    Id     int
    Shares []*big.Int
    mod *Modular
}

 func InitServer(j int, mod *Modular) (*Server){
   s_j := new(Server)
   s_j.Id = j
   s_j.mod = mod
   s_j.Shares = nil
   return s_j
 }

 func Set_share(s_j *Server, i int, share *big.Int){
   s_j.Shares = append(s_j.Shares,share)
  }
func GetServerId(s_j *Server) (int){
  return s_j.Id
}

func GetShares(s_j *Server) ([]*big.Int){
  return s_j.Shares
}

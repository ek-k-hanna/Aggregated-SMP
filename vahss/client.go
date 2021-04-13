package vahss

import(
  "math/big"
  //"github.com/ing-bank/zkrp/crypto/p256"
  //  "hannaekthesis/bulletproof"
)
type Client struct {
    Id     int
    t     int64
    Xi   *big.Int
    Ri   *big.Int
    //tau_i *p256.P256
    mod  *Modular
}

func GetClientId(c_i *Client) int{
  return c_i.Id
}
/*
func Generate_shares(c_i *Client, nr_servers int64) ([]*big.Int){
  var shares_Zp []*IntegerModP
  //var tau_i *p256.P256
  shares_Zp = gen_secret_share_additive_with_hash_functions(c_i.i, c_i.Xi, c_i.t, nr_servers,c_i.mod)
  //tau_i = gen_tau(c_i.xi,c_i.Ri,c_i.g, params)
  // convert []IntegerModP -> []big.Int
  var shares []*big.Int
  for _,share := range shares_Zp {
    shares = append(shares,share.num)
  }
  //c_i.tau_i = tau_i

  return shares
}
*/
/*
func Get_tau(c_i *Client) (*p256.P256){
  return c_i.tau_i
}
*/

func InitClient(id_client int, secret_input *big.Int, t int64, r_i *big.Int, mod *Modular) (*Client){
    c_i := new(Client)
    c_i.Id = id_client
    c_i.Xi = secret_input
    c_i.t = t
    c_i.Ri = r_i
    c_i.mod = mod
    return c_i
  }

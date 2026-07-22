package lib

import (
	"context"
	"crypto/ecdsa"
	"math"
	"os"
	"strconv"
	"strings"

	"fmt"
	"log"
	"math/big"

	"cointrade/store"
	token "cointrade/token"
	"cointrade/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

const (
	USDT_ADDR = "0xdac17f958d2ee523a2206206994597c13d831ec7"
)

type myerror struct {
	ErrorMsg string
}

func (m *myerror) Error() string {
	return m.ErrorMsg
}

var CoinAddressList = map[string]string{ //币种列表
	"usdt":  "0xdac17f958d2ee523a2206206994597c13d831ec7",
	"sushi": "0x6b3595068778dd592e39a122f4f5a5cf09c90fe2",
	"usdc":  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
	"uni":   "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984",
	"aave":  "0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9",
	"yfi":   "0x0bc529c00C6401aEF6D220BE8C6Ea1667F6Ad93e",
	"dai":   "0x6b175474e89094c44da98b954eedeac495271d0f",
	"link":  "0x514910771af9ca656af840dff83e8264ecf986ca",
	"LON":   "0x0000000000095413afc295d19edeb1ad7b71c952",
	"CRV":   "0xD533a949740bb3306d119CC777fa900bA034cd52",
	"WBTC":  "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
	"WETH":  "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
	"CONV":  "0xc834fa996fa3bec7aad3693af486ae53d8aa8b50",
	"inj":   "0xe28b3B32B6c345A34Ff64674606124Dd5Aceca30",
	"MKR":   "0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2",
	"ALPHA": "0xa1faa113cbe53436df28ff0aee54275c13b40975",
	"BAND":  "0xba11d00c5f74255f56a5e366f4f77f5a186d7f55",
	"snx":   "0xc011a73ee8576fb46f5e1c5751ca3b9fe0af2a6f",
	"comp":  "0xc00e94cb662c3520282e6f5717214004a7f26888",
	"sxp":   "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9",
	"FTT":   "0x50d1c9771902476076ecfc8b2a83ad6b9355a4c9",
	"ust":   "0xa47c8bf37f92abed4a126bda807a7b7498661acd",
	"TRIBE": "0xc7283b66eb1eb5fb86327f08e1b5816b0720212b",
	"wise":  "0x66a0f676479Cee1d7373f3DC2e2952778BfF5bd6",
	"RRAX":  "0x853d955acef822db058eb8505911ed77f175b99e",
	"CORE":  "0x62359Ed7505Efc61FF1D56fEF82158CcaffA23D7",
	"mir":   "0x09a3ecafa817268f77be1283176b946c4ff2e608",
	"DPI":   "0x1494ca1f11d487c2bbe4543e90080aeba4ba3c2b",
	"luna":  "0xd2877702675e6ceb975b4a1dff9fb7baf4c91ea9",
	"HEZ":   "0xEEF9f339514298C6A857EfCfC1A762aF84438dEE",
	"fxs":   "0x3432b6a60d23ca0dfca7761b7ab56459d9c964d0",
	"fei":   "0x956f47f50a910163d8bf957cf5846d573e7f87ca",
	"mkd":   "0xcdA431ae623ceb0D1F97039F6cA95960C599d688",
	"usdt1": "0xF60c2B9aeE608F41425fC5b8027fA8e7DDff1524",
}

type EthLib struct {
	Client         *ethclient.Client //客户端实例
	Type           string            //币种
	BlockHash      string            //最后一次交易的块HEX
	TransFerHeader string            //授权交易头
	TransHeader    string            //转账交易头
	ApproveHeader  string            //授权头
}
type ErcTransInfo struct {
	FromAddress string
	//ToAddress       string
	ContractAddress string
	//Amount          string
	AbiStruct *AbiDecodeStruct
	TxId      string
}
type AbiDecodeStruct struct {
	MethodCode string
	Params     [][]byte
}

func (m *EthLib) GetTokenAddress() string {
	address, ok := CoinAddressList[m.Type]
	if !ok {
		return USDT_ADDR
	}
	return address
}
func (m *EthLib) Close() {
	m.Client.Close()
}
func (m *EthLib) CreateClient() bool {
	rpcURL := strings.TrimSpace(os.Getenv("ETH_RPC_URL"))
	if rpcURL == "" {
		return false
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Println(err)
		return false
	}
	m.ApproveHeader = string(crypto.Keccak256([]byte("approve(address,uint256)"))[0:4])
	m.TransFerHeader = string(crypto.Keccak256([]byte("transferFrom(address,address,uint256)"))[0:4])
	m.TransHeader = string(crypto.Keccak256([]byte("transfer(address,uint256)"))[0:4])
	m.Client = client
	return true
}
func (m *EthLib) GetBlockNumber() uint64 {
	n, err := m.Client.BlockNumber(context.Background())
	if err != nil {
		return 0
	}
	return n
}
func (m *EthLib) GetAbiStruct(data []byte) *AbiDecodeStruct {
	if len(data) < 4 {
		return nil
	}
	rs := new(AbiDecodeStruct)
	rs.Params = make([][]byte, 0)
	methodcode := data[0:4]
	rs.MethodCode = string(methodcode)
	if rs.MethodCode != m.ApproveHeader && rs.MethodCode != m.TransFerHeader && rs.MethodCode != m.TransHeader {
		return nil
	}
	param_data := data[4:]
	param_data_len := len(param_data)
	//fmt.Println("datalen:", len(data))
	//fmt.Println("len:", len(param_data))
	for i := 0; i < param_data_len; i = i + 32 {
		if len(param_data) < i+32 {
			break
		}
		tmp_data := param_data[i : i+32]

		rs.Params = append(rs.Params, tmp_data)
	}
	return rs
}
func (m *EthLib) ApiToAmount(data []byte) float64 {
	//从ABI代码得到额度
	u_n := fmt.Sprintf("0x%x", m.ClearLeftPad(data))
	n, _ := strconv.ParseInt(u_n, 0, 64)
	n_f := big.NewFloat(float64(n))
	value := new(big.Float).Quo(n_f, big.NewFloat(math.Pow10(int(6))))
	v, _ := value.Float64()
	return v

}
func (m *EthLib) ApiToAddress(data []byte) string {
	//解析ABI代码成地址
	return fmt.Sprintf("0x%x", m.ClearPad(data))
}
func (m *EthLib) ClearPad(data []byte) []byte {
	rs := make([]byte, 0)
	n := len(data)
	for i := 0; i < n; i++ {
		if data[i] == 0 {
			continue
		}
		rs = append(rs, data[i])
	}
	//开始反向清除
	n = len(rs)
	for n > 0 {
		if rs[n-1] != 0 {
			break
		}
		n = n - 1
	}
	return rs[0:n]
}
func (m *EthLib) ClearLeftPad(data []byte) []byte {
	rs := make([]byte, 0)
	n := len(data)
	for i := 0; i < n; i++ {
		if data[i] == 0 {
			continue
		}
		rs = append(rs, data[i])
	}
	return rs
}

func (m *EthLib) GetTranslist(blocknumber uint64) []*ErcTransInfo {

	block, err := m.Client.BlockByNumber(context.TODO(), nil)
	if err != nil {
		return nil
	}
	ts := block.Transactions()
	rs := make([]*ErcTransInfo, 0)

	for _, v := range ts {

		//fmt.Println("trans:", v)
		tmp := new(ErcTransInfo)
		//json, _ := v.MarshalJSON()

		//fmt.Println(string(json))
		if v == nil || v.To() == nil {
			continue
		}
		/*chainID, err := m.Client.NetworkID(context.TODO())
		if err != nil {
			//log.Fatal(err)
			continue
		}*/
		ms, e := v.AsMessage(types.LatestSignerForChainID(v.ChainId()), v.GasPrice())
		if e != nil {
			fmt.Println("error", e.Error())
			fmt.Println(tmp.FromAddress)
			continue
		}

		tmp.FromAddress = strings.ToLower(ms.From().Hex())
		tmp.ContractAddress = strings.ToLower(v.To().String())
		//fmt.Println("hash:", v.Hash().String())

		tmp.AbiStruct = m.GetAbiStruct(v.Data())

		tmp.TxId = v.Hash().String()
		//bs, _ := v.MarshalJSON()
		//fmt.Println(string(bs))
		//fmt.Println(tmp.ToAddress)
		//fmt.Println(tmp.ToAddress)
		rs = append(rs, tmp)
	}
	return rs
}
func (m *EthLib) CreateWalletAddress() (string, string) {
	//生成新钱包地址，返回公钥与私钥
	privateKey, err := crypto.GenerateKey()

	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	priKey := hexutil.Encode(privateKeyBytes)[2:]  //私钥

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	return address, priKey
}

func (m *EthLib) GetBalance(address string) *big.Float { //获取地址余额 ETH余额
	account := common.HexToAddress(address)
	balance, err := m.Client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	b := new(big.Float)
	b.SetString(balance.String())
	return new(big.Float).Quo(b, big.NewFloat(math.Pow10(18)))
}
func (m *EthLib) GetBalanceOfUsdt(usdt_address string) *big.Float { //获取代币余额
	//获取USDT余额
	tokenAddress := common.HexToAddress(m.GetTokenAddress()) //查询USDT
	instance, err := token.NewToken(tokenAddress, m.Client)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress(usdt_address)
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	fmt.Println(decimals)
	if err != nil {
		log.Fatal(err)
	}

	fbal := new(big.Float)
	fbal.SetString(bal.String())
	utils.Log(bal)
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("balance: %f", value) // "balance: 74605500.647409"
	return value
}
func (m *EthLib) TransUsdt(from string, fromkey string, to string, value int) (bool, error) { //value是需要 进位decimal 后的大整数
	//转出USDT

	privateKey, err := crypto.HexToECDSA(fromkey)
	if err != nil {
		utils.Log("1111111111")
		return false, err
	}
	from_address := common.HexToAddress(from) //发送地址 和 私钥对应
	to_address := common.HexToAddress(to)     //接收地址 不需要私钥
	nonce, err := m.Client.NonceAt(context.Background(), from_address, nil)
	if err != nil {
		//log.Fatal("2222222222", err)
		utils.Log("22222222222")
		return false, err
	}
	gasPrice, err := m.Client.SuggestGasPrice(context.Background()) //交易需要收取的GAS值
	tokenAddress := common.HexToAddress(m.GetTokenAddress())        //得到合约地址
	tokenObj, err := token.NewToken(tokenAddress, m.Client)         //生成合约实例
	decimals, err := tokenObj.Decimals(&bind.CallOpts{})            //得到进位
	opts := bind.NewKeyedTransactor(privateKey)
	//opts.GasLimit = uint64(300000)        //设置GAS花费的上限值
	opts.Nonce = big.NewInt(int64(nonce)) //设置NONCE值
	opts.Value = big.NewInt(0)            //wei的花费
	opts.GasPrice = gasPrice              //设置花费的燃气
	opts.From = from_address              //设置发送地址
	opts.Context = context.Background()
	fmt.Println("decimals:", decimals)
	f_value := m.BigIntPow10(int64(value), int64(decimals))
	rs, err := tokenObj.Transfer(opts, to_address, f_value) //调用合约的转账函数
	if err != nil {
		//log.Fatal("3333333333", err)
		utils.Log("333333333333")
		return false, err
	}
	m.BlockHash = rs.Hash().String()
	return true, nil //成功发送
}
func (m *EthLib) BigIntPow10(n int64, pow_n int64) *big.Int {
	tmp := big.NewInt(n)
	for i := 1; i <= int(pow_n); i++ {
		tmp = tmp.Mul(big.NewInt(10), tmp)
	}
	return tmp
}
func (m *EthLib) ApproveTransUsdt(from string, spender string, spenderKey string, to string, value float64) (bool, error) { //from 授权地址 spender 被授权地址 spenderKey 被授权地址私钥 to 接收地址 value 数额
	//通过授权转出 这里只是使用授权者的地址 私钥等信息使用被授权者的

	privateKey, err := crypto.HexToECDSA(spenderKey)
	if err != nil {
		utils.Log("111111111111")
		return false, err
	}
	from_address := common.HexToAddress(from) //发送地址 和 私钥对应
	spender_address := common.HexToAddress(spender)
	fmt.Println("spender:", spender_address)
	to_address := common.HexToAddress(to) //接收地址 不需要私钥
	nonce, err := m.Client.NonceAt(context.Background(), spender_address, nil)
	fmt.Println("nonce", nonce)
	if err != nil {
		utils.Log("222222222222")
		return false, err
	}
	gasPrice, err := m.Client.SuggestGasPrice(context.Background()) //交易需要收取的GAS值
	tokenAddress := common.HexToAddress(m.GetTokenAddress())        //得到合约地址
	tokenObj, err := token.NewToken(tokenAddress, m.Client)         //生成合约实例
	if err != nil {
		return false, err
	}
	decimals, err := tokenObj.Decimals(&bind.CallOpts{})
	opts := bind.NewKeyedTransactor(privateKey)

	opts.GasLimit = uint64(80000)                                      //设置GAS花费的上限值
	opts.Nonce = big.NewInt(int64(nonce))                              //设置NONCE值
	opts.Value = big.NewInt(0)                                         //wei的花费
	opts.GasPrice = big.NewInt(int64(float64(gasPrice.Int64()) * 1.4)) //设置花费的燃气

	//opts.From = spender_address           //设置发送地址
	fmt.Println("gasprice:", gasPrice)
	opts.Context = context.Background()
	utils.Log("decimal", decimals)
	//val := big.NewInt(int64(value) * int64(math.Pow10(int(decimals))))
	f_value := m.BigIntPow10(int64(value), int64(decimals))
	utils.Log("val:", value)
	//_, err = tokenObj.Transfer(opts, to_address, big.NewInt(int64(value)*int64(math.Pow10(int(decimals))))) //调用合约的转账函数
	t, err := tokenObj.TransferFrom(opts, from_address, to_address, f_value)

	if err != nil {
		utils.Log("333333333333")
		return false, err
	}
	utils.Log("cost", t.Cost())
	utils.Log("gas:", t.Gas())
	utils.Log("hash:", t.Hash())
	utils.Log("chainid:", t.ChainId())
	utils.Log("value:", t.Value())
	utils.Log(t.Data())
	m.BlockHash = t.Hash().String()
	return true, err
}
func (m *EthLib) GetAddressFromPrivateKey(private_key string) (string, error) {
	//从密钥得到钱包地址
	privateKey, err := crypto.HexToECDSA(private_key)
	if err != nil {
		return "", err
	}
	publickey := privateKey.Public()
	publicKeyECDSA, ok := publickey.(*ecdsa.PublicKey)
	if !ok {

		return "", &myerror{ErrorMsg: "error public key"}
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, nil
}

func (m *EthLib) CheckApprove(from_address string, to_address string) (bool, error) {
	//检测是否已经授权
	from_address_hex := common.HexToAddress(from_address)
	to_address_hex := common.HexToAddress(to_address)
	tokenAddress := common.HexToAddress(m.GetTokenAddress())
	tokenObj, err := token.NewToken(tokenAddress, m.Client) //生成合约实例
	if err != nil {
		return false, err
	}
	rs, err := tokenObj.Allowance(nil, from_address_hex, to_address_hex)
	fmt.Println(rs)
	if err != nil {
		return false, err
	}
	if rs.Int64() != 0 {
		return true, nil
	}
	return false, &myerror{ErrorMsg: "small"}
}

func (m *EthLib) PubErc20Contrct(address string, priKey string) error { //密钥
	//发布ERC20代币
	privateKey, _ := crypto.HexToECDSA(priKey)
	auth := bind.NewKeyedTransactor(privateKey)
	//nonce, err := m.Client.NonceAt(context.Background(), common.HexToAddress(address), nil)
	/*if err != nil {
		utils.Log("222222222222")
		return err
	}*/
	gasPrice, err := m.Client.SuggestGasPrice(context.Background()) //交易需要收取的GAS值
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("su gas:", gasPrice.Uint64())
	auth.GasPrice = gasPrice
	abi, _ := store.StoreMetaData.GetAbi()
	a, tx, _, err := bind.DeployContract(auth, *abi, common.FromHex(store.StoreBin), m.Client, "Tether USD", "USDT", big.NewInt(6000000000))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("address:", a.Hex())
	fmt.Println("tx:", tx.Hash().Hex())
	return nil
}
func (m *EthLib) GetTransState(txid string) int {
	rs, err := m.Client.TransactionReceipt(context.TODO(), common.HexToHash(txid))
	if err == nil {
		return int(rs.Status)
	}
	return 0
}
func (m *EthLib) GetTrans(txid string) (trans *ErcTransInfo, ispending bool, err error) {
	t, isp, err := m.Client.TransactionByHash(context.TODO(), common.HexToHash(txid))
	if err != nil {
		return nil, false, err
	}
	if ispending {
		return nil, isp, err
	}
	rs := new(ErcTransInfo)
	rs.AbiStruct = m.GetAbiStruct(t.Data())
	return rs, isp, err
}

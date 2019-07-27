// sendAPI project sendAPI.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/golang/protobuf/proto"
	"github.com/andyYuanFZM/wallet/walletapi/sendAPI/sender/bty/crypto"
	_ "github.com/andyYuanFZM/wallet/walletapi/sendAPI/sender/bty/crypto/secp256k1"
	"github.com/andyYuanFZM/wallet/walletapi/sendAPI/sender/bty/types"

	"encoding/hex"
	"io/ioutil"
	"bytes"
)

const (
	SECP256K1 = 1
	ED25519   = 2
	SM2       = 3
)

var (
	Jrpc_Url = "http://47.96.14.136:8801"
	privateKey = "3990969DF92A5914F7B71EEB9A4E58D6E255F32BF042FEA5318FC8B3D50EE6E8"
	toAddress = "1BBxz3iRkEniD9efX8zh28ty1PsvBLeurF"
)

func SignRawTransaction(data []byte, prive string) (result []byte, err error) {

	sigtx, err := signRawTx(data, prive)
	if err != nil {
		fmt.Println("SignRawTransaction ", "Err:", err)
		return nil, err
	}
	signjson := make(map[string]interface{}, 0)
	signjson["id"] = 1
	signjson["result"] = sigtx
	signjson["error"] = err
	resultdata, err := json.Marshal(signjson)
	if err != nil {
		return nil, err
	}
	return resultdata, nil
}

//签名
func signRawTx(tx []byte, privkey string) (string, error) {
	var txdata types.Transaction
	err := json.Unmarshal(tx, &txdata)
	if err != nil {

		err = types.Decode(tx, &txdata)
		if err != nil {
			return "", err
		}
	}

	hexprivkey, err := FromHex(privkey)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	cr, err := crypto.New(GetSignatureTypeName(SECP256K1))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	priv, err := cr.PrivKeyFromBytes(hexprivkey)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	txdata.Signature = nil

	data := Encode(&txdata)
	pub := priv.PubKey()
	sign := priv.Sign(data)
	txdata.Signature = &types.Signature{SECP256K1, pub.Bytes(), sign.Bytes()}

	return hex.EncodeToString(Encode(&txdata)), nil

}

func FromHex(s string) ([]byte, error) {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
		if len(s)%2 == 1 {
			s = "0" + s
		}
		return hex.DecodeString(s)
	}
	return nil, nil
}

func ToHex(b []byte) string {
	hexString := hex.EncodeToString(b)
	// Prefer output of "0x0" instead of "0x"
	if len(hexString) == 0 {
		hexString = "0"
	}
	return "0x" + hexString
}

func Encode(data proto.Message) []byte {
	b, err := proto.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}

func Decode(data []byte, msg proto.Message) error {
	return proto.Unmarshal(data, msg)
}

func GetSignatureTypeName(signType int) string {
	if signType == 1 {
		return "secp256k1"
	} else if signType == 2 {
		return "ed25519"
	} else if signType == 3 {
		return "sm2"
	} else {
		return "unknow"
	}
}

func main() {

	// 调用RPC接口构造交易
	txHex, _:= createRawTranferTx(0.1)

	txByteData, err := FromHex(txHex)
	if err != nil {
		fmt.Println("Error FromHex" + err.Error())
	}
	var tx types.Transaction
	err = types.Decode(txByteData, &tx)
	if err != nil {
		fmt.Println("Error Decode" + err.Error())
	}

	// 本地签名交易
	result,_ := SignRawTransaction(txByteData, privateKey)
	var resultdata = make(map[string]interface{})
	err = json.Unmarshal(result, &resultdata)
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	if hextx, ok := resultdata["result"]; ok {
		fmt.Println(hextx.(string))
	}
}


func createRawTranferTx(amount float64) (string, error) {
	postdata := fmt.Sprintf(`{"jsonrpc":"2.0","id":4,"method":"Chain33.CreateRawTransaction",
   "params":[{"to":"%v","amount":%v,"fee":%v,"note":"%v","isWithdraw":%v}]}`,
		toAddress, int64(amount*1e8), 0.001*1e8, "test_sign", false)
	resp, err := http.Post(getJrpc(), "application/json", bytes.NewBufferString(postdata))
	if err != nil {
		fmt.Println("err:", err.Error())
		return "", err
	}
	fmt.Printf("postdata:%v\n", postdata)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err:", err.Error())
		return "", err
	}
	var txdata = make(map[string]interface{})
	fmt.Println("resp:", string(data))
	err = json.Unmarshal(data, &txdata)
	if err != nil {
		fmt.Println("err:", err.Error())
		return "", err
	}
	if hextx, ok := txdata["result"]; ok {
		return hextx.(string), nil
	}
	return "", fmt.Errorf("not have result!")

}


func getJrpc() string {
	return Jrpc_Url
}
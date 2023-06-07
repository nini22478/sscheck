package util

import (
	"bytes"
	"check_vpn/mylog"
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func AnyGet(url, met, keyt string) (decryptedData []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to send HTTP request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	// 读取响应主体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}
	//mylog.Logf("ret:body:%v", string(body))
	// 对响应数据进行解密
	decodedData, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		fmt.Printf("Failed to decode response data: %v\n", err)
		return
	}

	key := []byte(keyt) // 定义对称加密的密钥，长度必须为16、24或32字节
	block, err := aes.NewCipher(key)
	if err != nil {
		mylog.Logf("Failed to create AES cipher: %v\n", err)
		return
	}

	decryptedData = make([]byte, len(decodedData))
	for i := 0; i < len(decodedData); i += block.BlockSize() {
		block.Decrypt(decryptedData[i:i+block.BlockSize()], decodedData[i:i+block.BlockSize()])
	}
	//jsonData := string(decryptedData)
	//cleanedJSON := strings.ReplaceAll(jsonData, "\x02", "")

	//decryptedData = []byte(cleanedJSON)
	//mylog.Logf("pre:%v\nnew:%v")
	return
}
func GetWgRet(url, meto, key string) (WgBody, error) {
	decryptedData, err := AnyGet(url, meto, key)
	// 发送HTTP GET请求
	mylog.Logf("unDecrypted data: %v", decryptedData)

	var users WgBody
	//newdecryptedData := bytes.Replace(decryptedData, []byte(`\`), []byte(""), -1)
	//rth, err := simplejson.NewJson(decryptedData)
	//items := rth.Get("test").MustArray()
	//for _, i := range items {
	//	nm := map(i)
	//	users.Test = append(users.Test, WgItem{
	//		UserDNS: i["user_DNS"],
	//	})
	//	mylog.Logf("Decrypted data: %v", i)
	//}
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b := bytes.NewBuffer(decryptedData)
	d := json.NewDecoder(b)
	err = d.Decode(&users)

	//err = json.Unmarshal(decryptedData, &users)
	if err != nil {
		mylog.Logf("Error:%v", err)
		return users, err
	}

	mylog.Logf("Decrypted data: %v", users.Test)
	return users, err
}
func GetSsRet(url, meto, key string) (SsBody, error) {
	decryptedData, err := AnyGet(url, meto, key)
	// 发送HTTP GET请求

	var users SsBody
	//newdecryptedData := bytes.Replace(decryptedData, []byte(`\`), []byte(""), -1)
	mylog.Logf("Decrypted data: %v", string(decryptedData))
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(decryptedData, &users)
	if err != nil {
		mylog.Logf("Error:%v", err)
		return users, err
	}

	mylog.Logf("Decrypted data: %v", users.Test)
	return users, err
}

type SsBody struct {
	Test []SsItem
}
type SsItem struct {
	Ip     string
	Port   string
	Method string
	Paskey string
}
type WgBody struct {
	Test []WgItem
}
type WgItem struct {
	AllowedIPs      string `json:"AllowedIPs"`
	ServerPublicKey string `json:"server_PublicKey"`
	UserAddress     string `json:"user_Address"`
	UserDNS         string `json:"user_DNS"`
	UserEndpoint    string `json:"user_Endpoint"`
	UserPrivateKey  string `json:"user_PrivateKey"`
}

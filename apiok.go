package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	// "reflect"
	"time"
	// "strconv"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type configuration struct {
    APIKey string
	SecretKey string
	Passphrase string
	Timestamp string
}

var conf = configuration{}
var baseUrl = "https://www.okex.win/"

func main() {
	fmt.Println("hello, world!")

	file, _ := os.Open("apiok_conf.json")
    defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&conf)
    if err != nil {
        fmt.Println("Error:", err)
    }

	get("api/v5/account/account-position-risk")

	// http.HandleFunc("/test", test)
	http.ListenAndServe("0.0.0.0:8089", nil)
}

func get(url string){
	client := &http.Client{}

    //提交请求
    reqest, err := http.NewRequest("GET", baseUrl+url, nil)

	if err != nil {
        fmt.Println(err)
    }

	//reflect.TypeOf(time.Now().Unix())
	//2020-12-08T09:08:57.715Z
	//2021-05-13T14:28:50.171Z
	conf.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	fmt.Println(conf.Timestamp + "GET" + url)

	sign := hmacSha256(conf.Timestamp + "GET" + url, conf.SecretKey)

	//增加header选项
    reqest.Header.Add("OK-ACCESS-KEY", conf.APIKey)
    reqest.Header.Add("OK-ACCESS-SIGN", sign)
    reqest.Header.Add("OK-ACCESS-TIMESTAMP", conf.Timestamp)
	reqest.Header.Add("OK-ACCESS-PASSPHRASE", conf.Passphrase)

	//处理返回结果
    response, _ := client.Do(reqest)
    defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
    //fmt.Println(string(body))
    fmt.Printf("Get request result: %s\n", string(body))
}

func hmacSha256(data string, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}


// ===================  demo  =============================

type Auth struct {
    Username string `json:"username"`
    Pwd string `json:"password"`
}

func test(writer http.ResponseWriter,  request *http.Request){
	// var result  Resp
	// result.Code = "200"
	// result.Msg = "登录成功"
	// request.ParseForm()
    // username, _ :=  request.Form["username"][0]
    // pwd, _ :=  request.Form["password"][0]
	// if err := json.NewEncoder(writer).Encode(result); err != nil {
    //     log.Fatal(err)
    // }

	var auth Auth
	if err:=json.NewDecoder(request.Body).Decode(&auth); err != nil {
		fmt.Println(err)
		// log.Fatal(err)
	}

	fmt.Println(auth.Username)
}



package main

import (
	"bytes"
	"compress/gzip"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var numRequests = 1000 //总请求数
// var uri = "http://localhost:11001"
var uri = "http://d-game.test.rancigame.com" //test

// 完整的加密函数实现
func encryptData(data interface{}, key string) (string, error) {
	// 1. JSON 序列化
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON marshal error: %v", err)
	}

	// 2. Gzip 压缩
	var gzipBuffer bytes.Buffer
	gzWriter := gzip.NewWriter(&gzipBuffer)
	if _, err := gzWriter.Write(jsonData); err != nil {
		return "", fmt.Errorf("gzip error: %v", err)
	}
	gzWriter.Close()

	// 3. 第一次 Base64 编码
	base64Once := base64.StdEncoding.EncodeToString(gzipBuffer.Bytes())

	// 5. DES-ECB 加密
	encrypted, err := desECBEncrypt([]byte(base64Once), key)
	if err != nil {
		return "", fmt.Errorf("DES encryption error: %v", err)
	}

	// 6. 最终 Base64 编码
	result := base64.StdEncoding.EncodeToString(encrypted)
	return result, nil
}

// DES-ECB 加密实现
func desECBEncrypt(data []byte, key string) ([]byte, error) {
	// 处理密钥（必须为 8 字节）
	keyBytes := []byte(key)
	if len(keyBytes) != 8 {
		// 调整密钥长度
		if len(keyBytes) > 8 {
			keyBytes = keyBytes[:8]
		} else {
			// 填充到 8 字节
			padding := make([]byte, 8-len(keyBytes))
			keyBytes = append(keyBytes, padding...)
		}
	}

	// 创建 DES 密码块
	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	// 数据填充（PKCS5）
	blockSize := block.BlockSize()
	data = pkcs5Padding(data, blockSize)

	// ECB 模式加密
	encrypted := make([]byte, len(data))
	for i := 0; i < len(data); i += blockSize {
		block.Encrypt(encrypted[i:i+blockSize], data[i:i+blockSize])
	}

	return encrypted, nil
}

// 解密函数实现 - 与 encryptData 对应
func decryptData(encryptedData string, key string) (map[string]interface{}, error) {
	// 1. Base64 解码（最终结果）
	decoded, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("Base64 decode error: %v", err)
	}

	// 2. DES-ECB 解密
	decrypted, err := desECBDecrypt(decoded, key)
	if err != nil {
		return nil, fmt.Errorf("DES decryption error: %v", err)
	}

	// 3. 第一次 Base64 解码（加密前的 Base64）
	base64Once, err := base64.StdEncoding.DecodeString(string(decrypted))
	if err != nil {
		return nil, fmt.Errorf("First Base64 decode error: %v", err)
	}

	// 4. Gzip 解压缩
	gzReader, err := gzip.NewReader(bytes.NewReader(base64Once))
	if err != nil {
		return nil, fmt.Errorf("Gzip reader error: %v", err)
	}
	defer gzReader.Close()

	jsonData, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, fmt.Errorf("Gzip decompress error: %v", err)
	}

	// 5. JSON 反序列化
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %v", err)
	}

	return result, nil
}

// DES-ECB 解密实现
func desECBDecrypt(encryptedData []byte, key string) ([]byte, error) {
	// 处理密钥（必须为 8 字节）
	keyBytes := []byte(key)
	if len(keyBytes) != 8 {
		if len(keyBytes) > 8 {
			keyBytes = keyBytes[:8]
		} else {
			padding := make([]byte, 8-len(keyBytes))
			keyBytes = append(keyBytes, padding...)
		}
	}

	// 创建 DES 密码块
	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	// ECB 模式解密
	decrypted := make([]byte, len(encryptedData))
	blockSize := block.BlockSize()

	for i := 0; i < len(encryptedData); i += blockSize {
		block.Decrypt(decrypted[i:i+blockSize], encryptedData[i:i+blockSize])
	}

	// 去除 PKCS5 填充
	decrypted = pkcs5Unpadding(decrypted)

	return decrypted, nil
}

// PKCS5 去除填充
func pkcs5Unpadding(data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	padding := int(data[len(data)-1])
	if padding > len(data) {
		return data
	}

	return data[:len(data)-padding]
}

// PKCS5 填充
func pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// 正确的 MD5 计算
func getHeaderSign(data map[string]string, signKey string) string {
	// 按照 PHP 代码的顺序拼接字符串
	tobeStr := data["channel"] +
		data["did"] +
		data["install-id"] +
		data["noise"] +
		data["package-name"] +
		data["run-id"] +
		data["timestamp"] +
		data["version"] +
		signKey

	// 计算 MD5
	hash := md5.Sum([]byte(tobeStr))
	return hex.EncodeToString(hash[:])
}

// 压力测试函数
func stressTest() {
	fmt.Println("=== 开始压力测试 ===")

	// 配置参数
	// url := "https://game-nuanchu.qmhd87.com/api/v1/user/card/drawCard"

	// url := "http://d-game.test.rancigame.com/api/v1/user/card/drawCard"
	// token := "Bearer eyJpdiI6IlpxRnhIOW1lVmg2RXpUcU9KYUFqUWc9PSIsInZhbHVlIjoid2JVUTdMVnlKbkR0dmZPNEsrc0o5b045WWl2bWM2Yjd4ajhxdjNRRGN2Wm9rYldTSUZBZmtIQmU3Q28ydTdPcEIxT2ZTVkFuK254YjdFMWtzSUxFQ3c9PSIsIm1hYyI6ImY2YWZlODcxYWNjZGMwNjI3MWIxYzczNzkxMjBlOTQ1OTk1ODdhZGZlOTAxOGJkNDAzM2RlNTJjNTgxOGM2OTEiLCJ0YWciOiIifQ=="

	url := "/api/v1/user/card/drawCard"
	token := "Bearer eyJpdiI6IndDZ3RYRzJoaVpETHZFNVVsQjVReGc9PSIsInZhbHVlIjoiWVFGUzcxSUVyR3haRE5NWS9aR3ZTck5YSlpJa2RmVnNpelgxRXc4WUVYT2dVUDlTdUZ3YzdrZS9Ia2pTQkY5d295VENrd21Fci8vMDZHcGNqTmtSUGc9PSIsIm1hYyI6IjcwY2ZmZDkyYzE3MjU5YzJiNGFhNTc4OTlhZTVlMWYyNGRkZTliN2RlMGE0OGUwMzYwZjZmM2U1YmQxZDA0NjIiLCJ0YWciOiIifQ=="

	numRequests := 5
	lastID := 1100
	concurrent := 10 // 并发数

	// 生成测试数据
	requests := make([]map[string]interface{}, numRequests)
	for i := 0; i < numRequests; i++ {
		requests[i] = map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id":            lastID + i,
					"packagestar":   3,
					"isgoldPackage": 0,
					"isXLPackage":   0,
					"cardfrom":      "None",
					"packageId":     1,
					"addRate":       1,
				},
			},
		}
	}

	// 创建通道用于并发控制
	jobs := make(chan map[string]interface{}, numRequests)
	results := make(chan *testResult, numRequests)

	// 启动时间
	startTime := time.Now()

	// 启动工作协程
	for w := 0; w < concurrent; w++ {
		go worker(w, url, token, jobs, results)
	}

	// 发送任务到通道
	for _, req := range requests {
		jobs <- req
	}
	close(jobs)

	// 收集结果
	successCount := 0
	failCount := 0

	for i := 0; i < numRequests; i++ {
		result := <-results
		if result.Success {
			successCount++
			fmt.Printf("✓ 请求 %d 成功: %s\n", result.ID, result.Message)
		} else {
			failCount++
			fmt.Printf("✗ 请求 %d 失败: %s\n", result.ID, result.Error)
		}

		// 显示进度
		progress := float64(i+1) / float64(numRequests) * 100
		fmt.Printf("进度: %.1f%% [%d/%d]\r", progress, i+1, numRequests)
	}

	// 计算统计信息
	elapsed := time.Since(startTime)

	fmt.Println("\n=== 压力测试结果 ===")
	fmt.Printf("总请求数: %d\n", numRequests)
	fmt.Printf("成功数量: %d\n", successCount)
	fmt.Printf("失败数量: %d\n", failCount)
	fmt.Printf("成功率: %.2f%%\n", float64(successCount)/float64(numRequests)*100)
	fmt.Printf("总耗时: %.2f 秒\n", elapsed.Seconds())
	fmt.Printf("平均响应时间: %.2f 毫秒\n", float64(elapsed.Milliseconds())/float64(numRequests))
	fmt.Printf("QPS: %.2f\n", float64(numRequests)/elapsed.Seconds())
}

// 测试结果结构体
type testResult struct {
	ID      int
	Success bool
	Message string
	Error   string
}

// 创建 headers
func createHeaders(token string) map[string]string {
	// 准备请求头
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	noise := strconv.Itoa(rand.Intn(900000) + 100000)

	headers := map[string]string{
		"timestamp":     timestamp,
		"authorization": token,
		"noise":         noise,
		"did":           "postman",
		"version":       "1.0.0",
		"channel":       "PD0301",
		"package-name":  "postman",
		"install-id":    "szy-postman",
		"run-id":        "szy-postman",
	}

	// 计算签名
	headers["sign"] = getHeaderSign(headers, "xstxa+s")
	return headers
}

// 工作协程函数
func worker(id int, url, token string, jobs <-chan map[string]interface{}, results chan<- *testResult) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for job := range jobs {
		result := &testResult{ID: id}

		// 准备请求头
		// timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		// noise := strconv.Itoa(rand.Intn(900000) + 100000)
		// headers := map[string]string{
		// 	"timestamp":     timestamp,
		// 	"authorization": token,
		// 	"noise":         noise,
		// 	"did":           "postman",
		// 	"version":       "1.0.0",
		// 	"channel":       "PD0301",
		// 	"package-name":  "postman",
		// 	"install-id":    "szy-postman",
		// 	"run-id":        "szy-postman",
		// }
		// // 计算签名
		// headers["sign"] = getHeaderSign(headers, "xstxa+s")

		headers := createHeaders(token)

		// 加密数据
		encryptedData, err := encryptData(job, "S2xasd_2")
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("加密失败: %v", err)
			results <- result
			continue
		}

		// 准备请求体
		requestBody := map[string]string{
			"eData": encryptedData,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("JSON 序列化失败: %v", err)
			results <- result
			continue
		}

		// 创建请求
		req, err := http.NewRequest("POST", uri+url, bytes.NewBuffer(jsonBody))
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("创建请求失败: %v", err)
			results <- result
			continue
		}

		// 设置请求头
		for key, value := range headers {
			req.Header.Set(key, value)
		}
		req.Header.Set("Content-Type", "application/json")

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("请求失败: %v", err)
			results <- result
			continue
		}

		// 读取响应
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("读取响应失败: %v", err)
			results <- result
			continue
		}

		// fmt.Println("响应:", string(body))

		// 解析响应
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("JSON 解析失败: %v", err)
			results <- result
			continue
		}

		// 检查响应
		if code, ok := response["code"].(float64); ok && code == 0 {
			result.Success = true
			if msg, ok := response["message"].(string); ok {
				result.Message = msg
			}
		} else {
			result.Success = false
			if msg, ok := response["message"].(string); ok {
				result.Error = msg
			} else {
				result.Error = "未知错误"
			}
		}

		results <- result

		// 添加延迟避免服务器过载
		time.Sleep(50 * time.Millisecond)
	}
}

// 测试解密功能
func testDecrypt() {
	// 假设这是加密后的数据
	encryptedStr := "ce86nwzgOyOrpn/wVH70HeWag1eBQqJuOvjxox82QRusddF0faaTH9oYn+OPoOm5JzuR2sZMxPmnHgrLDD29ors0LL/wbk32fgPOkQEuMmF5Cayo+rKg/+EmznpPHs/zNDsZM+dcj5oxHS8tqOjSiN41mTGqYuZ0ikInULX94gOiv/oN4R6kuQjeToapepLKXwVbUztn9yh0X4pkh95MUpsrnTkWzSigTUk/mPyJrk/f5buipMX9fveENIBuTopxaDehXvsqt11a1DqNumvm6gM/ZYVBPZQl7kK/+sGN/SAP729ME/ppAYMompUZwaOz7OsO7VUWktxlZdOq3PDFH0hOwNFh/stqc8qVDMb48FZtb5cGfylsuD3/9oDG0cI0AnbvuXq2jQmJ+iTjqzA6O/DVRjpLkjQx8gQPpHv0M3JrH1fKK6ulbcXTEnR4FEowRiB0Yn40WQZdGaQhIdvP4WEIdHe2McYmUScbtBiOWMDayXCQsfO0t3XoOr1JwF7cjsWWNQqp7L2UI5y7+KEvIzH+0gNwDb2pQegBpPdMsdQgAvQuy5jxNtDP2CJyxiye+GS6CJ9I1UKc3utk6iDvPkRya557FijTsRtJ2PokdK0EBs0RwIAVc3sTkh7rM/QEHOsVVVH/jLpBEiUMQqo4aSHv+idJ5ZPymrQGTUGaeDoHfYnEwoqGnQvpojfPjk9p+8MWj9vdxFPxFTQe6WAwYVUC8K57HuzqAgnM4BpuKhNIEVmksY6zjQ=="

	// 解密数据
	decryptedData, err := decryptData(encryptedStr, "S2xasd_2")
	if err != nil {
		fmt.Printf("解密失败: %v\n", err)
		return
	}

	fmt.Printf("解密成功: %+v\n", decryptedData)
}

// 创建用户
func login(account string) string {
	headers := createHeaders("")
	url := "/api/v1/login"

	// 创建用户
	user := map[string]interface{}{
		"account": account,
		"type":    1,
	}

	// 加密数据
	encryptedData, err := encryptData(user, "S2xasd_2")
	if err != nil {
		fmt.Println("加密失败: %v", err)
		return ""
	}

	// 准备请求体
	requestBody := map[string]string{
		"eData": encryptedData,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("JSON 序列化失败: %v", err)
		return ""
	}

	// 创建请求
	req, err := http.NewRequest("POST", uri+url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("创建请求失败: %v", err)
		return ""
	}
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败: %v", err)
		return ""
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("读取响应失败: %v", err)
		return ""
	}

	// fmt.Println("响应:", string(body))

	// 解析响应
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("JSON 解析失败: ", err)

		if err.Error() == "invalid character '<' looking for beginning of value" {
			fmt.Println("响应:", string(body))

			//中断程序
			os.Exit(1)
		}

		return ""
	}

	// 检查响应
	if code, ok := response["code"].(float64); ok && code == 0 {
		eData := response["eData"].(string)
		oData, err := decryptData(eData, "S2xasd_2")
		if err != nil {
			fmt.Println("解密失败: ", err)
			return ""
		}
		// fmt.Println("登录成功: %+v", oData)
		// fmt.Println("token: ", oData["token"])
		if token, ok := oData["token"].(string); ok {
			return token
		}
		return ""
	} else {
		fmt.Println("登录失败: ", response)
		panic("登录失败")
		return ""
	}
}

func main() {
	// 设置随机种子
	// rand.Seed(time.Now().UnixNano())

	// testDecrypt()

	// stressTest()

	ts()
}

// 登录，拿到多个token，然后批量请求
func ts() {
	tokenArr := []string{}
	// accounts := []string{}
	// accountNum := 100
	// for i := 0; i < accountNum; i++ {
	// 	account := fmt.Sprintf("a-test%d", i)
	// 	accounts = append(accounts, account)
	// }
	// for _, account := range accounts {
	// 	token := login(account)
	// 	if token != "" {
	// 		tokenArr = append(tokenArr, token)
	// 	}
	// }

	accountNum := 100
	for i := 0; i < accountNum; i++ {
		account := fmt.Sprintf("test%d", i)
		fmt.Println("account-->: ", account)
		token := login(account)
		if token != "" {
			tokenArr = append(tokenArr, token)
		}
	}
	fmt.Println("tokenArr: ", tokenArr)

	// 抽卡测试
	tsCardDraw(tokenArr)

	// 存档测试
	tsUserSave(tokenArr)
}

func tsCardDraw(tokenArr []string) {
	fmt.Println("\n=== 抽卡——压测开始 ===")

	url := "/api/v1/user/card/drawCard"

	// numRequests := 1000         //总请求数
	lastID := 2000              // 查数据库验重最大值
	concurrent := len(tokenArr) // 并发数，正常一个用户一个线程
	// concurrent := 10 // 并发数

	// 生成测试数据
	requests := make([]map[string]interface{}, numRequests)
	for i := 0; i < numRequests; i++ {
		requests[i] = map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id":            lastID + i,
					"packagestar":   3,
					"isgoldPackage": 0,
					"isXLPackage":   0,
					"cardfrom":      "None",
					"packageId":     1,
					"addRate":       1,
				},
			},
		}
	}

	// 创建通道用于并发控制
	jobs := make(chan map[string]interface{}, numRequests)
	results := make(chan *testResult, numRequests)

	// 启动时间
	startTime := time.Now()

	// 启动工作协程
	for w := 0; w < concurrent; w++ {
		// go worker(w, url, tokenArr[w%len(tokenArr)], jobs, results)
		go worker(w, url, tokenArr[w], jobs, results)
	}

	// 发送任务到通道
	for _, req := range requests {
		jobs <- req
	}
	close(jobs)

	// 收集结果
	successCount := 0
	failCount := 0

	for i := 0; i < numRequests; i++ {
		result := <-results
		if result.Success {
			successCount++
			fmt.Printf("✓ 请求 %d 成功: %s\n", result.ID, result.Message)
		} else {
			failCount++
			fmt.Printf("✗ 请求 %d 失败: %s\n", result.ID, result.Error)
		}

		// 显示进度
		progress := float64(i+1) / float64(numRequests) * 100
		fmt.Printf("进度: %.1f%% [%d/%d]\r", progress, i+1, numRequests)
	}

	// 计算统计信息
	elapsed := time.Since(startTime)

	fmt.Println("\n=== 抽卡——压测结果 ===")
	fmt.Printf("总请求数: %d\n", numRequests)
	fmt.Printf("成功数量: %d\n", successCount)
	fmt.Printf("失败数量: %d\n", failCount)
	fmt.Printf("成功率: %.2f%%\n", float64(successCount)/float64(numRequests)*100)
	fmt.Printf("总耗时: %.2f 秒\n", elapsed.Seconds())
	fmt.Printf("平均响应时间: %.2f 毫秒\n", float64(elapsed.Milliseconds())/float64(numRequests))
	fmt.Printf("QPS: %.2f\n", float64(numRequests)/elapsed.Seconds())
}

func tsUserSave(tokenArr []string) {
	fmt.Println("\n=== 存档——压测开始 ===")

	url := "/api/v1/user/"

	// numRequests := 1000         //总请求数
	concurrent := len(tokenArr) // 并发数，正常一个用户一个线程

	// 生成测试数据
	requests := make([]map[string]interface{}, numRequests)
	for i := 0; i < numRequests; i++ {
		jsonStr := `{
  "_id": {
    "$oid": "67f3babada8375c9eb0cfedd"
  },
  "uid": 212347000,
  "energy": 4645,
  "gems": 642,
  "star": 467,
  "area": "CN",
  "score": 947,
  "rankAt": 1750299070,
  "lps": {
    "ver": 12,
    "exp": 30,
    "lv": 8,
    "story": 1,
    "map": 1,
    "mapLvs": [
      0,
      0,
      0,
      0,
      0,
      0,
      0,
      0
    ]
  },
  "bindInfo": {
    "1": {
      "type": 1,
      "at": 1744026298
    }
  },
  "regInfo": {
    "at": 1744026298,
    "ip": "172.19.0.1",
    "did": "postman",
    "installId": "szy-postman",
    "runId": "szy-postman",
    "channel": "PD0301",
    "version": "1.0.0",
    "localAt": 0,
    "tzOfs": 8
  },
  "loginInfo": {
    "at": 1764731973,
    "ip": "113.110.203.13",
    "did": "123",
    "installId": "123",
    "runId": "123",
    "channel": "PD0301",
    "version": "1"
  },
  "authMd5": "3503a5dc2d33f041c6d60780a038c62d",
  "statInfo": {
    "loginNum": 179,
    "loginDay": 44,
    "amn": 4,
    "amnAt": 1764677536,
    "noticeAt": 1762573925,
    "money": 145200,
    "snapAt": 1760152804,
    "snapVer": 165,
    "payAt": 1761287690,
    "smMoney": 16800,
    "amnCsAt": 1764677536
  },
  "wxInfo": {
    "openid": "oiF3S5GdGO0EO-ltKOliANlaHftI",
    "at": 1755608049,
    "sessionkey": "8A2XZlQx9YMKRKgnRAYFMg=="
  },
  "updatedAt": 1764731973,
  "createdAt": 1744026298,
  "avatar": "xxx",
  "frame": "xxxx",
  "dayRank": {
    "topTimes": 2,
    "continueDays": 3,
    "addDate": "20250620"
  },
  "abs": {
    "ver": 12,
    "vars": [
      {
        "id": 1001,
        "at": 1761045750,
        "exp": 1761062399
      },
      {
        "id": 1205,
        "at": 1761045750,
        "exp": 1761062399
      },
      {
        "id": 1203,
        "at": 1761045750,
        "exp": 1761062399
      },
      {
        "id": 1201,
        "at": 1761045750,
        "exp": 1761062399
      },
      {
        "id": 1202,
        "at": 1761045750,
        "exp": 1761062399
      },
      {
        "id": 1204,
        "at": 1761045750,
        "exp": 1761062399
      }
    ]
  },
  "btls": {
    "ver": 658,
    "needVars": [
      {
        "id": 1026001,
        "roleId": 2,
        "star": 47,
        "eleVars": [
          2703,
          2702
        ],
        "cardId": 0
      },
      {
        "id": 1027001,
        "roleId": 3,
        "star": 51,
        "eleVars": [
          2702,
          2702,
          2702
        ],
        "cardId": 0
      },
      {
        "id": 1034001,
        "roleId": 5,
        "star": 165,
        "eleVars": [
          3506,
          3404
        ],
        "cardId": 0
      },
      {
        "id": 97531,
        "type": 2,
        "roleId": 7,
        "star": 1,
        "eleVars": [
          2908
        ],
        "cardId": 0
      },
      {
        "id": 1020010,
        "roleId": 2,
        "star": 42,
        "eleVars": [
          1704,
          1702
        ],
        "cardId": 0
      },
      {
        "id": 1020011,
        "roleId": 3,
        "star": 104,
        "eleVars": [
          1205,
          1705
        ],
        "cardId": 0
      },
      {
        "id": 1020012,
        "roleId": 4,
        "star": 130,
        "eleVars": [
          1206,
          1103
        ],
        "cardId": 0
      }
    ],
    "sceneVars": [
      {
        "id": 2706,
        "pos": 21,
        "cover": 2
      },
      {
        "id": 1402,
        "pos": 22,
        "cover": 2
      },
      {
        "id": 2601,
        "pos": 23,
        "cover": 2
      },
      {
        "id": 2104,
        "pos": 24,
        "cover": 2
      },
      {
        "id": 2704,
        "pos": 25,
        "cover": 2
      },
      {
        "id": 2103,
        "pos": 26,
        "cover": 2
      },
      {
        "id": 1403,
        "pos": 27,
        "cover": 2
      },
      {
        "id": 2204,
        "pos": 28,
        "cover": 2
      },
      {
        "id": 1802,
        "pos": 29,
        "cover": 1
      },
      {
        "id": 1502,
        "pos": 30,
        "cover": 1
      },
      {
        "id": 1801,
        "pos": 31,
        "cover": 1
      },
      {
        "id": 1301,
        "pos": 32,
        "cover": 1
      },
      {
        "id": 2303,
        "pos": 33,
        "cover": 2
      },
      {
        "id": 2102,
        "pos": 34,
        "cover": 2
      },
      {
        "id": 1503,
        "pos": 35,
        "cover": 1
      },
      {
        "id": 1706,
        "pos": 40,
        "cover": 1
      },
      {
        "id": 2201,
        "pos": 41,
        "cover": 2
      },
      {
        "id": 1501,
        "pos": 42,
        "cover": 1
      },
      {
        "id": 2001,
        "pos": 48,
        "cover": 1
      },
      {
        "id": 1705,
        "cover": 1
      },
      {
        "id": 2004,
        "pos": 6,
        "cover": 1
      },
      {
        "id": 2301,
        "pos": 7,
        "cover": 1
      },
      {
        "id": 1701,
        "pos": 11,
        "cover": 1
      },
      {
        "id": 2701,
        "pos": 13,
        "cover": 1
      },
      {
        "id": 2401,
        "pos": 14,
        "cover": 1
      },
      {
        "id": 1703,
        "pos": 17,
        "cover": 1
      },
      {
        "id": 1305,
        "pos": 18,
        "cover": 2
      },
      {
        "id": 1702,
        "pos": 19,
        "cover": 1
      },
      {
        "id": 2501,
        "pos": 20,
        "cover": 2
      },
      {
        "id": 1302,
        "pos": 49,
        "cover": 2
      },
      {
        "id": 2204,
        "pos": 50,
        "cover": 1
      },
      {
        "id": 1901,
        "pos": 51,
        "cover": 1
      },
      {
        "id": 1603,
        "pos": 52,
        "cover": 2
      },
      {
        "id": 2402,
        "pos": 53,
        "cover": 2
      },
      {
        "id": 2002,
        "pos": 54,
        "cover": 2
      },
      {
        "id": 2106,
        "pos": 55,
        "cover": 2
      },
      {
        "id": 2103,
        "pos": 56,
        "cover": 2
      },
      {
        "id": 2203,
        "pos": 57,
        "cover": 2
      },
      {
        "id": 1404,
        "pos": 58,
        "cover": 2
      },
      {
        "id": 1903,
        "pos": 59,
        "cover": 2
      },
      {
        "id": 2202,
        "pos": 60,
        "cover": 2
      },
      {
        "id": 2602,
        "pos": 61,
        "cover": 2
      },
      {
        "id": 1403,
        "pos": 62,
        "cover": 2
      },
      {
        "id": 2702,
        "pos": 8
      },
      {
        "id": 4004,
        "num": 1,
        "pos": 1
      },
      {
        "id": 1005,
        "at": 1755608009,
        "num": 89,
        "pos": 46,
        "drt": 1,
        "currCd": 75
      },
      {
        "id": 1101,
        "pos": 36
      },
      {
        "id": 1206,
        "pos": 10
      },
      {
        "id": 1201,
        "pos": 16
      },
      {
        "id": 1204,
        "pos": 12
      },
      {
        "id": 1103,
        "pos": 2
      },
      {
        "id": 9906,
        "num": 5,
        "pos": 5,
        "currCd": 600
      },
      {
        "id": 1202,
        "pos": 37
      },
      {
        "id": 1204,
        "pos": 38
      },
      {
        "id": 1203,
        "pos": 43
      }
    ],
    "popVars": [],
    "propVars": [
      {
        "id": 301,
        "at": 1749743025,
        "num": 1
      },
      {
        "id": 2702,
        "at": 0,
        "num": 99
      }
    ],
    "needCDVars": [
      {
        "type": 1,
        "canUseNum": 0,
        "cAt": 0,
        "dNum": 123
      },
      {
        "type": 2,
        "canUseNum": 0,
        "cAt": 0,
        "dNum": 123
      }
    ]
  },
  "buffs": {
    "ver": 82,
    "vars": {
      "1": {
        "id": 1,
        "vAt": 1761062399
      },
      "2": {
        "id": 2,
        "vAt": 1761062399
      },
      "3": {
        "id": 3,
        "vAt": 1761062399
      },
      "4": {
        "id": 4,
        "vAt": 2145888000,
        "num": 1750301563
      },
      "5": {
        "id": 5,
        "vAt": 2145888000,
        "num": 1761046363
      },
      "17": {
        "id": 17,
        "vAt": 1761062399,
        "num": 2
      },
      "203": {
        "id": 203,
        "vAt": 1761045809
      },
      "2071": {
        "id": 2071,
        "vAt": 1761132152,
        "num": 1
      },
      "2072": {
        "id": 2072,
        "vAt": 1761132152,
        "num": 1
      },
      "2074": {
        "id": 2074,
        "vAt": 1761046892,
        "num": 1
      },
      "2075": {
        "id": 2075,
        "vAt": 1761049412,
        "num": 1
      }
    }
  },
  "bwhs": {
    "ver": 1,
    "boxNum": 16,
    "vars": []
  },
  "eles": {
    "ver": 1,
    "vars": [
      {
        "id": 1001
      },
      {
        "id": 1002
      },
      {
        "id": 1003
      },
      {
        "id": 1004
      },
      {
        "id": 1201
      },
      {
        "id": 1202
      },
      {
        "id": 1203
      },
      {
        "id": 1204
      },
      {
        "id": 1205
      },
      {
        "id": 1206
      },
      {
        "id": 1207
      },
      {
        "id": 1208
      }
    ]
  },
  "frames": {
    "ver": 19,
    "vars": [
      {
        "id": 4
      },
      {
        "id": 3
      }
    ]
  },
  "genVal": 640,
  "guides": {
    "ver": 73,
    "vars": {
      "97": {
        "id": 97,
        "num": 2
      },
      "98": {
        "id": 98,
        "num": 1
      },
      "101": {
        "id": 101,
        "num": 16
      },
      "102": {
        "id": 102,
        "val": "[0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0]",
        "num": 18
      },
      "104": {
        "id": 104,
        "val": "0_1"
      },
      "105": {
        "id": 105,
        "val": "0_1"
      },
      "106": {
        "id": 106,
        "val": "0_1"
      },
      "107": {
        "id": 107,
        "val": "[0,0,0]",
        "num": 3
      },
      "110": {
        "id": 110,
        "val": "0_1"
      },
      "115": {
        "id": 115
      },
      "126": {
        "id": 126,
        "num": 1
      },
      "127": {
        "id": 127,
        "num": 1750348799
      },
      "143": {
        "id": 143,
        "val": "[{\"sid\":13,\"snum\":[0,0,0,0,0,0,0,0]}]"
      },
      "144": {
        "id": 144,
        "val": "[0,0,0,0]",
        "num": 4
      },
      "145": {
        "id": 145,
        "val": "[1,1]",
        "num": 2
      },
      "148": {
        "id": 148,
        "val": "[0,0,0,0,0]",
        "num": 5
      },
      "149": {
        "id": 149,
        "val": "[0]",
        "num": 1
      },
      "151": {
        "id": 151,
        "val": "[0,0,0,0,0,0]",
        "num": 6
      },
      "153": {
        "id": 153,
        "val": "{\"sid\":13,\"buy\":0,\"pop\":0,\"cards\":[],\"pid\":0}"
      },
      "154": {
        "id": 154,
        "val": "[0,0,0,0,0,0,0]",
        "num": 7
      },
      "155": {
        "id": 155,
        "val": "[1,1,1,0,0,0,0,0]",
        "num": 8
      },
      "215": {
        "id": 215,
        "num": 1755172869
      },
      "701": {
        "id": 701,
        "num": 261
      },
      "1102": {
        "id": 1102,
        "val": "[]"
      },
      "1109": {
        "id": 1109,
        "val": "0_0_0"
      },
      "1203": {
        "id": 1203,
        "val": "1761045750",
        "num": 1
      }
    }
  },
  "labels": {
    "ver": 7,
    "vars": []
  },
  "medals": {
    "ver": 7,
    "vars": []
  },
  "propVal": 446,
  "recentStatInfo": {},
  "sts": {
    "mus": 1,
    "aud": 1,
    "vib": 1,
    "qck": 1
  },
  "ver": 168
}`

		var jsonData map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &jsonData)
		if err != nil {
			fmt.Println("JSON 解析失败: ", err)
			panic(err)
		}
		requests[i] = jsonData
	}

	// 创建通道用于并发控制
	jobs := make(chan map[string]interface{}, numRequests)
	results := make(chan *testResult, numRequests)

	// 启动时间
	startTime := time.Now()

	// 启动工作协程
	for w := 0; w < concurrent; w++ {
		// go worker(w, url, tokenArr[w%len(tokenArr)], jobs, results)
		go worker(w, url, tokenArr[w], jobs, results)
	}

	// 发送任务到通道
	for _, req := range requests {
		jobs <- req
	}
	close(jobs)

	// 收集结果
	successCount := 0
	failCount := 0

	for i := 0; i < numRequests; i++ {
		result := <-results
		if result.Success {
			successCount++
			fmt.Printf("✓ 请求 %d 成功: %s\n", result.ID, result.Message)
		} else {
			failCount++
			fmt.Printf("✗ 请求 %d 失败: %s\n", result.ID, result.Error)
		}

		// 显示进度
		progress := float64(i+1) / float64(numRequests) * 100
		fmt.Printf("进度: %.1f%% [%d/%d]\r", progress, i+1, numRequests)
	}

	// 计算统计信息
	elapsed := time.Since(startTime)

	fmt.Println("\n=== 存档——压测结果 ===")
	fmt.Printf("总请求数: %d\n", numRequests)
	fmt.Printf("成功数量: %d\n", successCount)
	fmt.Printf("失败数量: %d\n", failCount)
	fmt.Printf("成功率: %.2f%%\n", float64(successCount)/float64(numRequests)*100)
	fmt.Printf("总耗时: %.2f 秒\n", elapsed.Seconds())
	fmt.Printf("平均响应时间: %.2f 毫秒\n", float64(elapsed.Milliseconds())/float64(numRequests))
	fmt.Printf("QPS: %.2f\n", float64(numRequests)/elapsed.Seconds())
}

package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/monaco-io/request"
)

func main() {
	fmt.Println("shop tools")
	itemid := "11721678694"
	keyword := url.QueryEscape("เครื่องพิมพ์ฉลาก")
	limit := 100
	newest := 0
	maxNewest := 1200 //PC60一页
	url := "https://th.xiapibuy.com/api/v4/search/search_items?by=relevancy&keyword=" + keyword + "&limit=" + strconv.Itoa(limit) + "&newest=" + strconv.Itoa(newest) + "&order=desc&page_type=search&scenario=PAGE_GLOBAL_SEARCH&version=2"

	for i := newest; i <= maxNewest; i = i + limit {
		time.Sleep(5)
		if handleContent(url, itemid, newest, keyword) {
			break
		}
	}

}

func handleContent(url string, itemid string, newest int, keyword string) bool {
	var result interface{}
	param := md5.Sum([]byte("by=relevancy&keyword=" + keyword + "&limit=100&newest=" + strconv.Itoa(newest) + "&order=desc&page_type=search&version=2"))
	paramMd5 := string(fmt.Sprintf("%x", param))
	s1 := md5.Sum([]byte("55b03" + paramMd5 + "55b03"))
	k := string(fmt.Sprintf("%x", s1))

	// fmt.Println(k)

	// $header  = array(
	// 	'if-none-match-: 55b03-'.$k,
	// 	'user-agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36'
	// );

	c := request.Client{
		URL:    url,
		Method: "GET",
		Header: map[string]string{"if-none-match-": "55b03-" + k, "user-agent": "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"},
	}
	resp := c.Send().Scan(&result)
	if !resp.OK() {
		// handle error
		fmt.Println(resp.Error())
	}
	res, _ := json.Marshal(result)
	strData := string(res)
	// fmt.Println(strData)
	// utils.WriteFile("test", strData)
	isContain := strings.Contains(strData, itemid)
	if isContain {
		regex := regexp.MustCompile(`"itemid":\d+,"label`)
		arr := regex.FindAllString(strData, -1)

		regex2 := regexp.MustCompile(`"image":"(.*?)"`)
		arr2 := regex2.FindAllStringSubmatch(strData, -1)

		rank := 0
		for i, v := range arr {
			if strings.Contains(v, itemid) {
				rank = newest + i + 1
				fmt.Println("排名：", rank)
				fmt.Println("图片：https://cf.shopee.co.th/file/" + arr2[i][1] + "_tn")
			}
		}
		fmt.Println("https://th.xiapibuy.com/search?keyword=" + keyword + "&page=" + strconv.Itoa((rank-1)/60))
		return true
	}
	return false
}

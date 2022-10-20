package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("start")
	basePath := "G:/work/test/"
	list := [][]string{{"aa", "bb", "https://qiniu.tangzhijiao.com/闲鱼小店_store_logo.jpg"},
		{"11", "cc", "https://qiniu.tangzhijiao.com/闲鱼小店_store_logo.jpg"}}
	tt := []string{"bb", "cc", "https://qiniu.tangzhijiao.com/闲鱼小店_store_logo.jpg"}
	list = append(list, tt)
	// url := "https://qiniu.tangzhijiao.com/闲鱼小店_store_logo.jpg"
	// dirName := "aaa"
	// rename := "bbb"
	for i := 0; i < len(list); i++ {
		dir := basePath + list[i][0]
		rename := list[i][1]
		url := list[i][2]
		err := checkAndCreateDir(dir)
		dir += "/"
		if err != nil {
			fmt.Println(err)
		}
		err = download(dir, url, rename)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func download(dir string, url string, rename string) error {
	fileName := ""
	if rename != "" {
		idx := strings.LastIndex(url, ".")
		fileName = rename + url[idx:]
	} else {
		idx := strings.LastIndex(url, "/")
		fileName = url[idx+1:]
	}

	path := dir + fileName

	//判断文件是否存在
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	file, err := http.Get(url)
	if err != nil {
		return err
	}
	defer file.Body.Close()
	content, err := ioutil.ReadAll(file.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, content, 0666)
	if err != nil {
		return err
	}
	return nil
}

func checkAndCreateDir(dir string) error {
	f, err := os.Stat(dir)
	if err == nil {
		if !f.IsDir() {
			//创建文件夹
			err = os.Mkdir(dir, 0666)
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	}
	if os.IsNotExist(err) {
		//创建文件夹
		err = os.Mkdir(dir, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

package utils

import (
	"fmt"
	"io"
	"os"
)

func WriteFile(filename string, contents string) {
	var f *os.File
	var err error
	if checkFileIsExist(filename) { //如果文件存在
		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err = os.Create(filename) //创建文件
		fmt.Println("文件不存在已创建")
	}
	if err != nil {
		fmt.Println("文件处理报错:", err)
		return
	}
	n, err := io.WriteString(f, contents) //写入文件，字符串
	if err != nil {
		fmt.Println("文件处理报错2:", err)
		return
	}
	fmt.Printf("写入 %d 个字符", n)
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

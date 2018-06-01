package main

import (
	"bufio"
	"io/ioutil"
	"io"
	"os"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
    "github.com/gomodule/redigo/redis"
    "github.com/mahonia"
)


type t struct{
    b string
}

func (b *t) Test (a string) string{
    return b.b + a
}

func main() {
    //test1()
    //write_file()
    read_file()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func write_file() {
    var writeString = "测试nnnn"
    var filename = "./output1.txt"
    var f *os.File
    var err1 error

    /*****************第一种方法 使用 io.WriteString 写入文件 *********************/
    if checkFileIsExist(filename) { //如果文件存在
        f, err1 = os.OpenFile(filename, os.O_APPEND, 0666)  //打开文件
        fmt.Println("文件存在")
    } else {
        f, err1 = os.Create(filename) //创建文件
        fmt.Println("文件不存在已创建")
    }
    CheckErr(err1)
    n, err1 := io.WriteString(f, writeString)   //写入文件，字符串
    CheckErr(err1)
    fmt.Printf("写入 %d 个字符n", n)

    /*****************************  第二种方式: 使用 ioutil.WriteFile 写入文件 ***********************************************/
	var d1 = []byte(writeString)
	err2 := ioutil.WriteFile("./output2.txt", d1, 0666) //写入文件(字节数组)
	check(err2)

	/*****************************  第三种方式:  使用 File(Write,WriteString) 写入文件 ***********************************************/
	f, err3 := os.Create("./output3.txt") //创建文件
	check(err3)
	defer f.Close()
	n2, err3 := f.Write(d1) //写入文件(字节数组)
	check(err3)
	fmt.Printf("写入 %d 个字节n", n2)
	n3, err3 := f.WriteString("writesn") //写入文件(字节数组)
	fmt.Printf("写入 %d 个字节n", n3)
	f.Sync()

	/***************************** 第四种方式:  使用 bufio.NewWriter 写入文件 ***********************************************/
	w := bufio.NewWriter(f) //创建新的 Writer 对象
	n4, err3 := w.WriteString("bufferedn")
	fmt.Printf("写入 %d 个字节n", n4)
	w.Flush()
	f.Close()
}

func checkFileIsExist(filename string) bool {
    var exist = true
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        exist = false
    }
    return exist
}

func read_file() {
    file, err := os.Open("./test.txt")
    if err != nil {
        fmt.Println(err)
    }

    //文件ex7.txt的编码是gb18030
    decoder := mahonia.NewDecoder("utf8")
    if decoder == nil {
        fmt.Println("编码不存在")
    }
    buf := make([]byte, 1024)
    for {
        len, _ := file.Read(buf)

        if len == 0 {
            break
        }

        fmt.Println(string(buf))
    }
    file.Close()
}

func test1(){
    var bb t
    bb.b = "sadfsadf"
    fmt.Println(bb.b)
    var ab = bb.Test("bbbbb")
    fmt.Println(ab)

    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8")
    if err != nil {
        fmt.Println(err)
    }

    _, err2 := db.Exec("INSERT INTO userinfo (username, departname, created) VALUES (?, ?, ?)","lily","销售","2016-06-21")
    CheckErr(err2)

    c, err := redis.Dial("tcp", "127.0.0.1:6379")
    CheckErr(err)

    defer c.Close()

    _, err = c.Do("SET", "mykey", "superWang")
    CheckErr(err)

    username, err := redis.String(c.Do("GET", "mykey"))
    if err != nil {
        fmt.Println("redis get failed:", err)
    } else {
        fmt.Printf("Get mykey: %v \n", username)
    }
}

func CheckErr(err error){
    if err != nil {
        fmt.Println(err)
    }
}
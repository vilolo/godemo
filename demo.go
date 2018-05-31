package main

import  "fmt"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"


type t struct{
    b string
}

func (b *t) Test (a string) string{
    return b.b + a
}

func main() {
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
}

func CheckErr(err error){
    if err != nil {
        fmt.Println(err)
    }
}
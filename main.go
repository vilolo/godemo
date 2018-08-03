package main

import (
	"strconv"
	"strings"
	"path/filepath"
	"path"
	//"bufio"
	"fmt"
	"os"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
)

func main() {
	list1 := getFileList()
	if list1 != nil {
		for i := range list1 {
			createHtml(list1[i])
		}
	}
	
}

func createHtml(file_name string){
	xlsx, err := excelize.OpenFile(file_name)
	if err != nil {
		fmt.Println(err)
		return
	}

	// cell := xlsx.GetCellValue("Sheet1", "B2")
	// fmt.Println(cell)

	var html string

	html = "<html>"

	sheet_index := xlsx.GetActiveSheetIndex()
	sheet_name := xlsx.GetSheetName(sheet_index)
	rows := xlsx.GetRows(sheet_name)

	// if len(rows) == 0 {
	// 	rows = xlsx.GetRows("批量上传模板(2)")
	// }
	temp := 1	//1=模板1，否则2
	for i, row := range rows {
		if i == 0 {
			if row[0] == "Product Unique ID" {
				temp = 2
			}
			continue
		}
		// for _, colCell := range row {
		// 	if colCell[0:3] == "http" {
		// 		html += "<img src='"+colCell+"' />"
		// 	}
		// }

		// for j, _ := range row {
		// 	if j == 2 {
		// 		html += "<hr>"
		// 		title := row[j]
		// 		html += "NO:"+(string)(i+1)+" >> "+"title: "+title
		// 	}
		// 	if j == 10 {
		// 		html += "<img src='"+row[j]+"' />"
		// 	}
		// }

		if temp == 1 {
			title := row[2]
			arr := []string{row[11],row[12],row[13],row[14],row[15],row[16],row[17],row[18]}
			html += "Line:"+strconv.Itoa(i+1)+" >> "+" SKU NO: <font color='red'>"+row[1]+"</font> >>> title: "+title+"<br>"
			for j := range arr {
				html += "<img src='"+arr[j]+"' width='200' />"
			}
			html += "<hr>"
		}else{
			title := row[2]
			arr := []string{row[22],row[23],row[24],row[25],row[26]}
			html += "Line:"+strconv.Itoa(i+1)+" >> "+" Product Unique ID: <font color='red'>"+row[0]+"</font> >>> title: "+title+"<br>"
			for j := range arr {
				html += "<img src='"+arr[j]+"' width='200' />"
			}
			html += "<hr>"
		}
	}

	nn := filepath.Base(file_name)
	
	var filenameWithSuffix string
	filenameWithSuffix = path.Base(nn)
	
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix)
	
	var filenameOnly string
    filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)

	html += "</html>"
	d1 := []byte(html)
    err1 := ioutil.WriteFile(filenameOnly+".html", d1, 0644)
    check(err1)
}

func getFileList() []string{
	dir, error := os.OpenFile("./", os.O_RDONLY, os.ModeDir)
	if error != nil {
		defer dir.Close()
		fmt.Println(error.Error())
		return nil
	}
	// fileinfo, _ := dir.Stat()
	// fmt.Println(fileinfo.IsDir())
	names, _ := dir.Readdir(-1)

	slice1 := make([]string, 2, 5)
	for _, name := range names {
		//fmt.Println(name.Name(), "目录?", name.IsDir())
		//判断文件后缀
		//fmt.Println(path.Ext(name.Name()))
		if path.Ext(name.Name()) == ".xlsx" {
			slice1 = append(slice1, name.Name())
		}
	}
	return slice1
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}
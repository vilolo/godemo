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
	test2()
}

func test1() {
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

func test2() {
	list := []string{"/uploads/article/timg.jpg",
		"/uploads/article/radial_machine.png",
		"/uploads/article/u=1666882005,4218080920&fm=26&gp=0.jpg",
		"/uploads/article/QQ截图20200813111156.png",
		"/uploads/article/axial_machine.jpg",
		"/uploads/article/conveyor1.jpg",
		"/uploads/article/loaderNunloader.jpg",
		"/uploads/article/mini_reject_conveyor.jpg",
		"/uploads/article/feeder.png",
		"/uploads/article/yamaha_nozzle4.jpg",
		"/uploads/article/samsung_nozzle1.jpg",
		"/uploads/logo/logo.png",
		"/uploads/img-desc/new_radial_banner.jpg",
		"/uploads/img-desc/AA_banner.jpg",
		"/uploads/img-desc/xiaoJUKI_nozzle2.png",
		"/uploads/img-desc/SMT.jpg",
		"/uploads/img-desc/field_machine2.jpg",
		"/uploads/img-desc/field_machine.jpg",
		"/uploads/img-desc/radial_feeder2.png",
		"/uploads/img-desc/W12c.png",
		"/uploads/img-desc/nozzles.jpg",
		"/uploads/img-desc/xiao_1head_&_4_guide_jaw_clipping.jpg",
		"/uploads/img-desc/小_DSC7373.png",
		"/uploads/img-desc/BHSA.jpg",
		"/uploads/img-desc/1xiao_DSC6150.jpg",
		"/uploads/img-desc/小banner.jpg",
		"/uploads/img-desc/service.png",
		"/uploads/img-desc/field_stack_feeder.jpg",
		"/uploads/img-desc/bowl_feeder.jpg",
		"/uploads/img-desc/JUKI_NOZZLES.jpg",
		"/uploads/img-desc/Gripper.jpg",
		"/uploads/img-desc/lable_feeder.png",
		"/uploads/img-desc/lable_feeder.png",
		"/uploads/img-desc/board_handling2.jpg",
		"/uploads/img-desc/radial_feeder2.png",
		"/uploads/img-desc/axial_machine_videocover.jpg",
		"/uploads/img-desc/conveyor1.jpg"}
	for i := 0; i < len(list); i++ {
		idx := strings.LastIndex(list[i], "/")
		dir := "G:/work/test/wen" + list[i][:idx]
		err := checkAndCreateDir(dir)
		if err != nil {
			fmt.Println(err)
		}
		dir += "/"
		download(dir, "https://adm.smtfan.com"+list[i], "")
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

package main

import (
	"fmt"
	"os"
	"time"
)

var watchDir []string = []string{}

type Dir struct {
	ParentNode string
	ChildNodes []string
	Name       string
	code       string
}

func getLastTime() string {

}

func getFileInfo(path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fileInfo.ModTime())
	now := time.Now()
	fmt.Println("-->", now.Before(fileInfo.ModTime()))
	//fmt.Println(fileInfo.Mode())
}

//监测到没有文件夹信息的文件(1.新建文件;2.配置文件丢失)
//扫描监测目录下内容.建立数据结构
func NewStart(dir string) {

}

//检测监测文件夹有无更改
func checkRoofDir(watchDir []string) {
	//监测文件夹信息记录文件
	for _, dir := range watchDir {
		if os.IsNotExist(dir + "/init.ame") {
			fmt.Println("file not exist")

		}
		fileInfo, err := os.Stat(dir)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//检测文件更改状态
func checkModify() bool {

}

func main() {
	getFileInfo("newTest.go")
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type Dir struct {
	Id         int64
	ParentId   int64
	Path       string
	Name       string
	Children   []Dir
	IsDir      bool
	ModifyTime time.Time
}

var pathSep string = string(os.PathSeparator)

var data syscall.Win32finddata

func getFileId(path string) int64 {
	uintpath := syscall.StringToUTF16Ptr(path)
	handle, err := syscall.FindFirstFile(uintpath, &data)
	defer syscall.CloseHandle(handle)
	if err != nil {
		fmt.Println("--->", err)
	}
	return data.CreationTime.Nanoseconds()
}

func pathWalk(path string, parentId int64) []Dir {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	var children []Dir
	var childPath string
	for _, dir := range dirs {
		childPath = path + pathSep + dir.Name()
		children = append(children, getFileInfo(childPath, parentId))
	}
	return children
}

func getFileInfo(path string, parentId int64) Dir {
	fileInfo, errs := os.Stat(path)
	if errs != nil {
		fmt.Println(errs)
	}
	var self Dir
	self = Dir{
		Id:         getFileId(path),
		ParentId:   parentId,
		Path:       filepath.Dir(path),
		Name:       fileInfo.Name(),
		IsDir:      fileInfo.IsDir(),
		ModifyTime: fileInfo.ModTime(),
	}
	if self.IsDir {
		self.Children = pathWalk(path, self.Id)
	}
	return self
}

func main() {
	all := getFileInfo("./tt", -1)
	draw(all)
}

func draw(dir Dir) {
	var jsonbyte []byte
	var err error

	jsonbyte, err = json.Marshal(dir)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonbyte))
}

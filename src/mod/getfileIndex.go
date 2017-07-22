package main

import (
	"fmt"
	"syscall"
)

//获取文件在磁盘分区上的唯一索引值
func getFileIndex(path string) string {
	var fileInfo syscall.ByHandleFileInformation
	handle, err := syscall.Open(path, syscall.O_RDONLY, uint32(32))
	if err != nil {
		fmt.Println(err)
	}
	err = syscall.GetFileInformationByHandle(handle, &fileInfo)
	syscall.CloseHandle(handle)
	if err != nil {
		fmt.Println(err)
		return "-1"
	}

	return fmt.Sprintf("%d%d", fileInfo.FileIndexHigh, fileInfo.FileIndexLow)
}

func main() {
	fmt.Println(getFileIndex("F:/newyu.txt"))
}

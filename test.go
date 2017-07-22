package main

import (
	"fmt"
	//"strconv"
	//"sys"
	"syscall"
	//"time"
)

func main() {

	/*	//var handle syscall.Handle
		//var data syscall.Win32finddata
		s := "C:\\test.txt"
		//us, err := syscall.UTF16PtrFromString(s)
		//if err != nil {
		//	fmt.Println(err)
		//}
		var fileinfo syscall.ByHandleFileInformation
		//handle, _ := syscall.FindFirstFile(us, &data)
		//var sa syscall.SecurityAttributes
		handle, erro := syscall.Open(s, syscall.O_RDONLY, uint32(32)) //syscall.CreateFile(us, syscall.GENERIC_READ, syscall.FILE_SHARE_READ, &sa, syscall.OPEN_ALWAYS, 0, 0)
		if erro != nil {
			fmt.Println(erro)
		}
		err := syscall.GetFileInformationByHandle(handle, &fileinfo)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("modify:")
		fmt.Println(s + ":")
		fmt.Println("highIndex:", fileinfo.FileIndexHigh)
		fmt.Println("lowIndex:", fileinfo.FileIndexLow)
		err = syscall.Close(handle)
		if err != nil {
			fmt.Println(err)
		}
	*/

	path := "tt"

	uintpath := syscall.StringToUTF16Ptr(path)
	var data syscall.Win32finddata
	_, errs := syscall.FindFirstFile(uintpath, &data)
	if errs != nil {
		fmt.Println("--->", errs)
	}
	fmt.Println(data.CreationTime.Nanoseconds())
	//handle, err := syscall.Open(path, syscall.O_RDONLY, uint32(32))
	/*var sa syscall.SecurityAttributes
	handle, err := syscall.CreateFile(uintpath, syscall.GENERIC_ALL, syscall.FILE_SHARE_READ, &sa, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		fmt.Println(err)
	}*/
	/*tmptime := syscall.NsecToFiletime(int64(3000000000))
	syscall.Filetime.
		err = syscall.SetFileTime(handle, &data.LastWriteTime, &tmptime, &data.LastWriteTime)
	fmt.Println(data.LastWriteTime)*/
	//var filetime syscall.NsecToFiletime(nsec)
	/*
		if err != nil {
			fmt.Println(err)
		}*/

}

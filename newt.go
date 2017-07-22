package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	err := ioutil.WriteFile("./test/h.txt", []byte("data"), 0666)
	fmt.Println(err)
}

package main

import (
	"bytes"
	"fmt"
	//"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("tree")
	cmd.Path = "H"
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(out.String())
}

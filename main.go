package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("git", "ls-files")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Errorf("ERROR on git ls-files! %v", err))
		return
	}
	fmt.Printf("Out is:\n%s", out)
}

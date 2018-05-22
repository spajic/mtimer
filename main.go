package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/djherbis/times"
)

func main() {
	cmd := exec.Command("git", "ls-files")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Errorf("ERROR on git ls-files! %v", err))
		return
	}
	fmt.Printf("Out is:\n%s", out)

	t, err := times.Stat("main.go")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("main.go mtime: ", t.ModTime())

	mtime := time.Date(2006, time.February, 1, 3, 4, 5, 0, time.UTC)
	atime := mtime
	if err := os.Chtimes("main.go", atime, mtime); err != nil {
		log.Fatal(err)
	}
}

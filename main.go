package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/djherbis/times"
)

func main() {
	cmd := exec.Command("git", "ls-files")
	files, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Errorf("ERROR on git ls-files! %v", err))
		return
	}

	out, err := os.Create("mtimer.dat")
	if err != nil {
		panic("Error opening file")
	}
	defer out.Close()

	scanner := bufio.NewScanner(bytes.NewReader(files))
	for scanner.Scan() {
		fileName := scanner.Text()
		t, err := times.Stat(fileName)
		if err != nil {
			fmt.Println("WARNING: ", err.Error())
			continue
		}
		out.WriteString(fileName + "\n")
		mtime := t.ModTime()
		// if err := os.Chtimes(fileName, mtime, mtime); err != nil {
		// 	log.Fatal(err)
		// }
		out.WriteString(mtime.String() + "\n")
	}
}

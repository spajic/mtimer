package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/djherbis/times"
)

var store bool
var apply bool

func init() {
	flag.BoolVar(&store, "store", false, "-store to store mtimes")
	flag.BoolVar(&apply, "apply", false, "-apply to apply stored mtimes from mtimer.dat")
}

func main() {
	flag.Parse()
	if store && !apply {
		fmt.Println("STORE MODE")
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
			out.WriteString(mtime.String() + "\n")
		}
	} else if apply {
		fmt.Println("APPLY MODE")
		fileHandle, err := os.Open("mtimer.dat")
		if err != nil {
			fmt.Println("CANT OPEN FILE mtimer.dat! FINISH")
			return
		}
		defer fileHandle.Close()
		fileScanner := bufio.NewScanner(fileHandle)
		for fileScanner.Scan() {
			fileName := fileScanner.Text()
			fileScanner.Scan()
			fileMtimeText := fileScanner.Text()
			fileMtime, err := time.Parse("2018-05-22 14:32:08 +0300 MSK", fileMtimeText)
			if err != nil {
				fmt.Println("WARNING: Can't parse mtime", fileMtimeText, "for file", fileName)
				continue
			}

			if err := os.Chtimes(fileName, fileMtime, fileMtime); err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		fmt.Println("Use with -store or -apply flag")
	}
}

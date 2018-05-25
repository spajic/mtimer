package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/djherbis/times"
)

var path string
var store bool
var apply bool

func init() {
	flag.StringVar(&path, "path", "", "path to folder with files")
	flag.BoolVar(&store, "store", false, "to store mtimes")
	flag.BoolVar(&apply, "apply", false, "to apply stored mtimes from mtimer.dat")
}

func main() {
	flag.Parse()
	if path == "" {
		fmt.Println("PATH NOT SPECIFIED! EXIT NOW!")
		return
	}
	if store && !apply {
		fmt.Println("STORE MODE")
		fmt.Println("Working with path =", path)
		cmd := exec.Command("find", ".", "-and", "-not", "-path", "./node_modules*", "-and", "-not", "-path", "./tmp*")
		cmd.Dir = path
		files, err := cmd.Output()
		if err != nil {
			fmt.Println(fmt.Errorf("ERROR on list files. EXIT NOW, %v", err))
			return
		}

		pathToMtimerDat := path + "/mtimer.dat"
		fmt.Println("Create file", pathToMtimerDat)
		out, err := os.Create(pathToMtimerDat)
		if err != nil {
			panic("Error creating file " + pathToMtimerDat + ". EXIT NOW")
		}
		defer out.Close()

		scanner := bufio.NewScanner(bytes.NewReader(files))
		for scanner.Scan() {
			fileName := strings.Replace(scanner.Text(), ".", "", 1)
			t, err := times.Stat(path + fileName)
			if err != nil {
				fmt.Println("WARNING: ", err.Error())
				continue
			}
			out.WriteString(fileName + "\n")
			mtime := t.ModTime()
			out.WriteString(mtime.Format(string(time.RFC3339)) + "\n")
		}
		fmt.Println("FINISHED SUCCESSFULLY")
	} else if apply {
		fmt.Println("APPLY MODE")
		fileHandle, err := os.Open("/busfor/releases/building/mtimer.dat")
		if err != nil {
			fmt.Println("CANT OPEN FILE mtimer.dat! FINISH")
			return
		}
		defer fileHandle.Close()
		fileScanner := bufio.NewScanner(fileHandle)
		for fileScanner.Scan() {
			fileName := "/busfor/releases/building/" + fileScanner.Text()
			fileScanner.Scan()
			fileMtimeText := fileScanner.Text()
			fileMtime, err := time.Parse(time.RFC3339, fileMtimeText)
			if err != nil {
				fmt.Println("WARNING: Can't parse mtime", fileMtimeText, "for file", fileName, err)
				continue
			}

			if err := os.Chtimes(fileName, fileMtime, fileMtime); err != nil {
				fmt.Println("WARNING:", err)
				continue
			}
		}
		fmt.Println("FINISHED SUCCESSFULLY")
	} else {
		fmt.Println("Use with -store or -apply flag")
	}
}

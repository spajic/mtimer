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

var filespath string
var timespath string
var store bool
var apply bool

func init() {
	flag.StringVar(&filespath, "filespath", "", "path to folder with files")
	flag.StringVar(&timespath, "timespath", "", "path to folder with mtimer.dat")
	flag.BoolVar(&store, "store", false, "to store mtimes")
	flag.BoolVar(&apply, "apply", false, "to apply stored mtimes from mtimer.dat")
}

func main() {
	flag.Parse()
	if filespath == "" {
		fmt.Println("filespath NOT SPECIFIED! EXIT NOW!")
		return
	}
	if timespath == "" {
		fmt.Println("timespath NOT SPECIFIED! EXIT NOW!")
		return
	}
	if store && !apply {
		fmt.Println("STORE MODE")
		fmt.Println("Working with path =", filespath)
		cmd := exec.Command("find", ".", "-not", "-path", "./node_modules*", "-and", "-not", "-path", "./tmp*")
		cmd.Dir = filespath
		files, err := cmd.Output()
		if err != nil {
			fmt.Println(fmt.Errorf("ERROR on list files. EXIT NOW, %v", err))
			return
		}

		pathToMtimerDat := timespath + "/mtimer.dat"
		fmt.Println("Create file", pathToMtimerDat)
		out, err := os.Create(pathToMtimerDat)
		if err != nil {
			panic("Error creating file " + pathToMtimerDat + ". EXIT NOW")
		}
		defer out.Close()

		readFiles := 0
		scanner := bufio.NewScanner(bytes.NewReader(files))
		for scanner.Scan() {
			fileName := strings.Replace(scanner.Text(), ".", "", 1)
			t, err := times.Stat(filespath + fileName)
			if err != nil {
				fmt.Println("WARNING: ", err.Error())
				continue
			}
			out.WriteString(fileName + "\n")
			mtime := t.ModTime()
			out.WriteString(mtime.Format(string(time.RFC3339)) + "\n")
			readFiles++
		}
		fmt.Println("Stored mtimes of", readFiles, "files")
		fmt.Println("FINISHED SUCCESSFULLY")
	} else if apply {
		fmt.Println("APPLY MODE")
		pathToMtimerDat := timespath + "/mtimer.dat"
		fmt.Println("Applying mtimes from", pathToMtimerDat)
		fmt.Println("Updating files in", filespath)
		fileHandle, err := os.Open(pathToMtimerDat)
		if err != nil {
			fmt.Println("CANT OPEN FILE mtimer.dat! FINISH")
			return
		}
		defer fileHandle.Close()
		fileScanner := bufio.NewScanner(fileHandle)
		updatedCount := 0
		for fileScanner.Scan() {
			fileName := filespath + fileScanner.Text()
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
			updatedCount++
		}
		fmt.Println("Updated mtimes of", updatedCount, "files")
		fmt.Println("FINISHED SUCCESSFULLY")
	} else {
		fmt.Println("Use with ---store or --apply flag")
	}
}

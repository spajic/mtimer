package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/djherbis/times"
)

const version = "0.1.0"

var filespath string
var timespath string
var ignoreFolders string
var store bool
var apply bool
var showVersion bool
var showHelp bool

const help = `mtimer cat store and apply mtimes of file in given directory

Usage examples:
"mtimer --store --filespath=/path/to/files --timespath=/path/to_mtimer_dat" - store mtimes of file from filespath to mtimer.dat
"mtimer --store --ignore=node_modules,tmp,.git --filespath=/path --timespath=/path" - ignore specified subfolders
"mtimer --apply --filespath=/path/to/files --timespath=/path/to_mtimer_dat" - apply mtimes from mtimer.dat to files in filespath
"mtimer --version" - show version and exit
"mtimer --help" - show this help
`

func init() {
	flag.StringVar(&filespath, "filespath", "", "path to folder with files")
	flag.StringVar(&timespath, "timespath", "", "path to folder with mtimer.dat")
	flag.StringVar(&ignoreFolders, "ignore", "", "ignore folders")
	flag.BoolVar(&store, "store", false, "store mtimes to mtimer.dat")
	flag.BoolVar(&apply, "apply", false, "apply stored mtimes from mtimer.dat")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showHelp, "help", false, "show help")
}

func main() {
	flag.Parse()
	if showHelp {
		showHelpAndExit()
	}
	if showVersion {
		showVersionAndExit()
	}
	if filespath == "" || timespath == "" {
		showErrorMessageAndExit("Need to specify --filespath and --timespath.")
	}
	if !store && !apply || (store && apply) {
		showErrorMessageAndExit("Need to specify --store or --apply mode.")
	}

	logStart()
	if store {
		storeMtimes()
	} else {
		applyMtimes()
	}
}

func showHelpAndExit() {
	fmt.Println(help)
	os.Exit(0)
}
func showVersionAndExit() {
	fmt.Println(version)
	os.Exit(0)
}
func showErrorMessageAndExit(message string) {
	fmt.Println(message, "Exit now. Call 'mtimer --help' for help.")
	os.Exit(1)
}

func checkErrOrExitWithMessage(err error, msg string) {
	if err == nil {
		return
	}
	fmt.Println(fmt.Errorf(msg+". Fatal error: %v. Exit now.", err))
	os.Exit(1)
}

func logStart() {
	var modeString string
	if store {
		modeString = "store"
	} else {
		modeString = "apply"
	}
	fmt.Println("Start mtimer in", modeString, "mode")
	fmt.Println("filespath =", filespath)
	fmt.Println("timespath =", timespath)
	fmt.Println("ignoreFolders =", ignoreFolders)
}

func pathToMtimerDat() string {
	return timespath + "/mtimer.dat"
}

// for flag ignoreFolders="first_folder,second_folder"
// returns []string like [".", "-not", "-path", "./first_folder*", "-and", "-not", "-path", "./second_folder", "-type", "f"]
func findArgs() []string {
	result := []string{"."}
	if ignoreFolders != "" {
		folders := strings.Split(ignoreFolders, ",")
		for i, folder := range folders {
			result = append(result, "-not")
			result = append(result, "-path")
			result = append(result, "./"+folder+"*")
			if i < (len(folders) - 1) {
				result = append(result, "-and")
			}
		}
	}
	result = append(result, "-type")
	result = append(result, "f")
	return result
}

func fileHash(path string) (string, error) {
	fileHandle, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("mtimer warning: failed to open file %s for hashsum: %v", path, err)
		return "ERROR", err
	}
	defer fileHandle.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, fileHandle); err != nil {
		err = fmt.Errorf("mtimer warning: problem calculatiing hash for %s: %v", path, err)
		return "ERROR", err
	}
	return hex.EncodeToString(hash.Sum(nil)[:16]), nil
}

func storeMtimes() {
	listFilesCmd := exec.Command("find", findArgs()...)
	fmt.Println("List files with find", findArgs())
	listFilesCmd.Dir = filespath
	files, err := listFilesCmd.Output()
	checkErrOrExitWithMessage(err, "Error in getting files list")

	fmt.Println("Create file", pathToMtimerDat())
	out, err := os.Create(pathToMtimerDat())
	checkErrOrExitWithMessage(err, "Error creating mtimer.dat file")
	defer out.Close()

	readFilesCount := 0
	scanner := bufio.NewScanner(bytes.NewReader(files))
	for scanner.Scan() {
		// TODO: obtain convenient output from find command directly
		// Change ./file/name to /file/name
		fileName := strings.Replace(scanner.Text(), ".", "", 1)

		t, err := times.Stat(filespath + fileName)
		if err != nil {
			fmt.Println("mtimer warning: ", err.Error())
			continue
		}

		hash, err := fileHash(filespath + fileName)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// write filename
		out.WriteString(fileName + "\n")

		// write mtime
		mtime := t.ModTime()
		out.WriteString(mtime.Format(string(time.RFC3339)) + "\n")

		// write hash
		out.WriteString(hash + "\n")

		readFilesCount++
	}
	fmt.Println("Successfully stored mtimes of", readFilesCount, "files")
}

// If hash of file have changed - do not update mtime
func applyMtimes() {
	fileHandle, err := os.Open(pathToMtimerDat())
	checkErrOrExitWithMessage(err, "Error opening mtimer.dat file")
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)
	updatedFilesCount := 0
	for fileScanner.Scan() {
		fileName := filespath + fileScanner.Text()

		fileScanner.Scan()
		fileMtimeText := fileScanner.Text()

		fileScanner.Scan()
		storedHash := fileScanner.Text()

		fileMtime, err := time.Parse(time.RFC3339, fileMtimeText)
		if err != nil {
			fmt.Println("mtimer warning: can't parse mtime", fileMtimeText, "for file", fileName, err)
			continue
		}

		hash, err := fileHash(fileName)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if storedHash != hash {
			fmt.Printf("mtimer info: file %s changed, not updating mtime\n", filespath+fileName)
			continue
		}

		if err := os.Chtimes(fileName, fileMtime, fileMtime); err != nil {
			fmt.Println("mtimer warning: can't update mtime", err)
			continue
		}
		updatedFilesCount++
	}
	fmt.Println("Successfully Updated mtimes of", updatedFilesCount, "files")
}

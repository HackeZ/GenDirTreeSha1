package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/hackez/GenDirTreeSha1/FileSHACount"
)

var igFile = flag.String("f", "", "-f used to set ignore file list.")
var igDir = flag.String("d", "", "-d used to set ignore dir list.")
var dirRoot = flag.String("r", "", "-r used to set dir root.")
var maxGoroutineNum = flag.Int64("g", 1024, "-g used to set max of running goroutine number.")

var ignoreFileList []string
var ignoreDirList []string

// init get -r -f -d Params.
func init() {

	// 添加多核支持，适合当前 CPU 计算密集的场景。
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	// init ignore dir and file list.
	ignoreFileList = make([]string, 0)
	ignoreDirList = make([]string, 0)

	if len(*igFile) != 0 {
		for _, v := range strings.Split(*igFile, ",") {
			ignoreFileList = append(ignoreFileList, v)
		}
		fmt.Println("ignoreFileList -->", ignoreFileList)
	}
	if len(*igDir) != 0 {
		for _, v := range strings.Split(*igDir, ",") {
			ignoreDirList = append(ignoreDirList, v)
		}
		fmt.Println("ignoreDirList -->", ignoreDirList)
	}

}

func main() {
	// If dir root not set, exit.
	if len(*dirRoot) == 0 {
		fmt.Println("Please Set Up Dir Root.")
		os.Exit(0)
	}

	// Open result.txt
	file, err := os.OpenFile("./result.txt", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	// Clear Up result.txt
	file.Truncate(0)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	t1 := time.Now()
	// Get Result and Save it to writeMap.
	writeMap, err := FileSHACount.GenDirTreeSHA1(*dirRoot, ignoreDirList, ignoreFileList, *maxGoroutineNum)
	// Handle GenDirTreeSHA1 Error.
	if err != nil {
		fmt.Println(err)
		return
	}
	t2 := time.Now()

	// Write Result into File.
	for sha1, f := range writeMap {
		file1 := *f
		_, err = file.WriteString(fmt.Sprintf("%s, %s, %d Byte\n", file1.Name(), sha1, file1.Size()))
		if err != nil {
			panic(err)
		}
	}

	// Done.
	fmt.Println("Generate Dir tree SHA1 Done, Check your result.txt!")
	fmt.Println("Used Time: ", t2.Sub(t1))
}

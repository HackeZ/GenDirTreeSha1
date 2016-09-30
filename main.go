package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/hackez/gendirtreesha1/FileSHACount"
)

const (
	version = "1.3.1"
)

var (
	igFile          = kingpin.Flag("ignoreFile", "set up ignore file list. split by ',' and support wildcards.").Short('f').Default("").String()
	igDir           = kingpin.Flag("ignoreDir", "set up ignore directory list. split by ',' and support wildcards.").Short('d').Default("").String()
	dirRoot         = kingpin.Flag("root", "set up directory root.").Short('r').Required().String()
	maxGoroutineNum = kingpin.Flag("maxG", "set up max of running goroutine number.").Short('g').Default("2048").Int64()
)

var ignoreFileList []string
var ignoreDirList []string

// init get -r -f -d Params.
func init() {

	// 添加多核支持，适合当前 CPU 计算密集的场景。
	runtime.GOMAXPROCS(runtime.NumCPU())

	kingpin.Version(version)
	kingpin.Parse()

	// init ignore dir and file list.
	ignoreFileList = make([]string, 0)
	ignoreDirList = make([]string, 0)

	if len(*igFile) != 0 {
		for _, v := range strings.Split(*igFile, ",") {
			ignoreFileList = append(ignoreFileList, v)
		}
		fmt.Println("ignoreFileList =>", ignoreFileList)
	}
	if len(*igDir) != 0 {
		for _, v := range strings.Split(*igDir, ",") {
			ignoreDirList = append(ignoreDirList, v)
		}
		fmt.Println("ignoreDirList =>", ignoreDirList)
	}

}

func main() {
	// If dir root not set, exit.
	if len(*dirRoot) == 0 {
		fmt.Println("Please Use \"-r\" to Set Up Dir Root.")
		os.Exit(0)
	}
	fmt.Println("Max Goroutine:", *maxGoroutineNum)

	// Open result.txt
	file, err := os.OpenFile("./result.txt", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	// Clear Up result.txt
	file.Truncate(0)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	t1 := time.Now()
	// Get Result and Save it to writeChan.
	writeChan, err := FileSHACount.GenDirTreeSHA1(*dirRoot, ignoreDirList, ignoreFileList, *maxGoroutineNum)
	// Handle GenDirTreeSHA1 Error.
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write Result into File.
	for {
		sf, isOpen := <-writeChan
		if !isOpen {
			break
		}
		f := *sf.FileInfo
		_, err = file.WriteString(fmt.Sprintf("%s, %s, %d Byte\n", f.Name(), sf.SHA1, f.Size()))
		if err != nil {
			panic(err)
		}
	}

	t2 := time.Now()
	// Done.
	fmt.Println("Generate Dir tree SHA1 Done, Check your result.txt!")
	fmt.Println("Used Time: ", t2.Sub(t1))
}

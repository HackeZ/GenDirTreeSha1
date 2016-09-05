package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"GenDirTreeSha1/FileSHACount"
)

var igFile = flag.String("f", "", "-f used to set ignore file list.")
var igDir = flag.String("d", "", "-d used to set ignore dir list.")
var dirRoot = flag.String("r", "", "-r used to set dir root.")

var ignoreFileList []string
var ignoreDirList []string

// init get -r -f -d Params.
func init() {
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

	FileSHACount.GenDirTreeSHA1(*dirRoot, ignoreDirList, ignoreFileList)

	// Done.
	fmt.Println("Generate Dir tree SHA1 Done, Check your result.txt!")
}

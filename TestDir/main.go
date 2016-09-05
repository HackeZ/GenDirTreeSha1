package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"GenDirTreeSha1/glob"
)

// 单线程版本

var igFile *string = flag.String("f", "", "-f used to set ignore file list.")
var igDir *string = flag.String("d", "", "-d used to set ignore dir list.")
var dirRoot *string = flag.String("r", "", "-r used to set dir root.")

var ignoreFileList []string
var ignoreDirList []string

// GetDirTree return dir tree.
func GetDirTree(dirRoot string, file *os.File) {
	err := filepath.Walk(dirRoot, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		// Compile Dir Name.
		if f.IsDir() {
			for _, value := range ignoreDirList {
				if ok := glob.Match(value, []byte(f.Name()), false); ok {
					// Compiled! Return nil.
					return filepath.SkipDir
				}
			}
			return nil
		}
		// Compile File Name
		for _, value := range ignoreFileList {
			if ok := glob.Match(value, []byte(f.Name()), false); ok {
				// Compiled! Return nil.
				return nil
			}
		}
		sha, buf := sha1.New(), make([]byte, 1024*1024*16*16)
		thisF, _ := os.Open(path)
		defer thisF.Close()
		for {
			n, err := thisF.Read(buf)
			sha.Write(buf[:n])
			if err == io.EOF {
				break
			}
		}
		_, err = file.WriteString(fmt.Sprintf("%s, %x, %d Byte\n", f.Name(), sha.Sum(nil), f.Size()))
		if err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func init() {
	flag.Parse()

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
	if len(*dirRoot) == 0 {
		fmt.Println("Please Set Up Dir Root.")
		os.Exit(0)
	}
	file, err := os.OpenFile("./result.txt", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	file.Truncate(0)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	GetDirTree(*dirRoot, file)
}

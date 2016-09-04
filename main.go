package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"GenDirTreeSha1/glob"
)

var igFile *string = flag.String("f", "", "-f used to set ignore file list.")
var igDir *string = flag.String("d", "", "-d used to set ignore dir list.")
var dirRoot *string = flag.String("r", "", "-r used to set dir root.")

var ignoreFileList []string
var ignoreDirList []string

var writeFileChan chan map[string]*os.FileInfo
var wg sync.WaitGroup

// GetDirTree return dir tree.
func GetDirTree(dirRoot string, resultFile *os.File) {
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
		go Sha1Sum(path, f)

		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	// When Walk Done, Close Chan to Exit Main Process.
	wg.Wait()
	close(writeFileChan)
}

func Sha1Sum(path string, file os.FileInfo) {
	wg.Add(1)
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
	fileSha1 := sha.Sum(nil)
	fileString := fmt.Sprintf("%x", fileSha1)
	sendFileInfo := make(map[string]*os.FileInfo)
	sendFileInfo[fileString] = &file
	writeFileChan <- sendFileInfo
	wg.Done()
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

	writeFileChan = make(chan map[string]*os.FileInfo, 100)
	writeMap := make(map[string]*os.FileInfo)
	GetDirTree(*dirRoot, file)

	time.Sleep(10 * time.Millisecond)

	for {
		result, isClose := <-writeFileChan
		if !isClose {
			break
		}
		for sha1, f := range result {
			fmt.Println(result)
			writeMap[sha1] = f
		}
	}

	fmt.Println(writeMap)

	// Write Result into File.
	for sha1, f := range writeMap {
		file1 := *f
		_, err = file.WriteString(fmt.Sprintf("%s, %s, %d Byte\n", file1.Name(), sha1, file1.Size()))
		if err != nil {
			panic(err)
		}
	}
}

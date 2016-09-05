package FileSHACount

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"GenDirTreeSha1/glob"
)

var ignoreFileList []string
var ignoreDirList []string

var writeFileChan chan map[string]*os.FileInfo
var wg sync.WaitGroup

// GenDirTreeSHA1
func GenDirTreeSHA1(path string, ignoreDir, ignoreFile []string) {
	ignoreDirList, ignoreFileList = ignoreDir, ignoreFile
	// Open result.txt
	file, err := os.OpenFile("./result.txt", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	// Clear Up result.txt
	file.Truncate(0)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	// Send Count File SHA1 between gorotinue.
	writeFileChan = make(chan map[string]*os.FileInfo, 100)
	// writeMap Save Result from writeFileChan
	writeMap := make(map[string]*os.FileInfo)
	// Get Dir Tree Start.
	getDirTree(path, file)

	// Main Process Sleep a little, make writeFileChan not empty.
	time.Sleep(10 * time.Millisecond)

	// Get Result Start.
	for {
		result, isOpen := <-writeFileChan
		if !isOpen {
			break
		}
		for sha1, f := range result {
			writeMap[sha1] = f
		}
	}

	// Write Result into File.
	for sha1, f := range writeMap {
		file1 := *f
		_, err = file.WriteString(fmt.Sprintf("%s, %s, %d Byte\n", file1.Name(), sha1, file1.Size()))
		if err != nil {
			panic(err)
		}
	}
}

// getDirTree return dir tree.
func getDirTree(dirRoot string, resultFile *os.File) {
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

		// Get a legal File to Count SHA1.
		go sHA1Sum(path, f)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

	// When Walk Done, Close Chan to Exit Main Process After Count SHA1 gorotinue Done.
	wg.Wait()
	close(writeFileChan)
}

// sHA1Sum Get File info and Count SHA1.
func sHA1Sum(path string, file os.FileInfo) {
	// Start a New gorotinue.
	wg.Add(1)

	// Read File Content.
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

	// Use sha Count File SHA1.
	fileSha1 := sha.Sum(nil)
	fileString := fmt.Sprintf("%x", fileSha1)

	// init a sendFileInfo to Send writeFileChan.
	sendFileInfo := make(map[string]*os.FileInfo)
	sendFileInfo[fileString] = &file
	// Send to writeFileChan.
	writeFileChan <- sendFileInfo
	// release sendFileInfo.
	sendFileInfo = nil

	// This gorotinue Done.
	wg.Done()
}

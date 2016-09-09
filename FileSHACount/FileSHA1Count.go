package FileSHACount

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hackez/GenDirTreeSha1/glob"
)

var ignoreFileList []string
var ignoreDirList []string

// Set Max of Gorotine Number.
var maxGoroutineChan chan int

var writeFileChan chan map[string]*os.FileInfo
var wg sync.WaitGroup

// GenDirTreeSHA1 ...
// @param path, ignoreDir, ignoreFile, maxGoroutineNum
// @return resultMap, error
func GenDirTreeSHA1(path string, ignoreDir, ignoreFile []string, maxGNum int64) (map[string]*os.FileInfo, error) {

	// set of Max Gorotine Number.
	if maxGNum <= 0 {
		return nil, errors.New("Max Gorotine Number must greater than 0")
	}
	if maxGNum >= 100 {
		return nil, errors.New("Max Gorotine Number must less than 100")
	}
	maxGoroutineChan = make(chan int, maxGNum)

	// set of ignoreDirList and ignoreFileList.
	ignoreDirList, ignoreFileList = ignoreDir, ignoreFile

	// Send Count File SHA1 between gorotinue.
	writeFileChan = make(chan map[string]*os.FileInfo, 100)
	// writeMap Save Result from writeFileChan
	writeMap := make(map[string]*os.FileInfo)
	// Get Dir Tree Start.
	go getDirTree(path)

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

	return writeMap, nil
}

// getDirTree return dir tree.
func getDirTree(dirRoot string) {
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

		// check G chan full or not.
		maxGoroutineChan <- 1
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
	close(maxGoroutineChan)
}

// sHA1Sum Get File info and Count SHA1.
func sHA1Sum(path string, file os.FileInfo) {
	// Start a New gorotinue.
	wg.Add(1)

	// Read File Content.
	sha, buf := sha1.New(), make([]byte, 1024*16)
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
	fileString := fmt.Sprintf("%x", sha.Sum(nil))

	// init a sendFileInfo to Send writeFileChan.
	sendFileInfo := make(map[string]*os.FileInfo)
	sendFileInfo[fileString] = &file
	// Send to writeFileChan.
	writeFileChan <- sendFileInfo
	// release sendFileInfo.
	sendFileInfo = nil

	// This gorotinue Done.
	<-maxGoroutineChan
	wg.Done()
}

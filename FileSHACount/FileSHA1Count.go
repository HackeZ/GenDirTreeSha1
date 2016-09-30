package FileSHACount

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/hackez/gendirtreesha1/glob"
)

var ignoreFileList []string
var ignoreDirList []string

// Set Max of Gorotine Number.
var maxGoroutineChan chan int

var writeFileChan chan SHAFile
var wg sync.WaitGroup

// SHAFile Save A File Info With SHA1.
type SHAFile struct {
	FileInfo *os.FileInfo
	SHA1     string
}

// GenDirTreeSHA1 ...
// @param path, ignoreDir, ignoreFile, maxGoroutineNum
// @return <-chan, error
func GenDirTreeSHA1(path string, ignoreDir, ignoreFile []string, maxGNum int64) (<-chan SHAFile, error) {

	// set of Max Gorotine Number.
	if maxGNum <= 0 {
		return nil, errors.New("Max Gorotine Number must greater than 0")
	}
	maxGoroutineChan = make(chan int, maxGNum)

	// set of ignoreDirList and ignoreFileList.
	ignoreDirList, ignoreFileList = ignoreDir, ignoreFile

	// Send Count File SHA1 between gorotinue.
	writeFileChan = make(chan SHAFile, 100)
	// Get Dir Tree Start.
	go getDirTree(path)

	return writeFileChan, nil
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
	sha, buf := sha1.New(), make([]byte, 1024)
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
	var sendFileInfo SHAFile
	// sendFileInfo[fileString] = &file
	sendFileInfo.FileInfo = &file
	sendFileInfo.SHA1 = fileString
	// Send to writeFileChan.
	writeFileChan <- sendFileInfo

	// This gorotinue Done.
	<-maxGoroutineChan
	wg.Done()
}

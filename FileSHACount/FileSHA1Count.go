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
func GenDirTreeSHA1(path string, ignoreDir, ignoreFile []string) map[string]*os.FileInfo {
	ignoreDirList, ignoreFileList = ignoreDir, ignoreFile

	// Send Count File SHA1 between gorotinue.
	writeFileChan = make(chan map[string]*os.FileInfo, 100)
	// writeMap Save Result from writeFileChan
	writeMap := make(map[string]*os.FileInfo)
	// Get Dir Tree Start.
	getDirTree(path)

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

	return writeMap
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

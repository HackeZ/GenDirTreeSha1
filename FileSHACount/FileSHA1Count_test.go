package FileSHACount

import (
	"os"
	"testing"
)

func TestFileSHA1Count(t *testing.T) {
	// init Params
	ignoreDir := []string{"HideDir*"}
	ignoreFile := []string{"*.js", "*.t?t"}
	// GenDirTreeSHA1
	resultMap := GenDirTreeSHA1("../TestDir", ignoreDir, ignoreFile)

	if len(resultMap) != 7 {
		t.Fatal("Result Length not Corrent!")
	}

	// init testMap
	testMap := make(map[string]*os.FileInfo)
	testMap["19cf1e53712298df8cb75ea4b817b50d9a7671c4"] = nil
	testMap["a653c1dfd892652dfd439eebca90038c198aa5f9"] = nil
	testMap["89c46b621fe09daf8633b9a91952c40f7dbfc641"] = nil
	testMap["d80f835b13b0b66e8dbea4384ff5025517478520"] = nil
	testMap["753636cbe5ded693ce83dbf470a0b563ea4aa0ac"] = nil
	testMap["82d51cac5de01733455973610709e021adae4fae"] = nil
	testMap["9ce996c7481105375a22677284121ecaafe4d682"] = nil

	// Check SHA1 Corrent or Not.
	for k := range resultMap {
		if _, isExist := testMap[k]; !isExist {
			t.Log(k)
			t.Fatal("File SHA1 Count Wrong!")
		}
	}
}

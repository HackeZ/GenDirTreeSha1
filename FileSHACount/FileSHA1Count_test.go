package FileSHACount

import "testing"

func TestFileSHA1Count(t *testing.T) {
	// init Params
	ignoreDir := []string{"HideDir*"}
	ignoreFile := []string{"*.js", "*.t?t"}
	// GenDirTreeSHA1
	resultChan, err := GenDirTreeSHA1("../TestDir", ignoreDir, ignoreFile, int64(10))

	if err != nil {
		t.Fatal(err)
	}

	var resultLen int

	for {
		res, isOpen := <-resultChan
		if !isOpen {
			break
		}

		switch res.SHA1 {
		case "19cf1e53712298df8cb75ea4b817b50d9a7671c4", "a653c1dfd892652dfd439eebca90038c198aa5f9", "89c46b621fe09daf8633b9a91952c40f7dbfc641", "d80f835b13b0b66e8dbea4384ff5025517478520", "753636cbe5ded693ce83dbf470a0b563ea4aa0ac", "82d51cac5de01733455973610709e021adae4fae", "9ce996c7481105375a22677284121ecaafe4d682":
			resultLen++
		default:
			t.Fatal("File SHA1 Count Wrong!")
		}
	}

	if resultLen != 7 {
		t.Fatal("Result Length not Corrent!")
	}
}

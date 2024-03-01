package obs

import (
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	file, err := os.Open("../../../static/uploadfile/log.txt")
	if err != nil {
		t.Error(err)
		return
	}

	fileName := file.Name()
	userId := int64(1)
	newFileName := GetNewFileName(fileName, userId)
	accessUrl, err := Upload(newFileName, file)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(accessUrl)
}

package fs

import (
	"fmt"
	"testing"
)

func Test_CreateFile(t *testing.T) {
	filename := "C:\\Users\\14595\\Desktop\\储存\\重要资料\\新建文本文档.txt"
	//file := GetLastFile(filename)
	//
	//fmt.Println(file)
	////file := GetLastFile(filename)
	//
	////dir := filename[:-len(file)]
	////
	////fmt.Println(dir)
	//
	//fmt.Println(filename[:len(filename)-len(file)])
	//
	//os.MkdirAll("C:\\Users\\14595\\Desktop\\储存\\重要资料\\", os.ModePerm)
	_ = MkdirAllFile(filename)
}

func Test_ab(t *testing.T) {
	filePaths := []string{
		"C:\\Users\\14595\\Desktop\\FSFC\\fsfc_windows\\1.txt",
		"C:\\Users\\14595\\Desktop\\FSFC\\重要资料",
		"C:\\Users\\14595\\Desktop\\毕业论文\\论文正文",
	}
	for _, filePath := range filePaths {
		//filePaths = append(filePaths, fs.AbsToRelaStore(filePath))
		fmt.Println(AbsToRelaStore(filePath))
	}
	//fmt.Println(filePaths)

}

package fs

import (
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

package fs

import (
	"fmt"
	"os"
	"strings"
)

func FileIsExist(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("文件不存在, 将创建新文件：", filename)
			return err
		}
	}

	return err
}

func GetLastFile(path string) string {
	seqList := strings.Split(path, "\\")
	lastDir := seqList[len(seqList)-1]

	return lastDir
}

func MkdirAllFile(filename string) error {
	file := GetLastFile(filename)

	dir := filename[:len(filename)-len(file)]

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("创建目录失败")
	}

	err = FileIsExist(filename)
	if err != nil {
		_, err := os.Create(filename)
		if err != nil {
			fmt.Println("创建文件失败")
		}
		return err
	}

	return err
}

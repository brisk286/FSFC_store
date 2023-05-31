package fs

import (
	"fmt"
	"fsfc_store/redis"
	"os"
	"strings"
	"time"
)

type Filesystem struct {
	LastSyncTime time.Time
}

var PrimFs Filesystem

func init() {
	PrimFs.LastSyncTime = time.Now()
}

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
	seqList := strings.Split(path, "/")
	lastDir := seqList[len(seqList)-1]

	return lastDir
}

func GetLastFile2(path string) string {
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
		_, err = os.Create(filename)
		if err != nil {
			fmt.Println("创建文件失败")
		}
		return err
	}

	return err
}

func AbsToRelaStore(abs string) string {
	return RelaToAbsRemotePaths(AbsToRela(abs))
}

// RelaToAbsRemotePaths  相对路径转换为linux绝对路径
func RelaToAbsRemotePaths(filenames string) string {
	// todo: 写成接口
	remotePath := "/go/project/syncDir"

	filenames = remotePath + "/" + filenames

	return filenames
}

func AbsToRela(absPath string) string {
	var RelaPath string
	var lastDirs []string

	LocalPath, err := redis.Client.SMembers("AllSaveSpace").Result()
	if err != nil {
		fmt.Println(err)
	}

	for _, localPath := range LocalPath {
		lastDir := "\\" + GetLastFile2(localPath)
		lastDirs = append(lastDirs, lastDir)
	}

	for _, lastDir := range lastDirs {
		if strings.Index(absPath, lastDir) != -1 {
			RelaPath = absPath[strings.Index(absPath, lastDir)+1:]
			RelaPath = strings.ReplaceAll(RelaPath, "\\", "/")
			break
		}
	}

	return RelaPath
}

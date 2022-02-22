package api

import (
	"fmt"
	"fsfc_store/fs"
	"fsfc_store/response"
	"fsfc_store/rsync"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"net/http"
)

func GetChangedFilesAndPostDataList(c *gin.Context) {
	var changedFiles []string
	err := c.ShouldBindBodyWith(&changedFiles, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	fmt.Println("接收到修改文件")

	var hashesFiles []rsync.FileBlockHashes
	for _, filename := range changedFiles {
		err := fs.MkdirAllFile(filename)
		if err != nil {
			panic("文件创建发生错误")
		}
		originalFile, err := ioutil.ReadFile(filename)
		if err != nil {
			panic("未找到远程端文件")
		}
		fmt.Println("读取远程文件成功", filename, originalFile)

		hashes := rsync.CalculateBlockHashes(originalFile)

		fmt.Println("hashes:", hashes)

		hashesFiles = append(hashesFiles, rsync.FileBlockHashes{Filename: filename, BlockHashes: hashes})
	}

	fmt.Println(hashesFiles)

	c.JSON(http.StatusOK, response.SuccessMsg(hashesFiles))
}

func GetRsyncOpsToRebuild(c *gin.Context) {
	fmt.Println("接收到RsyncOps")

	var rsyncOpsResp response.RsyncOpsResp

	err := c.ShouldBindBodyWith(&rsyncOpsResp, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	fmt.Println(rsyncOpsResp)

	filename := rsyncOpsResp.Filename
	rsyncOps := rsyncOpsResp.RsyncOps
	modifiedLength := rsyncOpsResp.ModifiedLength

	original, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("未找到远程端文件")
	} else {
		fmt.Println("找到远程端文件2")
	}
	result := rsync.ApplyOps(original, rsyncOps, modifiedLength)

	//写入文件
	err = ioutil.WriteFile(filename, result, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("写入文件成功")

	c.JSON(http.StatusOK, response.SuccessCodeMsg())
	return
}

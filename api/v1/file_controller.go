package v1

import (
	"fmt"
	"fsfc_store/fs"
	"fsfc_store/request"
	"fsfc_store/response"
	"fsfc_store/rsync"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"net/http"
	"os"
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
			fmt.Println(filename)
			panic("文件创建发生错误")
		}

		originalFile, err := ioutil.ReadFile(filename)
		if err != nil {
			panic("未找到远程端文件")
		}
		fmt.Println("读取远程文件成功", filename)

		fmt.Println("计算BlockHashes")
		hashes := rsync.CalculateBlockHashes(originalFile)

		hashesFiles = append(hashesFiles, rsync.FileBlockHashes{Filename: filename, BlockHashes: hashes})
	}

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

	filename := rsyncOpsResp.Filename
	rsyncOps := rsyncOpsResp.RsyncOps
	modifiedLength := rsyncOpsResp.ModifiedLength

	original, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("未找到远程端文件")
	} else {
		fmt.Println("找到远程端文件2")
	}

	fmt.Println("文件同步中:", filename)
	result := rsync.ApplyOps(original, rsyncOps, modifiedLength)
	err = ioutil.WriteFile(filename, result, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("同步文件成功")

	c.JSON(http.StatusOK, response.SuccessCodeMsg())
}

// RemotePath todo change  存储端应该默认自己所安装的目录，remotePath则设置在同级目录下
const RemotePath = "/var/test"

func MultiDownload(c *gin.Context) {
	var downloadFilenames []string
	err := c.ShouldBindBodyWith(&downloadFilenames, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	for _, filename := range downloadFilenames {
		file, err := os.Open(RemotePath + "/" + filename)
		if err != nil {
			/*c.JSON(http.StatusOK, gin.H{
			    "success": false,
			    "message": "失败",
			    "error":   "资源不存在",
			})*/
			c.Redirect(http.StatusFound, "/404")
			return
		}
		//结束后关闭文件
		defer file.Close()

		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(RemotePath + "/" + filename)
		return
	}
}

func GetFilesInfo(c *gin.Context) {
	filesInfoReq := new(request.FilesInfoReq)
	err := c.ShouldBindBodyWith(&filesInfoReq, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	var filesInfoResp response.FilesInfoResp

	files, _ := ioutil.ReadDir(filesInfoReq.DirPath)
	for _, file := range files {
		if file.IsDir() {
			filesInfoResp.Dirs = append(filesInfoResp.Dirs, file.Name())
		} else {
			filesInfoResp.FilesInfo = append(filesInfoResp.FilesInfo, response.RsyncFileInfo{
				Name:      file.Name(),
				Size:      file.Size(),
				RsyncTime: file.ModTime(),
			})
		}
	}

	c.JSON(http.StatusOK, response.SuccessMsg(filesInfoResp))
}

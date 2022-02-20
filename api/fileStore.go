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
	var hashesFiles []rsync.FileBlockHashes

	changedFiles := c.PostFormArray("changedFiles")

	fmt.Println(changedFiles)

	for _, filename := range changedFiles {
		_ = fs.MkdirAllFile(filename)
		originalFile, _ := ioutil.ReadFile(filename)

		hashes := rsync.CalculateBlockHashes(originalFile)
		hashesFiles = append(hashesFiles, rsync.FileBlockHashes{Filename: filename, BlockHashes: hashes})
	}

	c.JSON(http.StatusOK, response.SuccessMsg(hashesFiles))
}

func GetRsyncOpsToRebuild(c *gin.Context) {
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

	original, _ := ioutil.ReadFile(filename)

	result := rsync.ApplyOps(original, rsyncOps, modifiedLength)

	//写入文件
	err = ioutil.WriteFile(filename, result, 0644)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, response.SuccessCodeMsg())
	return
}

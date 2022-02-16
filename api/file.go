package api

import (
	"fsfc_store/response"
	"fsfc_store/rsync"
	"github.com/gin-gonic/gin"
	"io/ioutil"

	"net/http"
)

func GetChangedFilesAndPostDataList(c *gin.Context) {
	var hashesFiles []rsync.FileBlockHashes

	changedFiles := c.PostFormArray("changedFiles")

	for _, filename := range changedFiles {
		originalFile, _ := ioutil.ReadFile(filename)
		hashes := rsync.CalculateBlockHashes(originalFile)
		hashesFiles = append(hashesFiles, rsync.FileBlockHashes{Filename: filename, BlockHashes: hashes})
	}

	//hashesFiles := []string{"123", "123", "345"}

	c.JSON(http.StatusOK, response.SuccessMsg(hashesFiles))
}

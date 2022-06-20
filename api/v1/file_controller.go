package v1

import (
	"archive/zip"
	"fsfc_store/request"
	"fsfc_store/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RemotePath todo change  存储端应该默认自己所安装的目录，remotePath则设置在同级目录下
const RemotePath = "C:\\Users\\14595\\Desktop\\"

func MultiDownload(c *gin.Context) {
	downloadFilePaths := new(request.DownloadFilePath)
	err := c.ShouldBindBodyWith(&downloadFilePaths, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	err = os.MkdirAll(".\\RsyncFiles", 0777)
	if err != nil {
		return
	}

	for _, filename := range downloadFilePaths.FilePaths {
		//file, err := os.Open(RemotePath + "/" + filename)
		file, err := os.Open("C:\\Users\\14595\\Desktop\\储存\\重要资料\\" + filename)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success":  false,
				"protocol": "失败",
				"error":    "资源不存在",
			})
			c.Redirect(http.StatusFound, "/404")
			return
		}
		//结束后关闭文件
		defer file.Close()

		src, _ := os.Create(".\\RsyncFiles\\" + filename)
		_, err = io.Copy(src, file)
		if err != nil {
			return
		}
	}
	defer os.RemoveAll(".\\RsyncFiles")

	Zip(".\\RsyncFiles", ".\\RsyncFiles.zip")
	defer os.RemoveAll(".\\RsyncFiles.zip")

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"RsyncFiles.zip")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("response-type", "blob") // 以流的形式下载必须设置这一项，否则前端下载下来的文件会出现格式不正确或已损坏的问题
	//c.File(RemotePath + "/" + filename)
	//fmt.Println(RemotePath + filename)
	c.File("RsyncFiles.zip")
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
				Path:      filesInfoReq.DirPath + "\\" + file.Name(),
			})
		}
	}

	c.JSON(http.StatusOK, response.SuccessMsg(filesInfoResp))
}

// 打包成zip文件
func Zip(src_dir string, zip_file_name string) {

	// 预防：旧文件无法覆盖
	os.RemoveAll(zip_file_name)

	// 创建：zip文件
	zipfile, _ := os.Create(zip_file_name)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	filepath.Walk(src_dir, func(path string, info os.FileInfo, _ error) error {

		// 如果是源路径，提前进行下一个遍历
		if path == src_dir {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src_dir+`/`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
}

package v1

import (
	"archive/zip"
	"bufio"
	"fsfc_store/fs"
	"fsfc_store/redis"
	"fsfc_store/request"
	"fsfc_store/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func MultiDownload(c *gin.Context) {
	downloadFilePaths := new(request.DownloadFilePath)
	err := c.ShouldBindBodyWith(&downloadFilePaths, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	err = os.MkdirAll("./RsyncFiles", 0777)
	if err != nil {
		return
	}

	for _, fpath := range downloadFilePaths.FilePaths {
		fInfo, err := os.Stat(fpath)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success":  false,
				"protocol": "失败",
				"error":    "资源不存在",
			})
			c.Redirect(http.StatusFound, "/404")
			return
		}

		if fInfo.IsDir() {
			lastFile := fs.GetLastFile(fpath)
			err := Copy(fpath, "./RsyncFiles/"+lastFile)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"success":  false,
					"protocol": "失败",
					"error":    "文件夹创建失败",
				})
			}
		} else {
			file, _ := os.Open(fpath)

			//结束后关闭文件
			defer file.Close()

			lastFile := fs.GetLastFile(fpath)
			src, _ := os.Create("./RsyncFiles/" + lastFile)
			_, err = io.Copy(src, file)
			if err != nil {
				return
			}
		}

	}
	defer os.RemoveAll("./RsyncFiles")

	Zip("./RsyncFiles", "./RsyncFiles.zip")
	defer os.RemoveAll("./RsyncFiles.zip")

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"RsyncFiles.zip")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("response-type", "blob") // 以流的形式下载必须设置这一项，否则前端下载下来的文件会出现格式不正确或已损坏的问题
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
		uid := uuid.New()

		if file.IsDir() {
			filesInfoResp.Dirs = append(filesInfoResp.Dirs, response.RsyncDirsInfo{
				Id:   uid.String(),
				Name: file.Name(),
				Path: filesInfoReq.DirPath + "/" + file.Name(),
				Icon: "el-icon-arrow-down",
			})
		} else {
			filesInfoResp.Files = append(filesInfoResp.Files, response.RsyncFileInfo{
				Id:        uid.String(),
				Name:      file.Name(),
				Size:      float64(file.Size()) / (1024),
				RsyncTime: file.ModTime().Format("2006-01-02 15:04:05"),
				Path:      filesInfoReq.DirPath + "/" + file.Name(),
			})
		}
	}

	c.JSON(http.StatusOK, response.SuccessMsg(filesInfoResp))
}

//func GetBack(c *gin.Context) {
//	filesInfoReq := new(request.FilesInfoReq)
//	err := c.ShouldBindBodyWith(&filesInfoReq, binding.JSON)
//	if err != nil {
//		c.AbortWithStatusJSON(
//			http.StatusInternalServerError,
//			gin.H{"error": err.Error()})
//		return
//	}
//
//	seqList := strings.Split(filesInfoReq.DirPath, "/")
//	filesInfoReq.DirPath = strings.ReplaceAll(filesInfoReq.DirPath, "/"+seqList[len(seqList)-1], "")
//	fmt.Println(filesInfoReq.DirPath)
//
//	var filesInfoResp response.FilesInfoResp
//
//	files, _ := ioutil.ReadDir(filesInfoReq.DirPath)
//	for _, file := range files {
//		if file.IsDir() {
//			filesInfoResp.Dirs = append(filesInfoResp.Dirs, response.RsyncDirsInfo{
//				Name: file.Name(),
//				Path: filesInfoReq.DirPath + "/" + file.Name(),
//			})
//		} else {
//			filesInfoResp.Files = append(filesInfoResp.Files, response.RsyncFileInfo{
//				Name:      file.Name(),
//				Size:      float64(file.Size()) / (1024),
//				RsyncTime: file.ModTime().Format("2006-01-02 15:04:05"),
//				Path:      filesInfoReq.DirPath + "/" + file.Name(),
//			})
//		}
//	}
//
//	c.JSON(http.StatusOK, response.SuccessMsg(filesInfoResp))
//}

func GetAllSaveSpace(c *gin.Context) {
	var allSaveSpaceResp response.AllSaveSpaceResp
	allSaveSpace, err := redis.Rdb.LRange("AllSaveSpace", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusOK, response.FailCodeMsg())
	}
	allSaveSpaceResp.Dirs = allSaveSpace

	c.JSON(http.StatusOK, response.SuccessMsg(allSaveSpaceResp))
}

func AddSaveSpace(c *gin.Context) {
	addSaveSpaceReq := new(request.AddSaveSpaceReq)
	err := c.ShouldBindBodyWith(&addSaveSpaceReq, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	err = redis.Rdb.RPush("AllSaveSpace", addSaveSpaceReq.DirPath).Err()
	if err != nil {
		c.JSON(http.StatusOK, response.FailCodeMsg())
	}
	c.JSON(http.StatusOK, response.SuccessCodeMsg())
}

func DeleteSaveSpace(c *gin.Context) {
	deleteSaveSpaceReq := new(request.DeleteSaveSpaceReq)
	err := c.ShouldBindBodyWith(&deleteSaveSpaceReq, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	err = redis.Rdb.LRem("AllSaveSpace", 0, deleteSaveSpaceReq.DirPath).Err()
	if err != nil {
		c.JSON(http.StatusOK, response.FailCodeMsg())
	}
	c.JSON(http.StatusOK, response.SuccessCodeMsg())
}

//打包成zip文件
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

func Copy(from, to string) error {
	f, e := os.Stat(from)
	if e != nil {
		return e
	}
	if f.IsDir() {
		//from是文件夹，那么定义to也是文件夹
		if list, e := ioutil.ReadDir(from); e == nil {
			for _, item := range list {
				if e = Copy(filepath.Join(from, item.Name()), filepath.Join(to, item.Name())); e != nil {
					return e
				}
			}
		}
	} else {
		//from是文件，那么创建to的文件夹
		p := filepath.Dir(to)
		if _, e = os.Stat(p); e != nil {
			if e = os.MkdirAll(p, 0777); e != nil {
				return e
			}
		}
		//读取源文件
		file, e := os.Open(from)
		if e != nil {
			return e
		}
		defer file.Close()
		bufReader := bufio.NewReader(file)
		// 创建一个文件用于保存
		out, e := os.Create(to)
		if e != nil {
			return e
		}
		defer out.Close()
		// 然后将文件流和文件流对接起来
		_, e = io.Copy(out, bufReader)
	}
	return e
}

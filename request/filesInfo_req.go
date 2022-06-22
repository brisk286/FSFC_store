package request

type FilesInfoReq struct {
	DirPath string `json:"dirPath"`
}

type DownloadFilePath struct {
	FilePaths []string `json:"filePaths"`
}

type AddSaveSpaceReq struct {
	DirPath string `json:"dirPath"`
}

type DeleteSaveSpaceReq struct {
	DirPath string `json:"dirPath"`
}

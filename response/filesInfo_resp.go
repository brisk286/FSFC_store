package response

type FilesInfoResp struct {
	Dirs  []RsyncDirsInfo `json:"dirs"`
	Files []RsyncFileInfo `json:"files"`
}

type RsyncDirsInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Icon string `json:"icon"`
}

type RsyncFileInfo struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Size      float64 `json:"size"`
	RsyncTime string  `json:"rsyncTime"`
	Path      string  `json:"path"`
}

type AllSaveSpaceResp struct {
	Dirs []string `json:"dirs"`
}

package response

import "time"

type RsyncFileInfo struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	RsyncTime time.Time `json:"rsyncTime"`
}

type FilesInfoResp struct {
	Dirs      []string        `json:"dirs"`
	FilesInfo []RsyncFileInfo `json:"filesInfo"`
}

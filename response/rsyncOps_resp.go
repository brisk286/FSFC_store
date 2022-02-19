package response

import "fsfc_store/rsync"

type RsyncOpsResp struct {
	Filename       string          `json:"filename"`
	RsyncOps       []rsync.RSyncOp `json:"rsyncOps"`
	ModifiedLength int             `json:"ModifiedLength"`
}

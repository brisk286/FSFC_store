package rsync

import "crypto/md5"

type FileBlockHashes struct {
	Filename    string      `json:"filename"`
	BlockHashes []BlockHash `json:"blockHashes"`
}

// BlockHash hash块结构
type BlockHash struct {
	//哈希块下标
	Index int `json:"index"`
	//强哈希值
	StrongHash []byte `json:"strongHash"`
	//弱哈希值
	WeakHash uint32 `json:"weakHash"`
}

// RSyncOp An rsync operation (typically to be sent across the network). It can be either a block of raw data or a block index.
//rsync数据体
type RSyncOp struct {
	//操作类型
	OpCode int `json:"opCode"`
	//如果是DATA 那么保存数据
	Data []byte `json:"data"`
	//如果是BLOCK 保存块下标
	BlockIndex int `json:"blockIndex"`
}

//常量
const (
	// BLOCK 整块数据
	BLOCK = iota
	// DATA 单独修改数据
	DATA
)

const (
	// BlockSize 默认块大小
	//BlockSize = 1024 * 644
	BlockSize = 2
	// M 65536 弱哈希算法取模
	M = 1 << 16
)

// Returns the number of blocks for a given slice of content.
//计算文件需要块的数量
func getBlocksNumber(content []byte) int {
	blockNumber := len(content) / BlockSize
	if len(content)%BlockSize != 0 {
		blockNumber += 1
	}
	return blockNumber
}

// Returns the smaller of a or b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Returns a weak hash for a given block of data.
//弱hash
func weakHash(v []byte) (uint32, uint32, uint32) {
	var a, b uint32
	for i := range v {
		a += uint32(v[i])
		b += (uint32(len(v)-1) - uint32(i) + 1) * uint32(v[i])
	}
	return (a % M) + (1 << 16 * (b % M)), a % M, b % M
}

// Returns a strong hash for a given block of data
func strongHash(v []byte) []byte {
	h := md5.New()
	h.Write(v)
	return h.Sum(nil)
}

// CalculateBlockHashes Returns weak and strong hashes for a given slice.
//计算每个块的哈希值
//参数：全部数据内容
//返回：每个块组成的列表
func CalculateBlockHashes(content []byte) []BlockHash {
	blockHashes := make([]BlockHash, getBlocksNumber(content))
	for i := range blockHashes {
		initialByte := i * BlockSize
		endingByte := min((i+1)*BlockSize, len(content))
		// 确认每个块的定位
		block := content[initialByte:endingByte]
		//计算此块的弱hash
		weak, _, _ := weakHash(block)
		//保存到块哈希数组中
		blockHashes[i] = BlockHash{
			Index:      i,
			StrongHash: strongHash(block),
			WeakHash:   weak,
		}
	}
	return blockHashes
}

// ApplyOps Applies operations from the channel to the original content.
// Returns the modified content.
//根据通道接收到的信息，将数据组装发送
//参数：文件内容，数据操作体 通道， 本地文件大小
//返回:组装后的数据
func ApplyOps(content []byte, rSyncOps []RSyncOp, fileSize int) []byte {
	result := make([]byte, fileSize)

	var offset int

	for _, op := range rSyncOps {
		switch op.OpCode {
		case BLOCK:
			//copy：目标文件，源文件
			copy(result[offset:offset+BlockSize], content[op.BlockIndex*BlockSize:op.BlockIndex*BlockSize+BlockSize])
			offset += BlockSize
		//DATA是不定长的
		case DATA:
			copy(result[offset:], op.Data)
			offset += len(op.Data)
		}
	}

	return result
}

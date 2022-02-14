package rsync

import "crypto/md5"

type FileBlockHashes struct {
	Filename    string
	BlockHashes []BlockHash
}

// BlockHash hash块结构
type BlockHash struct {
	//哈希块下标
	index int
	//强哈希值
	strongHash []byte
	//弱哈希值
	weakHash uint32
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
			index:      i,
			strongHash: strongHash(block),
			weakHash:   weak,
		}
	}
	return blockHashes
}

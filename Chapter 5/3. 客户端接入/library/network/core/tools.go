package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

const (
	idSize  = aes.BlockSize / 2 // 64 bits
	keySize = aes.BlockSize     // 128 bits
)

var (
	ctr []byte
	n   int
	b   []byte
	c   cipher.Block
	m   sync.Mutex
)

// ID 标准数字唯一键
type ID uint64

func init() {
	buf := make([]byte, keySize+aes.BlockSize)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(err) // /dev/urandom had better work
	}
	c, err = aes.NewCipher(buf[:keySize])
	if err != nil {
		panic(err) // AES had better work
	}
	n = aes.BlockSize
	ctr = buf[keySize:]
	b = make([]byte, aes.BlockSize)
}

// GenerateID 返回一个随机生成的 64 字节 ID。这个函数进程安全
// 这个方法大概花费 13.29 ns, 不产出任何内存垃圾
func GenerateID() ID {
	m.Lock()
	if n == aes.BlockSize {
		c.Encrypt(b, ctr)
		for i := aes.BlockSize - 1; i >= 0; i-- { // increment ctr
			ctr[i]++
			if ctr[i] != 0 {
				break
			}
		}
		n = 0
	}
	id := *(*ID)(unsafe.Pointer(&b[n])) // zero-copy b/c we're arch-neutral
	n += idSize
	m.Unlock()
	return id
}

// VerifyAddress 验证地址是否正确
func VerifyAddress(address string) bool {
	var pair []string

	// 是否是 ipv6
	arr := strings.Split(address, "]:")
	if len(arr) == 2 {
		pair = make([]string, 2)
		pair[0] = arr[0][1:]
		pair[1] = arr[1]
	} else {
		// 验证格式
		pair = strings.Split(address, ":")
		if len(pair) != 2 {
			return false
		}
	}
	// 验证 ip
	if pair[0] != "" && net.ParseIP(pair[0]) == nil {
		return false
	}
	// 验证端口
	num, err := strconv.Atoi(pair[1])
	if err != nil {
		return false
	}
	if num < 1 || num > 65535 {
		return false
	}

	return true
}

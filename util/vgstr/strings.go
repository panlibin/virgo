package vgstr

import (
	"strconv"
	"strings"
)

// Hash 字符串哈希
func Hash(str string) uint32 {
	arrByte := []byte(str)
	var nHash uint32 = 1315423911
	for _, c := range arrByte {
		nHash ^= (nHash << 5) + uint32(c) + (nHash >> 2)
	}
	return nHash
}

// IsAskiiLowerLetter 判断是否小写字母
func IsAskiiLowerLetter(r rune) bool {
	return r >= 0x61 && r <= 0x7A
}

// IsAskiiUpperLetter 判断是否大写字母
func IsAskiiUpperLetter(r rune) bool {
	return r >= 0x41 && r <= 0x5A
}

// IsAskiiNumber 判断是否数字
func IsAskiiNumber(r rune) bool {
	return r >= 0x30 && r <= 0x39
}

// IsAskiiLetter 判断是否字母
func IsAskiiLetter(r rune) bool {
	return IsAskiiLowerLetter(r) || IsAskiiUpperLetter(r)
}

// IsAlphanumericOrUnderscore 判断是否字母数字下划线
func IsAlphanumericOrUnderscore(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !IsAskiiLetter(r) && !IsAskiiNumber(r) && r != '_'
	}) < 0
}

// EncodeHexToUpperString 16进制编码[]byte
func EncodeHexToUpperString(src []byte) string {
	const hextable = "0123456789ABCDEF"
	dst := make([]byte, len(src)<<1)
	for i, v := range src {
		dst[i*2] = hextable[v>>4]
		dst[i*2+1] = hextable[v&0x0f]
	}
	return string(dst)
}

// SplitToInt32Array 将字符串分割成[]int32
func SplitToInt32Array(s string, sep string) ([]int32, error) {
	arrStr := strings.Split(s, sep)
	ret := make([]int32, 0, len(arrStr))

	for _, v := range arrStr {
		if v != "" {
			item, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			ret = append(ret, int32(item))
		}
	}
	return ret, nil
}

// SplitToUint32Array 将字符串分割成[]uint32
func SplitToUint32Array(s string, sep string) ([]uint32, error) {
	arrStr := strings.Split(s, sep)
	ret := make([]uint32, 0, len(arrStr))

	for _, v := range arrStr {
		if v != "" {
			item, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return nil, err
			}
			ret = append(ret, uint32(item))
		}
	}
	return ret, nil
}

// SplitToInt64Array 将字符串分割成[]int64
func SplitToInt64Array(s string, sep string) ([]int64, error) {
	arrStr := strings.Split(s, sep)
	ret := make([]int64, 0, len(arrStr))

	for _, v := range arrStr {
		if v != "" {
			item, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			ret = append(ret, item)
		}
	}
	return ret, nil
}

package vgcrypt

import (
	"crypto/md5"
	"fmt"
	"strings"
	"vg_proj/virgo/util/vg_str"
)

// CalcSign 计算签名
func CalcSign(params []string, key string) string {
	strSrc := fmt.Sprintf("%s%s", strings.Join(params, ""), key)
	sum := md5.Sum([]byte(strSrc))
	return vg_str.EncodeHexToUpperString(append(sum[:]))
}

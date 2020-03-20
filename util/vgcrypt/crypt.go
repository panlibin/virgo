package vgcrypt

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/panlibin/virgo/util/vgstr"
)

// CalcSign 计算签名
func CalcSign(params []string, key string) string {
	strSrc := fmt.Sprintf("%s%s", strings.Join(params, ""), key)
	sum := md5.Sum([]byte(strSrc))
	return vgstr.EncodeHexToUpperString(append(sum[:]))
}

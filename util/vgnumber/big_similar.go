package vgnumber

import (
	"math"
	"strconv"
	"strings"
)

// BigIntSimilar 大数字模拟
type BigIntSimilar struct {
	value float64
	tail  int32
}

// SetNumber 设置数值
func (bi *BigIntSimilar) SetNumber(val float64, tail int32) {
	bi.value = val
	bi.tail = tail
	bi.normalize()
}

// SetNumberByString 设置数值
func (bi *BigIntSimilar) SetNumberByString(num string) {
	bi.value = 0
	bi.tail = 0
	sLen := len(num)
	if sLen == 0 {
		return
	}

	var i int
	if num[0] == '-' || num[0] == '+' {
		i = 1
	}

	var c uint8
	for ; i < sLen; i++ {
		c = num[i]
		if (c < 0x30 || c > 0x39) && c != 0x2E {
			break
		}
	}
	bi.value, _ = strconv.ParseFloat(num[:i], 64)

	if i < sLen {
		tailStartIdx := i
		tailLen := sLen - i
		var tailBase string
		tailVal := num[i]
		for ; i < sLen; i++ {
			if num[i] != tailVal {
				tailBase = num[tailStartIdx:i]
				tailVal = num[i]
				break
			}
		}
		if len(tailBase) == 0 {
			if tailLen == 1 {
				tailBase = num[tailStartIdx:i]
				tailVal = 0
			} else if tailLen == 2 {
				tailBase = "aa"
			} else if tailLen == 3 {
				tailBase = "aaa"
			}
		}
		if tailVal >= 0x61 && tailVal <= 0x7A {
			tailVal -= 0x61
		} else if tailVal >= 0x41 && tailVal <= 0x5A {
			tailVal -= 0x41 - 26
		}
		bi.tail = _NumberBase2Int[tailBase] + int32(tailVal)
	}

	bi.normalize()
}

// ToString 转字符串
func (bi *BigIntSimilar) ToString() string {
	strVal := strconv.FormatFloat(bi.value, 'f', 6, 64)
	if bi.tail == 0 {
		return strVal
	}

	tmpTail := bi.tail
	var strNum string

	pNumBase := getNumberBaseByInt(tmpTail)
	if pNumBase != nil {
		if pNumBase.repCount == 0 {
			strNum = strVal + pNumBase.prefix
		} else {
			tmpTail -= pNumBase.baseVal
			if tmpTail >= 0 && tmpTail <= 25 {
				tmpTail += 0x61
			} else if tmpTail >= 26 && tmpTail <= 51 {
				tmpTail += 0x41 - 26
			}
			strNum = strVal + pNumBase.prefix + strings.Repeat(string(tmpTail), pNumBase.repCount)
		}
	} else {
		strNum = strVal + "e" + strconv.FormatInt(int64(tmpTail*3), 10)
	}

	return strNum
}

// Add 加法
func (bi *BigIntSimilar) Add(v *BigIntSimilar) *BigIntSimilar {
	if bi.tail > v.tail {
		v = v.Clone()
		v.setTail(bi.tail)
	} else {
		bi.setTail(v.tail)
	}

	bi.value += v.value
	bi.normalize()
	return bi
}

// Sub 减法
func (bi *BigIntSimilar) Sub(v *BigIntSimilar) *BigIntSimilar {
	if bi.tail > v.tail {
		v = v.Clone()
		v.setTail(bi.tail)
	} else {
		bi.setTail(v.tail)
	}

	bi.value -= v.value
	bi.normalize()
	return bi
}

// Mul 乘法
func (bi *BigIntSimilar) Mul(v *BigIntSimilar) *BigIntSimilar {
	bi.value *= v.value
	bi.tail += v.tail
	bi.normalize()
	return bi
}

// MulFloat 乘法
func (bi *BigIntSimilar) MulFloat(v float64) *BigIntSimilar {
	bi.value *= v
	bi.normalize()
	return bi
}

// Div 除法
func (bi *BigIntSimilar) Div(v *BigIntSimilar) *BigIntSimilar {
	bi.value /= v.value
	bi.tail -= v.tail
	bi.normalize()
	return bi
}

// Pow 乘方
func (bi *BigIntSimilar) Pow(y float64) *BigIntSimilar {
	for y > 10 {
		y /= 10
		bi.value = math.Pow(bi.value, 10)
		bi.tail *= 10
		bi.normalize()
	}

	bi.value = math.Pow(bi.value, y)
	fTail := float64(bi.tail)
	fTail *= y
	bi.tail = int32(fTail)
	fTail -= float64(bi.tail)
	bi.value *= math.Pow(1000, fTail)
	bi.normalize()
	return bi
}

// Compare 比较 0.相等 1.左大 -1.右大
func (bi *BigIntSimilar) Compare(v *BigIntSimilar) int32 {
	if bi.tail > v.tail {
		return 1
	} else if bi.tail < v.tail {
		return -1
	}

	if math.Abs(bi.value-v.value) < 0.000000001 {
		return 0
	} else if bi.value > v.value {
		return 1
	} else if bi.value < v.value {
		return -1
	}
	return 0
}

// Clone 复制
func (bi *BigIntSimilar) Clone() *BigIntSimilar {
	return &BigIntSimilar{
		value: bi.value,
		tail:  bi.tail,
	}
}

// IsZero 是否为0
func (bi *BigIntSimilar) IsZero() bool {
	return bi.tail <= 0 && bi.value < 0.000000001
}

// IsNegative 是否负数
func (bi *BigIntSimilar) IsNegative() bool {
	return !bi.IsZero() && bi.value < 0
}

// Minus 取负数
func (bi *BigIntSimilar) Minus() {
	bi.value *= -1
}

// GetValue 获取数值
func (bi *BigIntSimilar) GetValue() float64 {
	return bi.value
}

// GetTail 获取幂数
func (bi *BigIntSimilar) GetTail() int32 {
	return bi.tail
}

// normalize 标准化数值1000以内
func (bi *BigIntSimilar) normalize() {
	if bi.value > 1000 || bi.value < -1000 {
		for bi.value > 1000 || bi.value < -1000 {
			bi.value /= 1000
			bi.tail++
		}
	} else if bi.value < 1 && bi.value > -1 {
		for bi.value < 1 && bi.value > -1 && bi.tail > 0 {
			bi.value *= 1000
			bi.tail--
		}
	}
}

func (bi *BigIntSimilar) setTail(tail int32) {
	if bi.tail == tail {
		return
	}

	offset := bi.tail - tail
	if offset < -3 {
		bi.value = 0
		bi.tail = tail
		return
	}

	bi.value *= math.Pow10(int(offset * 3))
	bi.tail = tail
}

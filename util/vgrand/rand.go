package vgrand

import "math/rand"

// RandInt31 获取范围内随机数
func RandInt31(min int32, max int32) int32 {
	return rand.Int31n(max-min) + min
}

// RandInt63 获取范围内随机数
func RandInt63(min int64, max int64) int64 {
	return rand.Int63n(max-min) + min
}

// RandUint32 获取范围内随机数
func RandUint32(min uint32, max uint32) uint32 {
	return rand.Uint32()%(max-min) + min
}

// RandUint64 获取范围内随机数
func RandUint64(min uint64, max uint64) uint64 {
	return rand.Uint64()%(max-min) + min
}

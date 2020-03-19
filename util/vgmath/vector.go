package vgmath

import (
	"fmt"
	"math"
)

// AxisX x轴
var AxisX *Vector3 = NewVector3(1, 0, 0)

// AxisY y轴
var AxisY *Vector3 = NewVector3(0, 1, 0)

// AxisZ z轴
var AxisZ *Vector3 = NewVector3(0, 0, 1)

// Vector3 向量
type Vector3 struct {
	x, y, z float64
}

// NewVector3 新建向量
func NewVector3(x, y, z float64) *Vector3 {
	pVec := new(Vector3)
	pVec.x, pVec.y, pVec.z = x, y, z
	return pVec
}

// Add 向量相加
func (v *Vector3) Add(vr *Vector3) *Vector3 {
	v.x += vr.x
	v.y += vr.y
	v.z += vr.z
	return v
}

// Sub 向量相减
func (v *Vector3) Sub(vr *Vector3) *Vector3 {
	v.x -= vr.x
	v.y -= vr.y
	v.z -= vr.z
	return v
}

// Dot 向量内积
func (v *Vector3) Dot(vr *Vector3) float64 {
	return v.x*vr.x + v.y*vr.y + v.z*vr.z
}

// Cross 向量外积
func (v *Vector3) Cross(vr *Vector3) *Vector3 {
	x := v.y*vr.z - vr.y*v.z
	y := vr.x*v.z - v.x*vr.z
	z := v.x*vr.y - vr.x*v.y
	v.x, v.y, v.z = x, y, z
	return v
}

// MulMatrix 向量乘以矩阵
func (v *Vector3) MulMatrix(mat *Matrix33) *Vector3 {
	x := v.x*mat.a[0] + v.y*mat.a[1] + v.z*mat.a[2]
	y := v.x*mat.a[3] + v.y*mat.a[4] + v.z*mat.a[5]
	z := v.x*mat.a[6] + v.y*mat.a[7] + v.z*mat.a[8]
	v.x, v.y, v.z = x, y, z
	return v
}

// MulNumber 向量乘以数字
func (v *Vector3) MulNumber(n float64) *Vector3 {
	v.x *= n
	v.y *= n
	v.z *= n
	return v
}

// DivNumber 向量除以数字
func (v *Vector3) DivNumber(n float64) *Vector3 {
	rn := 1 / n
	return v.MulNumber(rn)
}

// Negative 取反
func (v *Vector3) Negative() *Vector3 {
	v.x, v.y, v.z = -v.x, -v.y, -v.z
	return v
}

// ProjectLength 指定方向上的投影长度
func (v *Vector3) ProjectLength(dir *Vector3) float64 {
	nd := dir.Clone()
	nd.Normalize()
	return v.Dot(nd)
}

// Project 投影向量
func (v *Vector3) Project(dir *Vector3) *Vector3 {
	nd := dir.Clone()
	nd.Normalize()
	pLen := v.Dot(nd)
	nd.MulNumber(pLen)
	v.x, v.y, v.z = nd.x, nd.y, nd.z
	return v
}

// Normalize 标准化
func (v *Vector3) Normalize() *Vector3 {
	length := v.Length()
	if math.Abs(length-1) > 0.000000001 {
		v.DivNumber(length)
	}
	return v
}

// Clone 复制
func (v *Vector3) Clone() *Vector3 {
	return NewVector3(v.x, v.y, v.z)
}

// LengthSquare 长度平方
func (v *Vector3) LengthSquare() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

// Length 长度
func (v *Vector3) Length() float64 {
	return math.Sqrt(v.LengthSquare())
}

// ToString 转字符串
func (v *Vector3) ToString() string {
	return fmt.Sprintf("[%f, %f, %f]", v.x, v.y, v.z)
}

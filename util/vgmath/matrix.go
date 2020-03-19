package vgmath

import "math"

// Matrix33 3*3矩阵
type Matrix33 struct {
	a [9]float64
}

// MakeIdentityMatrix33 创建单位矩阵
func MakeIdentityMatrix33() *Matrix33 {
	pMat := new(Matrix33)
	pMat.a[0], pMat.a[4], pMat.a[8] = 1, 1, 1
	return pMat
}

// MakeRotateMatrix33 创建旋转矩阵
func MakeRotateMatrix33(axis *Vector3, angle float64) *Matrix33 {
	pMat := new(Matrix33)
	pMat.SetRotate(axis, angle)
	return pMat
}

// SetRotate 旋转矩阵
func (m *Matrix33) SetRotate(axis *Vector3, angle float64) *Matrix33 {
	angle = math.Pi * angle / 180
	s := math.Sin(angle)
	c := math.Cos(angle)

	if axis.x != 0 && axis.y == 0 && axis.z == 0 {
		if axis.x < 0 {
			s = -s
		}
		m.a[0], m.a[3], m.a[6] = 1, 0, 0
		m.a[1], m.a[4], m.a[7] = 0, c, s
		m.a[2], m.a[5], m.a[8] = 0, -s, c
	} else if axis.x == 0 && axis.y != 0 && axis.z == 0 {
		if axis.y < 0 {
			s = -s
		}
		m.a[0], m.a[3], m.a[6] = c, 0, -s
		m.a[1], m.a[4], m.a[7] = 0, 1, 0
		m.a[2], m.a[5], m.a[8] = s, 0, c
	} else if axis.x == 0 && axis.y == 0 && axis.z != 0 {
		if axis.z < 0 {
			s = -s
		}
		m.a[0], m.a[3], m.a[6] = c, s, 0
		m.a[1], m.a[4], m.a[7] = -s, c, 0
		m.a[2], m.a[5], m.a[8] = 0, 0, 1
	} else {
		axis = axis.Clone()
		axis.Normalize()
		nc := 1 - c
		xy := axis.x * axis.y
		yz := axis.y * axis.z
		xz := axis.x * axis.z
		xs := axis.x * s
		ys := axis.y * s
		zs := axis.z * s

		m.a[0], m.a[3], m.a[6] = axis.x*axis.x*nc+c, xy*nc+zs, xz*nc-ys
		m.a[1], m.a[4], m.a[7] = xy*nc-zs, axis.y*axis.y*nc+c, yz*nc+xs
		m.a[2], m.a[5], m.a[8] = xz*nc+ys, yz*nc-xs, axis.z*axis.z*nc+c
	}

	return m
}

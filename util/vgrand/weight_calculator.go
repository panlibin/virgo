package vgrand

import (
	"math/rand"
	"sort"
)

// WeightElement 权重元素
type WeightElement struct {
	weightSum int64
	content   interface{}
}

// WeightArray 权重元素数组
type WeightArray []*WeightElement

func (w WeightArray) Len() int {
	return len(w)
}

func (w WeightArray) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w WeightArray) Less(i, j int) bool {
	return w[i].weightSum < w[j].weightSum
}

// WeightCalculator 权重计算器
type WeightCalculator struct {
	totalWeight int64
	arrElement  WeightArray
}

// NewWeightCalculator 创建权重计算器
func NewWeightCalculator() *WeightCalculator {
	pObj := new(WeightCalculator)
	return pObj
}

// AddElement 添加元素
func (w *WeightCalculator) AddElement(weight int64, content interface{}) {
	w.totalWeight += weight
	pElement := new(WeightElement)
	pElement.weightSum = w.totalWeight
	pElement.content = content
	w.arrElement = append(w.arrElement, pElement)
}

// Random 根据权重获取随机元素
func (w *WeightCalculator) Random() interface{} {
	if w.totalWeight <= 0 {
		return nil
	}

	randKey := rand.Int63n(w.totalWeight)
	idx := sort.Search(len(w.arrElement), func(i int) bool {
		return w.arrElement[i].weightSum > randKey
	})

	if idx >= len(w.arrElement) {
		return nil
	}

	return w.arrElement[idx].content
}

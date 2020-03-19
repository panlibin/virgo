package vgnumber

import "sort"

var _NumberBase2Int = map[string]int32{
	"K":   1,
	"M":   2,
	"B":   3,
	"T":   4,
	"aa":  5,
	"aaa": 57,
	"$":   109,
	"$$":  161,
	"@":   213,
	"@@":  265,
	"#":   317,
	"##":  369,
	"%":   421,
	"%%":  473,
	"&":   525,
	"&&":  577,
}

type _NumberBase struct {
	repCount int
	prefix   string
	baseVal  int32
}

var _Int2NumberBase = []_NumberBase{
	{repCount: 1, prefix: "&&", baseVal: 577},
	{repCount: 2, prefix: "&", baseVal: 525},
	{repCount: 1, prefix: "%%", baseVal: 473},
	{repCount: 2, prefix: "%", baseVal: 421},
	{repCount: 1, prefix: "##", baseVal: 369},
	{repCount: 2, prefix: "#", baseVal: 317},
	{repCount: 1, prefix: "@@", baseVal: 265},
	{repCount: 2, prefix: "@", baseVal: 213},
	{repCount: 1, prefix: "$$", baseVal: 161},
	{repCount: 2, prefix: "$", baseVal: 109},
	{repCount: 3, prefix: "", baseVal: 57},
	{repCount: 2, prefix: "", baseVal: 5},
	{repCount: 0, prefix: "T", baseVal: 4},
	{repCount: 0, prefix: "B", baseVal: 3},
	{repCount: 0, prefix: "M", baseVal: 2},
	{repCount: 0, prefix: "K", baseVal: 1},
}

func getNumberBaseByInt(tail int32) *_NumberBase {
	if tail > 628 {
		return nil
	}
	idx := sort.Search(len(_Int2NumberBase), func(i int) bool {
		return _Int2NumberBase[i].baseVal <= tail
	})
	return &_Int2NumberBase[idx]
}

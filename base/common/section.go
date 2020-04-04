package common

import (
	"errors"
	"sync"
)

const PerSectionIdSize uint64 = 100000

var NotInRangeSet = errors.New("id not int the set range")

type RangeID struct {
	IdBegin uint64
	Size    uint64
}

type SectionID uint64

type Section struct {
	RangeID

	Mut    sync.RWMutex
	MaxSeq uint64
	SeqNum []uint64
}

func NewSection(id SectionID, maxSeq uint64) *Section {
	s := &Section{
		RangeID: RangeID{
			IdBegin: uint64(id) * PerSectionIdSize,
			Size:    PerSectionIdSize - 1,
		},
		MaxSeq: maxSeq,
		SeqNum: make([]uint64, PerSectionIdSize),
	}
	Memset(s.SeqNum, s.MaxSeq)
	return s
}

type RouterMap map[string][]SectionID

// 获取section在set里的位置
func CalcIndex(rangeId RangeID, id uint64) (bool, uint64) {
	if !CheckIDByRange(rangeId, id) {
		return false, 0
	}
	return true, (id - rangeId.IdBegin) / PerSectionIdSize
}

// 检查id是否在当前set里
func CheckIDByRange(rangeId RangeID, id uint64) bool {
	return id >= rangeId.IdBegin && id < rangeId.IdBegin+rangeId.Size
}

func GetSectionIDByUid(uid uint64) SectionID {
	id := uid / PerSectionIdSize
	if uid%PerSectionIdSize != 0 {
		id++
	}
	return SectionID(id)
}

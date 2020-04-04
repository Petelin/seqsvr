package common

import "testing"

func TestSection(t *testing.T) {
	// 0 - 42950
	println(GetSectionIDByUid(0))
	println(GetSectionIDByUid(PerSectionIdSize))
	println(GetSectionIDByUid(1<<32 - 1))
}

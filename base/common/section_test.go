package common

import (
	"fmt"
	"testing"
)

func TestSection(t *testing.T) {
	// 0 - 42950
	fmt.Println(1<<32 - 1)
	println(GetSectionIDByUid(0))
	println(GetSectionIDByUid(PerSectionIdSize - 1))
	println(GetSectionIDByUid(PerSectionIdSize))
	println(GetSectionIDByUid(PerSectionIdSize + 1))
	println(GetSectionIDByUid(1<<32 - 1))
}

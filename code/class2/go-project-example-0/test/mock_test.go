package test

import (
	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessFirstLine(t *testing.T) {
	firstLine := ProcessFirstLine()
	assert.Equal(t, "line00", firstLine)
}

func TestProcessFirstLineWithMock(t *testing.T) {
	// 在调用的过程中其实调用的是打桩函数
	monkey.Patch(ReadFirstLine, func() string {
		return "line110"
	})
	// 卸载桩
	defer monkey.Unpatch(ReadFirstLine)
	line := ProcessFirstLine()
	assert.Equal(t, "line000", line)
}

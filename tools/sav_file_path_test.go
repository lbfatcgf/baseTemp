package tools_test

import (
	"fmt"
	"testing"

	"codeup.aliyun.com/67c7c688484ca2f0a13acc04/baseTemp/tools"
)

func TestSaveFilePath(t *testing.T) {
	ps := []string{
		"../test.txt",
		"../s/test.txt",
		"../../test.txt",
		"../ss/../test.txt",
		"../ss/../sss/test.txt",
		"./test.txt",
		"test.txt",
		"C:/ol/../ss/../test.txt",
	}
	for _, p := range ps {
		fmt.Println(tools.SafeFilePath(p))
	}
}

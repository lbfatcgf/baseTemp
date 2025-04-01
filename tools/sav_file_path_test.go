package tools_test

import (
	"baseTemp/tools"
	"fmt"
	"testing"
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

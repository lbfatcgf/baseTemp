package tools

import (
	"path/filepath"
	"runtime"
	"strings"
)

// 过滤文件路径，防止路径穿越
//
// ../test.txt -> test.txt
//
// ../s/test.txt -> /s/test.txt
//
// ../../test.txt -> test.txt
//
// ../ss/../test.txt -> test.txt
//
// ../ss/../sss/test.txt -> /sss/test.txt
//
// ./test.txt -> test.txt
//
// test.txt -> test.txt
//
// C:/ol/../ss/../test.txt -> test.txt
func SafeFilePath(fp string) string {
	if runtime.GOOS == "windows" {
		volumeName := filepath.VolumeName(fp)
		if len(volumeName) > 0 {
			fp = fp[len(volumeName)+1:]
		}
	}
	cfp := filepath.Clean(fp)
	if strings.Contains(cfp, "..") {
		cfp = cfp[strings.LastIndex(cfp, "..")+2:]
	}
	return cfp
}

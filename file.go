package snake

import (
	"bytes"
	"os"
)

type snakefile struct {
	Input *os.File
}

// FileOperate ...
type FileOperate interface {
	Get() *os.File
	Text() String
}

// ---------------------------------------
// 输入 :

// File 初始化...
func File(f *os.File) FileOperate {
	return &snakefile{Input: f}
}

// ---------------------------------------
// 输出 :

// Get 获取文本...
func (sk *snakefile) Get() *os.File {
	return sk.Input
}

// Text 获取文本...
func (sk *snakefile) Text() String {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(sk.Input)
	if err != nil {
		// todo: 字符串转化错误消息
	}
	return Text(buf.String())
}

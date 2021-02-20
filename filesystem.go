package snake

import (
	"os"
	"path/filepath"

	"github.com/plandem/xlsx"
)

// FileSystem ...
type FileSystem interface {
	Add(str ...string) FileSystem             // 新增路径
	Dir() string                              // 返回目录路径
	Base() string                             // 返回路径中最后一个元素
	IsDir(dst ...string) bool                 // 判断是否为目录
	IsFile(dst ...string) bool                // 判断是否为文件
	Ls(opt ...string) []string                // 查看文件夹列表
	Find(opt ...string) []string              // 查找文件
	MkDir(dst ...string) bool                 // 新建文件夹
	MkFile(dst ...string) (FileOperate, bool) // 新建文件
	Open() (FileOperate, bool)                // 打开文件
	OpenXlsx() (Excel, bool)                  // 打开Excel Xlsx文件
	Exist(dst ...string) bool                 // 判断目录或文件是否存在
	Rm(dst ...string) bool                    // 删除目录或文件
	Rn(newname string) bool                   // 修改目录或文件名
	Mv(newpath string) bool                   // 移动目录或文件到指定位置
	Cp(dir string, overwrite bool) bool       // 拷贝目录或文件到指定位置
	// SameFile()                      // 文件对比
	// Chmod()                         // 设置权限
	// Chown()                         // 设置用户、用户组
	Ext() string // 返回文件扩展名
	Get() string // 返回路径
}

type snakeFileSystem struct {
	Path string
}

// ---------------------------------------
// 输入 :

// FS 初始化...
func FS(str ...string) FileSystem {
	sk := &snakeFileSystem{}
	return sk.Add(str...)
}

// Add 在字符串中追加文字...
func (sk *snakeFileSystem) Add(str ...string) FileSystem {
	if len(str) > 0 {
		for _, v := range str {
			sk.Path = filepath.Join(sk.Path, v)
		}
	}
	return sk
}

// ---------------------------------------
// 处理 :

// Cp 拷贝目录或文件
func (sk *snakeFileSystem) Cp(dir string, overwrite bool) bool {
	dst := FS(dir)

	// todo:目标存在则返回错误
	if dst.Add(sk.Base()).Exist() && overwrite == false {
		return false
	}

	// todo:目标与源相同
	if dst.Add(sk.Base()).Get() == sk.Get() {
		return false
	}

	if sk.IsFile() {
		// 覆盖拷贝文件
		// todo: 拷贝文件错误信息
		return _owcpfile(sk, dst.Add(sk.Base()))
	} else if sk.IsDir() {
		// 覆盖拷贝目录
		if dst.Add(sk.Base()).Exist() {
			dst.Rm()
		}
		for _, v := range sk.Find() {
			src := FS(v)
			if src.IsFile() {
				// todo: 拷贝文件错误信息
				_owcpfile(src, FS(dst.Get(), src.Get()))
			} else if sk.IsDir() {
				// todo: 拷贝目录错误信息
				FS(dst.Get(), src.Get()).MkDir()
			}
		}
	}

	return true
}

// Rm 删除目录及文件
func (sk *snakeFileSystem) Rm(dst ...string) bool {
	return os.RemoveAll(sk.pathdst(dst...)) == nil
}

// Open 打开文件
func (sk *snakeFileSystem) Open() (FileOperate, bool) {
	file, err := os.Open(sk.Path)
	return File(file), err == nil
}

// OpenXlsx 打开Xlsx文件
func (sk *snakeFileSystem) OpenXlsx() (Excel, bool) {
	file, err := xlsx.Open(sk.Path)
	return Xlsx(file), err == nil
}

// Rn 修改目录或文件名
func (sk *snakeFileSystem) Rn(newname string) bool {
	res := os.Rename(sk.Path, filepath.Join(sk.Dir(), newname)) == nil
	if res {
		sk.Path = filepath.Join(sk.Dir(), newname)
	}
	return res
}

// Mv 移动目录或文件到指定位置
func (sk *snakeFileSystem) Mv(newpath string) bool {
	res := os.Rename(sk.Path, filepath.Join(newpath, sk.Base())) == nil
	if res {
		sk.Path = filepath.Join(newpath, sk.Base())
	}
	return res
}

// Ext 扩展名
func (sk *snakeFileSystem) Ext() string {
	return filepath.Ext(sk.Path)
}

// MkDir 创建目录
func (sk *snakeFileSystem) MkDir(dst ...string) bool {
	return os.MkdirAll(sk.pathdst(dst...), os.ModePerm) == nil
}

// MkFile 创建文件
func (sk *snakeFileSystem) MkFile(dst ...string) (FileOperate, bool) {
	p := FS(sk.pathdst(dst...))
	if FS(p.Dir()).Exist() == false {
		sk.MkDir(p.Dir())
	}
	file, err := os.Create(p.Get())
	return File(file), err == nil
}

// Exist 判断文件或目录是否存在
func (sk *snakeFileSystem) Exist(dst ...string) bool {
	if _, err := os.Stat(sk.pathdst(dst...)); err != nil {
		return os.IsExist(err)
	}
	return true
}

// Ls 返回路径目录下内容
// dst为空时，则返回当前路径列表；
// 例子：
// snake.FS("./").LS()
// 返回：./路径下的目录与文件
// dst 可设置多个参数，参数可使用正则方式搜索当前路径下列表；
// 例子：
// snake.FS("./").LS("*.go")
// 返回：./路径下的扩展名为.go的所有文件或目录
func (sk *snakeFileSystem) Ls(opt ...string) []string {
	if len(opt) == 0 {
		return ls(sk.Path, "*")
	}
	return ls(sk.Path, opt...)
}

// Find 根据条件搜索路径目录下内容
// 功能与Ls()方法一直，区别在于Find可以对当前路径下所有目录遍历搜索并返回列表。
func (sk *snakeFileSystem) Find(opt ...string) []string {
	if len(opt) == 0 {
		return walkPath(sk.Path, "*")
	}
	return walkPath(sk.Path, opt...)
}

// Dir 获取目录名
func (sk *snakeFileSystem) Dir() string {
	return filepath.Dir(sk.Path)
}

// Base 返回路径中最后一个元素
func (sk *snakeFileSystem) Base() string {
	return filepath.Base(sk.Path)
}

// IsDir 判断是否是目录
func (sk *snakeFileSystem) IsDir(dst ...string) bool {
	if i, err := os.Stat(sk.pathdst(dst...)); os.IsExist(err) == false {
		return i.Mode().IsDir()
	}
	return false
}

// IsFile 判断是否是目录
func (sk *snakeFileSystem) IsFile(dst ...string) bool {
	if i, err := os.Stat(sk.pathdst(dst...)); os.IsExist(err) == false {
		return i.Mode().IsRegular()
	}
	return false
}

// pathdst 处理方法中dst数组，当dst数组为空时，输出Path值，不为空时，输出dst数组的第一个元素。
func (sk *snakeFileSystem) pathdst(dst ...string) string {
	if len(dst) > 0 {
		return dst[0]
	}
	return sk.Path
}

// Get 获取文本...
func (sk *snakeFileSystem) Get() string {
	return sk.Path
}

package snake

import (
	"github.com/plandem/xlsx"
)

type snakexlsx struct {
	Input *xlsx.Spreadsheet
}

// Excel ...
type Excel interface {
	Get() *xlsx.Spreadsheet
	Sheet(col, row int) string
	List(sheet xlsx.Sheet) []map[int]string
}

// ---------------------------------------
// 输入 :

// Xlsx 初始化...
func Xlsx(i *xlsx.Spreadsheet) Excel {
	return &snakexlsx{Input: i}
}

// ---------------------------------------
// 输出 :

// Get 获取文本...
func (sk *snakexlsx) Get() *xlsx.Spreadsheet {
	return sk.Input
}

// Sheet 获取Sheet...
func (sk *snakexlsx) Sheet(col, row int) string {
	xl := sk.Get()
	defer xl.Close()
	// 当前仅输出第一页
	sheet := xl.Sheet(0, xlsx.SheetModeStream)
	defer sheet.Close()
	return sheet.Cell(col, row).String()
}

// List 获取某页的内容列表 ...
func (sk *snakexlsx) List(sheet xlsx.Sheet) []map[int]string {
	var res []map[int]string
	defer sheet.Close()
	totalCols, totalRows := sheet.Dimension()
	totalRows = totalRows + 2
	for row := 2; row < totalRows; row++ {
		r := make(map[int]string)
		for col := 0; col < totalCols; col++ {
			r[col+1] = sheet.Cell(col, row).String()
		}
		res = append(res, r)
	}

	return res
}

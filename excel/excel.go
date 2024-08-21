package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
)

// CreateExcel  xlsx
func CreateExcel[T any](rows []*T) *excelize.File {
	xlsx := excelize.NewFile()
	var typeOf T
	t := reflect.TypeOf(typeOf)
	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 解析自定义标签
		if tag := field.Tag.Get("excel"); tag != "" {
			label := fmt.Sprintf("%s1", IntToColumnLabel(i+1))
			_ = xlsx.SetCellValue("Sheet1", label, tag)
			style := &excelize.Style{
				Fill: excelize.Fill{
					Type:    "pattern",
					Pattern: 1,
					Color:   []string{"#EEECE1"},
				},
				Font: &excelize.Font{
					Bold:   true,
					Family: "宋体",
					Size:   10,
				},
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
			}
			_style, err := xlsx.NewStyle(style)
			if err == nil {
				_ = xlsx.SetCellStyle("Sheet1", label, label, _style)
			}
		}
	}
	for k, val := range rows {
		valueOf := reflect.ValueOf(*val)
		for i := 0; i < valueOf.NumField(); i++ {
			label := fmt.Sprintf("%s%d", IntToColumnLabel(i+1), k+2)
			_ = xlsx.SetCellValue("Sheet1", label, valueOf.Field(i).Interface())

			style := &excelize.Style{
				Fill: excelize.Fill{
					Type:    "pattern",
					Pattern: 1,
				},
				Font: &excelize.Font{
					Bold:   true,
					Family: "宋体",
					Size:   10,
				},
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
			}
			_style, err := xlsx.NewStyle(style)
			if err == nil {
				_ = xlsx.SetCellStyle("Sheet1", label, label, _style)
			}
		}

	}
	return xlsx
}
func IntToColumnLabel(n int) string {
	var columnLabel string
	for n > 0 {
		n-- // Adjust for 1-based indexing
		remainder := n % 26
		columnLabel = string(rune('A'+remainder)) + columnLabel
		n /= 26
	}
	return columnLabel
}

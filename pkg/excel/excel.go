// Package excel 提供 Excel 导入导出工具。
//
// 导出：将 []struct 写入 .xlsx 文件返回字节流
// 导入：从 .xlsx 文件读出 []map[string]any（按表头匹配）
//
// 用法见 example_test.go。
package excel

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

// Export 把任意 slice 导出为 .xlsx 字节流
//   - sheetName：sheet 名（默认 "Sheet1"）
//   - headers：表头（中文名）按 struct 字段顺序排列；nil 时用 struct tag "excel" 或字段名
//   - fieldOrder：struct 字段导出顺序（按 reflect name）
//
// Example:
//
//	type User struct { Name string; Age int }
//	users := []User{{"alice", 30}, {"bob", 25}}
//	data, _ := excel.Export(users, []string{"姓名", "年龄"}, nil)
func Export(slice any, headers []string, fieldOrder []string) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "Sheet1"

	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected slice, got %s", v.Kind())
	}

	// 推导字段顺序
	if fieldOrder == nil {
		if v.Len() > 0 {
			elem := v.Index(0)
			for i := 0; i < elem.NumField(); i++ {
				fieldOrder = append(fieldOrder, elem.Type().Field(i).Name)
			}
		}
	}

	// 写表头
	if headers == nil {
		headers = fieldOrder
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 写数据
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		for j, fname := range fieldOrder {
			fieldVal := elem.FieldByName(fname)
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, formatValue(fieldVal))
		}
	}

	// 输出
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Import 把 .xlsx 字节流读为 []map[string]any
//   - 第一行是表头
//   - 每行是一个 map（key 是表头文本，value 是单元格内容）
//
// Example:
//
//	rows, err := excel.Import(fileBytes)
//	for _, row := range rows {
//	    fmt.Println(row["姓名"], row["年龄"])
//	}
func Import(data []byte) ([]map[string]any, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}

	header := rows[0]
	result := make([]map[string]any, 0, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		row := make(map[string]any, len(header))
		for j, h := range header {
			if j < len(rows[i]) {
				row[h] = rows[i][j]
			}
		}
		result = append(result, row)
	}
	return result, nil
}

// formatValue 把 reflect.Value 格式化成适合 Excel 的值
func formatValue(v reflect.Value) any {
	if !v.IsValid() {
		return ""
	}
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.Bool:
		return v.Bool()
	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok {
			return t.Format("2006-01-02 15:04:05")
		}
		return fmt.Sprintf("%v", v.Interface())
	default:
		if v.CanInt() {
			return strconv.FormatInt(v.Int(), 10)
		}
		return fmt.Sprintf("%v", v.Interface())
	}
}
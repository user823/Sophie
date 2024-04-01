package excelutil

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// 在目标位置处生成一个结构体的记录
const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// 可处理的类型：结构体、结构体指针、结构体切片、结构体指针切片
func WriteXlsx(file *excelize.File, sheet string, records interface{}) {
	index, _ := file.NewSheet(sheet)
	file.SetActiveSheet(index)
	t := reflect.TypeOf(records)

	if t.Kind() != reflect.Slice {
		_, err := WriteStruct(file, sheet, 0, 0, records)
		if err != nil {
			file.DeleteSheet(sheet)
		}
		return
	}

	// 处理切片类型
	s := reflect.ValueOf(records)
	index = 0
	for i := 0; i < s.Len(); i++ {
		elem := s.Index(i).Interface()
		_, err := WriteStruct(file, sheet, index, 0, elem)
		if err == nil {
			index++
		}
	}
}

// col 从0开始
func num2Col(num int) string {
	var result string
	for num >= 0 {
		result = string(alphabet[num%26]) + result
		num = num/26 - 1
	}
	return result
}

func pos(row, col int) string {
	return fmt.Sprintf("%s%d", num2Col(col), row+1)
}

func WriteStruct(file *excelize.File, sheet string, row, col int, record any) (int, error) {
	elemValue := reflect.ValueOf(record)
	// 如果是指针则获取值
	for elemValue.Kind() == reflect.Ptr {
		elemValue = elemValue.Elem()
	}
	elemType := elemValue.Type()
	origin := col
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("xlsx")
		mp := getKV(tag)

		// 忽略该字段
		if tag == "" || tag == "-" {
			continue
		}

		// 如果是指针则获取值
		fieldValue := elemValue.Field(i)
		for fieldValue.Type().Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		// 递归处理结构体类型
		switch fieldValue.Type().Kind() {
		case reflect.Struct:
			if _, ok := mp[INLINE]; !ok {
				continue
			}
			cnt, err := WriteStruct(file, sheet, row, col, elemValue.Field(i).Interface())
			if err != nil {
				return 0, err
			}
			col += cnt
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32,
			reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:

			// 首行设置
			if row == 0 {
				exportName := field.Name
				if name, ok := mp[NAME]; ok {
					exportName = name
				}
				err := file.SetCellValue(sheet, pos(row, col), exportName)
				if err != nil {
					return 0, nil
				}
			}
			err := file.SetCellValue(sheet, pos(row+1, col), elemValue.Field(i).Interface())
			if err != nil {
				return 0, err
			}
			// 处理每个kv对
			setStyle(file, sheet, row+1, col, mp)
			col++
		default:
			return 0, fmt.Errorf("不支持该类型: %s", field.Type.Kind())
		}
	}
	return col - origin, nil
}

// 处理kv对
func getKV(tag string) (res map[string]string) {
	res = make(map[string]string)
	for _, kv := range strings.Split(tag, SEPARATOR) {
		str := strings.Split(kv, ":")
		if len(str) == 2 {
			res[str[0]] = str[1]
		} else {
			res[str[0]] = ""
		}
	}
	return
}

// 设置单元格stype
func setStyle(file *excelize.File, sheetname string, row int, col int, mp map[string]string) {
	for k, v := range mp {
		switch k {
		case WIDTH:
			now, _ := file.GetColWidth(sheetname, num2Col(col))
			toSet, _ := strconv.ParseFloat(v, 64)
			file.SetColWidth(sheetname, num2Col(col), num2Col(col), math.Max(now, toSet))
		case HEIGHT:
			now, _ := file.GetRowHeight(sheetname, row)
			toSet, _ := strconv.ParseFloat(v, 64)
			file.SetRowHeight(sheetname, row, math.Max(now, toSet))
		case SUFFIX:
			content, _ := file.GetCellValue(sheetname, pos(row, col))
			file.SetCellValue(sheetname, pos(row, col), content+v)
		case READCONVERTEXP:
			content, _ := file.GetCellValue(sheetname, pos(row, col))
			file.SetCellValue(sheetname, pos(row, col), readConvert(v, content))
		}
	}
}

func readConvert(expiration string, content string) string {
	for _, kv := range strings.Split(expiration, ",") {
		str := strings.Split(kv, "=")
		if len(str) == 2 && str[0] == content {
			return str[1]
		}
	}
	return content
}

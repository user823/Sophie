package utils

import (
	"github.com/user823/Sophie/internal/schedule/predefined"
	"github.com/user823/Sophie/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

func CallTarget(invokeTarget string) (err error) {
	funcName, args := parseArgs(invokeTarget)
	argValues := make([]reflect.Value, len(args))

	if fn, ok := predefined.Targets[funcName]; ok {
		fnType := reflect.TypeOf(fn)
		for i := 0; i < fnType.NumIn(); i++ {
			argType := fnType.In(i)
			ptr := reflect.New(argType)
			argValue := ptr.Elem()
			switch argType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v, _ := strconv.ParseInt(args[i], 10, 64)
				argValue.SetInt(v)
			case reflect.Bool:
				v, _ := strconv.ParseBool(args[i])
				argValue.SetBool(v)
			case reflect.Float64, reflect.Float32:
				v, _ := strconv.ParseFloat(args[i], 64)
				argValue.SetFloat(v)
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint8, reflect.Uint64:
				v, _ := strconv.ParseUint(args[i], 10, 64)
				argValue.SetUint(v)
			case reflect.Complex64, reflect.Complex128:
				v, _ := strconv.ParseComplex(args[i], 128)
				argValue.SetComplex(v)
			case reflect.String:
				argValue.SetString(args[i])
			default:
				return errors.New("不支持的参数类型，仅支持整形、字符串、浮点型、布尔、复数类型 ")
			}
			argValues[i] = argValue
		}
		fnValue := reflect.ValueOf(fn)
		fnValue.Call(argValues)
		return nil
	}
	return errors.New("未找到调用目标")
}

func parseArgs(invokeTarget string) (string, []string) {
	// 去除首尾的空格
	invokeTarget = strings.TrimSpace(invokeTarget)

	// 查找函数名
	idx := strings.Index(invokeTarget, "(")
	if idx == -1 {
		return "", nil
	}
	funcName := invokeTarget[:idx]

	// 查找参数列表
	argsStr := invokeTarget[idx+1 : len(invokeTarget)-1]
	var result []string
	for argsStr != "" {
		argsStr = strings.Trim(argsStr, " ,")
		if argsStr == "" {
			break
		}
		if strings.Index(argsStr, "\"") == 0 {
			next := 1 + strings.Index(argsStr[1:], "\"")
			result = append(result, argsStr[:next+1])
			argsStr = argsStr[next+1:]
		} else if strings.Contains(argsStr, ",") {
			next := strings.Index(argsStr, ",")
			result = append(result, argsStr[:next])
			argsStr = argsStr[next+1:]
		} else {
			result = append(result, argsStr)
			argsStr = ""
		}
	}
	return funcName, result
}

// 检查调用目标参数是否合法
func CheckTarget(invokeTarget string) (err error) {
	funcName, args := parseArgs(invokeTarget)
	if funcName == "" {
		return errors.New("调用目标格式设置不正确")
	}
	if fn, ok := predefined.Targets[funcName]; ok {
		fnType := reflect.TypeOf(fn)
		if len(args) != fnType.NumIn() {
			return errors.New("调用参数设置不正确")
		}
		return nil
	}
	return errors.New("未找到调用目标")
}

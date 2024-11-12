package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// IsDateTimeDtFmt ...
// var := IsDateTimeDtFmt("20210131" + "000000") -- true : valid date, false : in valid dat
func IsDateTimeDtFmt(DtFmt string) bool {

	_, err := time.Parse("20060102150405", DtFmt)

	// fmt.Println(err)
	// fmt.Println(convTime)
	// fmt.Println(convTime.IsZero())

	return err == nil
}

// IsDateDtFmt ...
func IsDateDtFmt(DtFmt string) bool {

	_, err := time.Parse("20060102", DtFmt)

	// fmt.Println(err)
	// fmt.Println(convTime)
	// fmt.Println(convTime.IsZero())

	return err == nil
}

// SplitIndex ...
// SplitIndex("A|B|C", "|", 0) == "A"
// SplitIndex("A|B|C", "|", 1) == "B"
// SplitIndex("A|B|C", "|", 2) == "C"
// SplitIndex("A|B|C", "|", 3) == ""
func SplitIndex(s string, sep string, index int) string {
	if index < 0 {
		return ""
	}

	tempStr := strings.Split(s, sep)

	if len(tempStr) <= index {
		return ""
	}

	return tempStr[index]
}

// Min ...
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max ...
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// StringEllipsis ...
func StringEllipsis(str string, start int, length int) (string, bool) {

	//lll := len(str)
	//var r rune = str
	//lll := utf8.RuneCountInString(str)
	maxStr := Min(utf8.RuneCountInString(str), length)

	//fmt.Println(maxStr)

	return strings.TrimSpace(string([]rune(str)[start : maxStr+start])), (utf8.RuneCountInString(str) > length)
}

// ConvStringEllipsis ...
func ConvStringEllipsis(str string, start int, length int) string {

	ellStr, ell := StringEllipsis(str, start, length)

	var addStr string
	if ell == true {
		addStr = "..."
	} else {
		addStr = ""
	}

	return ellStr + addStr
}

func parseAndSet(target reflect.Value, val string) error {
	if !target.CanSet() {
		return fmt.Errorf("Cannot set %v to %v", target, val)
	}

	switch kind := target.Type().Kind(); kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			target.SetInt(intVal)
		}
		return err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		intVal, err := strconv.ParseUint(val, 10, 64)
		if err == nil {
			target.SetUint(intVal)
		}
		return err
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err == nil {
			target.SetFloat(floatVal)
		}
		return err
	case reflect.String:
		target.SetString(val)
		return nil
	}
	return fmt.Errorf("Field %v has type %v, cannot set to %v", target, target.Type(), val)
}

// func Invoke2(any interface{}, name string, args ...interface{}) {
// 	inputs := make([]reflect.Value, len(args))
// 	for i, _ := range args {
// 		inputs[i] = reflect.ValueOf(args[i])
// 	}

// 	reflect.ValueOf(any).MethodByName(name).Call(inputs)
// }

func Invoke(any interface{}, name string, args ...interface{}) (reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)

	if !method.IsValid() {
		return reflect.ValueOf(nil), fmt.Errorf("Not Found Method %s. ", name)
	}

	if method.IsNil() {
		return reflect.ValueOf(nil), fmt.Errorf("Is Nil Method %s. ", name)
	}

	if method.IsZero() {
		return reflect.ValueOf(nil), fmt.Errorf("Is Zero Method %s. ", name)
	}

	methodType := method.Type()
	numIn := methodType.NumIn()
	if numIn > len(args) {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have minimum %d params. Have %d", name, numIn, len(args))
	}
	if numIn != len(args) && !methodType.IsVariadic() {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have %d params. Have %d", name, numIn, len(args))
	}
	in := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		var inType reflect.Type
		if methodType.IsVariadic() && i >= numIn-1 {
			inType = methodType.In(numIn - 1).Elem()
		} else {
			inType = methodType.In(i)
		}
		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argValue.String())
		}
		argType := argValue.Type()
		if argType.ConvertibleTo(inType) {
			in[i] = argValue.Convert(inType)
		} else {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argType)
		}
	}
	return method.Call(in)[0], nil
}

// LoadConfigFile ...
func LoadConfigFile(filename string, cfg interface{}) int {

	// value := reflect.TypeOf(cfg)
	// fmt.Println("1st Value Type: ", value.Kind())

	//fmt.Println(reflect.ValueOf(cfg).Type())
	//cc := cfg.(reflect.TypeOf(cfg).Type())
	//cc := cfg.(tables.TableConfig)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf(fmt.Sprintf("[LoadConfiFile][Error] %s", err))
		return 1
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	var lineNum int = 0
	var errCnt int = 0
	for _, eachline := range txtlines {
		lineNum++

		if eachline == "" {
			continue
		}

		if eachline[0] == '#' {
			continue
		}

		splitString := strings.Split(eachline, "=")
		if len(splitString) != 2 {
			errCnt++
			fmt.Printf(fmt.Sprintf("[LoadConfiFile][%d][Error] %s -> invalid set", lineNum, eachline))
			continue
		}

		category := strings.TrimSpace(splitString[0])
		value := strings.TrimSpace(splitString[1])

		if len(category) <= 0 || len(value) <= 0 {
			errCnt++
			fmt.Printf(fmt.Sprintf("[LoadConfiFile][%d][Error] %s=%s -> category is nill or value is nill", lineNum, category, value))
			continue
		}

		v := reflect.ValueOf(cfg).Elem()

		var err error

		// field, ok := reflect.TypeOf(&cfg).Elem().FieldByName("LiveNvnApplyMaxCount") // not json:name
		// tag := string(field.Tag)
		// fmt.Print(tag)
		// fmt.Print(ok)

		err = parseAndSet(v.FieldByName(category), value)
		if err != nil {
			errCnt++
			fmt.Printf(fmt.Sprintf("[LoadConfiFile][%d][Error] %s=%s -> %s", lineNum, category, value, err))
		} else {
			fmt.Printf(fmt.Sprintf("[LoadConfiFile][%d][Load] %s=%s", lineNum, category, value))
		}

		// result := v.MethodByName("IsCheckValue").Call([]reflect.Value{})
		// if result != "" {
		// 	errCnt++
		// 	fmt.Printf(fmt.Sprintf("[LoadConfiFile][%d][Error] %s=%s -> %s inValid Value", lineNum, category, value, result))
		// }
	}

	return errCnt
}

// ConvertEucKR ...
func ConvertEucKR(codeStr string) string {
	var convCode string = ""

	var bufs bytes.Buffer
	wr := transform.NewWriter(&bufs, korean.EUCKR.NewDecoder())
	wr.Write([]byte(codeStr))
	wr.Close()
	convCode = bufs.String()

	return convCode
}

func StringDelimSplit(str string, paramChar string, paramValueChar string) map[string]string {
	entries := strings.Split(str, paramChar)

	m := make(map[string]string)
	for _, e := range entries {
		parts := strings.Split(e, paramValueChar)
		m[parts[0]] = parts[1]
	}

	return m
}

func GetStackTrace() string {
	var retStr = ""
	buf := make([]byte, 1<<20)

	stacklen := runtime.Stack(buf, true)
	retStr = string(buf[:stacklen])

	return retStr
}

func ToJsonString(v interface{}) string {
	conv, err := json.Marshal(v)
	if err != nil {
		logs.Error(fmt.Sprintf("[ConvertObjectToJson][Error] %v => %v ", v, err))
		return string(conv)
	}

	return string(conv)
}

// uniqueIntSlice := unique(intSlice).([]int)
func RemoveDuplicateValuesOld(src interface{}) interface{} {
	srcv := reflect.ValueOf(src)
	dstv := reflect.MakeSlice(srcv.Type(), 0, 0)
	visited := make(map[interface{}]struct{})
	for i := 0; i < srcv.Len(); i++ {
		elemv := srcv.Index(i)
		if _, ok := visited[elemv.Interface()]; ok {
			continue
		}
		visited[elemv.Interface()] = struct{}{}
		dstv = reflect.Append(dstv, elemv)
	}
	return dstv.Interface()
}

// testArray := []int{1,5,3,6,9,9,4,2,3,1,5}
// testArrayDup := make([]int, 0)

// err = utils.RemoveDuplicateValues(testArray, &testArrayDup)
// if err != nil {
// 	fmt.Printf(fmt.Sprintf("error filtering int slice: %v\n", err))
// }
func RemoveDuplicateValues(slice interface{}, filtered interface{}) error {

	// Check for slice of string
	if sliceOfString, ok := slice.([]string); ok {

		// If slice is slice of string filtered MUST also be slice of string
		filteredAsSliceOfString, ok := filtered.(*[]string)
		if !ok {
			return fmt.Errorf("filtered should be of type %T, got %T instead", &[]string{}, filtered)
		}

		keys := make(map[string]bool)

		for _, entry := range sliceOfString {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				*filteredAsSliceOfString = append(*filteredAsSliceOfString, entry)
			}
		}

	} else if sliceOfInt, ok := slice.([]int); ok {

		// If slice is slice of int filtered MUST also be slice of int
		filteredAsInt, ok := filtered.(*[]int)
		if !ok {
			return fmt.Errorf("filtered should be of type %T, got %T instead", &[]string{}, filtered)
		}

		keys := make(map[int]bool)

		for _, entry := range sliceOfInt {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				*filteredAsInt = append(*filteredAsInt, entry)
			}
		}

	} else {
		return fmt.Errorf("only slice of in or slice of string is supported")
	}

	return nil
}

func DtToFmtDt(dt string, dtFmt string) string {

	convDt, err := time.Parse("20060102150405", dt)
	if err != nil {
		return dt
	}

	return fmt.Sprintf(convDt.Format(dtFmt))
}

func CdnPathToImgSvrPath(cdnPath string, imgServer string, regDt string) string {

	convDt, err := time.Parse("20060102150405", regDt)
	if err != nil {
		return ""
	}

	dir, file := filepath.Split(cdnPath)

	return fmt.Sprintf("%v%vtemp_%v%v", imgServer, dir, convDt.Year(), file)
}

func DbRowToInt64(value interface{}, defValue int64) int64 {
	var resultValue int64 = defValue

	if value != nil {
		resultValue = value.(int64)
	}  						

	return resultValue;
}

func DbRowToString(value interface{}, defValue string) string {
	var resultValue string = defValue

	if value != nil {
		resultValue = value.(string)
	}  						

	return resultValue;
}
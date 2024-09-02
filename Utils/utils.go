package utils

import (
	"UmaAIChatServer/Utils/logx"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"math/rand"
)

const tag = "Utils"

var LangList = [...]string{
	"日本語",
	"简体中文",
	"English",
}

var TrList = [...]string{
	"Chinese",
	"English",
}

func ReadOrCreateFile(filename string, content []byte) ([]byte, error) {
	if _, err := os.Stat(filename); err == nil {
		data, err := os.ReadFile(filename)
		if err != nil {
			return []byte{}, fmt.Errorf("the file could not be read. 不能读取文件: %v", err)
		}
		return data, nil
	} else if os.IsNotExist(err) {
		err = os.WriteFile(filename, []byte(content), 0755)
		if err != nil {
			return []byte{}, fmt.Errorf("the file could not be created. 不能创建文件: %v", err)
		}
		return content, nil
	} else {
		return []byte{}, fmt.Errorf("error detecting the file. 检测文件出错: %v", err)
	}
}

func EasyFileRead(path string) (bool, []byte) {
	data, err := os.ReadFile(path)
	if err != nil {
		logx.Warn(tag, err)
		return false, nil
	}
	return true, data
}

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func GenerateID(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func StrToInt64(number string) int64 {
	out, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		logx.Warn(tag, "转换错误", number, err)
		return math.MinInt64
	}
	return out
}

func Int64ToStr(number int64) string {
	return strconv.FormatInt(number, 10)
}

func Str2Int(number string) int {
	int, err := strconv.Atoi(number)
	if err != nil {
		return math.MinInt
	}
	return int
}

func Int2Str(number int) string {
	return strconv.Itoa(number)
}

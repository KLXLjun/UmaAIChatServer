package utils

import (
	"fmt"
	"os"
)

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

package util

import (
	"os"
	"log"
	"io/ioutil"
	"strings"
	"path/filepath"
)

//判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建文件夹
func CreatedDir(dir string) {
	exist, err := PathExists(dir)
	if err != nil {
		log.Printf("get dir error![%v]\n", err)
		return
	}

	if exist {
		log.Printf("has dir![%v]\n", dir)
	} else {
		log.Printf("no dir![%v]\n", dir)
		// 创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Printf("mkdir failed![%v]\n", err)
		} else {
			log.Printf("mkdir success!\n")
		}
	}
}

func Yamls2Bytes(rootPath string, files []string) [][]byte {
	yamls := make([][]byte, len(files))
	for i, name := range files {
		yamlBytes, err := ioutil.ReadFile(filepath.Join(rootPath, name))
		if err != nil {
			log.Println(err.Error())
		}
		yamls[i] = yamlBytes

	}
	return yamls
}

func FirstUpper(org string) string {
	return strings.ToUpper(org[:1]) + org[1:]
}

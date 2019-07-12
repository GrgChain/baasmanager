package util

import (
	"strings"
	"io/ioutil"
	"path/filepath"
	"log"
	"bytes"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//首字母大写
func FirstUpper(org string) string {
	return strings.ToUpper(org[:1]) + org[1:]
}

//yaml to 字节数组
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

//yaml to json
func Yamls2Jsons(data [][]byte) [][]byte {
	jsons := make([][]byte, 0)
	for _, yamlBytes := range data {
		yamls := bytes.Split(yamlBytes, []byte("---"))
		for _, v := range yamls {
			if len(v) == 0 {
				continue
			}
			obj, err := yaml.ToJSON(v)
			if err != nil {
				log.Println(err.Error())
			}
			jsons = append(jsons, obj)
		}

	}
	return jsons
}
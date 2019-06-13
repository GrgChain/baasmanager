package util

import (
	"net/http"
	"strings"
	"io/ioutil"
	"net/url"
	"log"
	"encoding/json"
)

func PostJson(uri string, jsons interface{}) []byte {

	jsonParam, errs := json.Marshal(jsons) //转换成JSON返回的是byte[]
	if errs != nil {
		log.Println(errs.Error())
	}

	//发送请求
	req, err := http.NewRequest("POST", uri, strings.NewReader(string(jsonParam)))
	if err != nil {
		log.Println(err.Error())
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()
	//响应
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read failed:", err)
	}

	//返回结果
	return response
}

func PostForm(uri string, paras map[string][]string) []byte {

	resp, err := http.PostForm(uri, url.Values(paras))
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	return body

}

func Get(uri string) []byte {

	resp, err := http.Get(uri)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}

	return body

}

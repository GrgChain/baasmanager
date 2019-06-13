package httputil

import (
	"net/http"
	"strings"
	"io/ioutil"
	"github.com/jonluo94/commontools/log"
	"net/url"
	"github.com/jonluo94/commontools/json"
)

var logger = log.GetLogger("httputil", log.ERROR)

func PostJson(uri string, jsons interface{}) []byte {

	jsonParam, errs := json.Marshal(jsons) //转换成JSON返回的是byte[]
	if errs != nil {
		logger.Error(errs.Error())
	}

	//发送请求
	req, err := http.NewRequest("POST", uri, strings.NewReader(string(jsonParam)))
	if err != nil {
		logger.Error(err.Error())
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
	}
	defer resp.Body.Close()
	//响应
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Read failed:", err)
	}

	//返回结果
	return response
}

func PostForm(uri string, paras map[string][]string) []byte {

	resp, err := http.PostForm(uri, url.Values(paras))
	if err != nil {
		logger.Error(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
	}
	return body

}

func Get(uri string) []byte {

	resp, err := http.Get(uri)
	if err != nil {
		logger.Error(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
	}

	return body

}

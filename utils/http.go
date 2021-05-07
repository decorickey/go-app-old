package utils

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// DoSimpleGet 簡単にGETしたいとき
func DoSimpleGet(endpoint string) []byte {
	res, _ := http.Get(endpoint)
	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	return resBody
}

// DoSimplePost 簡単にPOSTしたいとき
func DoSimplePost(endpoint string, params map[string]string) []byte {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}
	res, _ := http.PostForm(endpoint, values)
	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	return resBody
}

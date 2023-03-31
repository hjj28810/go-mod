package request

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/hjj28810/go-mod/log"
	"github.com/hjj28810/go-mod/utility"
)

func DoHttpBase(url string, method string, data any, headers map[string]string) (result io.ReadCloser) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, bytes.NewReader(utility.ToJsonBody(data)))
	if err != nil {
		log.ErrorLogAsync("http request err", "", err)
		return
	}
	if _, ok := headers["Content-Type"]; !ok {
		request.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	if len(headers) > 0 {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}
	resp, err := client.Do(request)
	if err != nil {
		log.ErrorLogAsync("http do response err", "", err)
		return
	}

	return resp.Body
}

func DoHttpGen[T any](url string, method string, data any, headers map[string]string) T {
	respBody := DoHttpBase(url, method, data, headers)
	body, err := io.ReadAll(respBody)
	defer respBody.Close()
	if err != nil {
		panic(err)
	}
	// str := string(body)
	// fmt.Println(str)
	return utility.JsonBodyToObj[T](body)
}

func DoHttp(url string, method string, data any, headers map[string]string) string {
	respBody := DoHttpBase(url, method, data, headers)
	body, err := io.ReadAll(respBody)
	defer respBody.Close()
	if err != nil {
		panic(err)
	}
	return string(body)
}

type GetRequest struct {
	UrlValues url.Values
}

func (p *GetRequest) Init() *GetRequest {
	p.UrlValues = url.Values{}
	return p
}

func (p *GetRequest) InitFrom(reqParams *GetRequest) *GetRequest {
	if reqParams != nil {
		p.UrlValues = reqParams.UrlValues
	} else {
		p.UrlValues = url.Values{}
	}
	return p
}

func (p *GetRequest) AddParam(property string, value string) *GetRequest {
	if property != "" && value != "" {
		p.UrlValues.Add(property, value)
	}
	return p
}

func (p *GetRequest) BuildParams() string {
	return p.UrlValues.Encode()
}

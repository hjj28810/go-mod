package request

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hjj28810/go-mod/log"
	"github.com/hjj28810/go-mod/utility"
)

func DoHttpBase(url string, method string, data any, headers map[string]string) (result io.ReadCloser) {
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json;charset=utf-8"
	}
	client := &http.Client{}
	request, err := http.NewRequest(method, url, DataReader(data))
	if err != nil {
		log.ErrorLogAsync("http request err", "", err)
		return
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

func DataReader(data any) io.Reader {
	var reader io.Reader
	switch data := data.(type) {
	case string:
		reader = StringReader(data)
	default:
		reader = JsonReader(data)
	}
	return reader
}

func JsonReader(data any) *bytes.Reader {
	return bytes.NewReader(utility.ToJsonBody(data))
}

func StringReader(data string) *strings.Reader {
	return strings.NewReader(data)
}

package request

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/hjj28810/go-mod/log"
	"github.com/hjj28810/go-mod/utility"
)

func DoHttpBase(url string, method string, data any, headers map[string]string) (result io.ReadCloser) {
	if _, ok := headers["Content-Type"]; headers != nil && !ok {
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
	if resp.StatusCode != http.StatusOK {
		log.WarningLogAsync("http do response status", resp.Status)
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Status)
		return
	}
	if err != nil {
		log.ErrorLogAsync("http do response err", "", err)
		return
	}
	return resp.Body
}

func DownloadFile(url string, ext string) string {
	file, err := os.CreateTemp("", "file-*"+ext)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()
	reader := DoHttpBase(url, "GET", nil, nil)
	defer reader.Close()
	_, err = io.Copy(file, reader)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return file.Name()
}

func DoHttpGen[T any](url string, method string, data any, headers map[string]string) T {
	respBody := DoHttpBase(url, method, data, headers)
	body, err := io.ReadAll(respBody)
	defer respBody.Close()
	if err != nil {
		panic(err)
	}
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

func BinaryReader(data string) *strings.Reader {
	return strings.NewReader(data)
}

type MultipartFormField struct {
	IsMedia bool
	Name    string
	Value   string
}

func RequestMultiPart[T any](method, url string, fields []MultipartFormField) T {
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)

	for _, field := range fields {
		if field.IsMedia {
			fileWriter, e := writer.CreateFormFile(field.Name, field.Value)
			if e != nil {
				panic(e)
			}
			fh, e := os.Open(field.Value)
			if e != nil {
				panic(e)
			}
			defer fh.Close()

			if _, e = io.Copy(fileWriter, fh); e != nil {
				panic(e)
			}
		} else {
			partWriter, e := writer.CreateFormField(field.Name)
			if e != nil {
				panic(e)
			}
			valueReader := bytes.NewReader(utility.ToJsonBody(field.Value))
			if _, e = io.Copy(partWriter, valueReader); e != nil {
				panic(e)
			}
		}
	}

	// Close the form
	writer.Close()

	// Create a new HTTP request with the form data
	req, err := http.NewRequest(method, url, bodyBuf)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return utility.JsonBodyToObj[T](respBody)
}

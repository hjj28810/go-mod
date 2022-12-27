package utility

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/hjj28810/go-mod/log"
)

func IntArrJoinString(arr []int, spilt string) string {
	var temp = make([]string, len(arr))
	for i, v := range arr {
		temp[i] = strconv.Itoa(v)
	}
	return strings.Join(temp, spilt)
}

func StrToIntArr(str string) []int {
	arr := strings.Split(str, ",")
	var temp = make([]int, len(arr))
	for i, v := range arr {
		intV, _ := strconv.Atoi(v)
		temp[i] = intV
	}
	return temp
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ToJson(t interface{}) string {
	msgBody, err := json.Marshal(t)
	if nil != err {
		log.ErrorLogAsync("toJson fail", "", err)
	}
	return string(msgBody)
}

func ToJsonBody(t interface{}) []byte {
	msgBody, err := json.Marshal(t)
	if nil != err {
		log.ErrorLogAsync("toJsonBody fail", "", err)
	}
	return msgBody
}

func JsonToObj[T any](str string) (result T) {
	var data []byte = []byte(str)
	err := json.Unmarshal(data, &result)
	if nil != err {
		log.ErrorLogAsync("jsonToObj fail", "", err)
	}
	return result
}

func XMLToObj[T any](str string) (result T) {
	var data []byte = []byte(str)
	err := xml.Unmarshal(data, &result)
	if nil != err {
		log.ErrorLogAsync("XMLToObj fail", "", err)
	}
	return result
}

func JsonBodyToObj[T any](data []byte) (result T) {
	err := json.Unmarshal(data, &result)
	if nil != err {
		log.ErrorLogAsync("jsonBodyToObj fail", "", err)
	}
	return result
}

func XMLBodyToObj[T any](data []byte) (result T) {
	err := xml.Unmarshal(data, &result)
	if nil != err {
		log.ErrorLogAsync("XMLBodyToObj fail", "", err)
	}
	return result
}

func SubString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length

	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}

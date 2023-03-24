package utility

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"reflect"
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
	return string(ToJsonBody(t))
}

func ToJsonBody(t interface{}) []byte {
	msgBody, err := json.Marshal(t)
	if nil != err {
		log.ErrorLogAsync("toJsonBody fail", "", err)
		panic(err)
	}
	return msgBody
}

func ToXML(t interface{}) string {
	return string(ToXMLBody(t))
}

func ToXMLBody(t interface{}) []byte {
	msgBody, err := xml.Marshal(t)
	if nil != err {
		log.ErrorLogAsync("toXMLBody fail", "", err)
		panic(err)
	}
	return msgBody
}

func JsonToObj[T any](str string) T {
	var data []byte = []byte(str)
	return JsonBodyToObj[T](data)
}

func XMLToObj[T any](str string) T {
	var data []byte = []byte(str)
	return XMLBodyToObj[T](data)
}

func JsonBodyToObj[T any](data []byte) (result T) {
	err := json.Unmarshal(data, &result)
	if nil != err {
		log.ErrorLogAsync("jsonBodyToObj fail", reflect.TypeOf(result).Name(), err)
		panic(err)
	}
	return result
}

func XMLBodyToObj[T any](data []byte) (result T) {
	err := xml.Unmarshal(data, &result)
	if nil != err {
		log.ErrorLogAsync("XMLBodyToObj fail", reflect.TypeOf(result).Name(), err)
		panic(err)
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

func SliceExcept[T any](slice1, slice2 []T, compare func(T, T) bool) []T {
	result := []T{}
	for _, v1 := range slice1 {
		found := false
		for _, v2 := range slice2 {
			if compare(v1, v2) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, v1)
		}
	}
	return result
}

func SliceIntersect[T any](slice1, slice2 []T, compare func(T, T) bool) []T {
	result := []T{}
	for _, v1 := range slice1 {
		for _, v2 := range slice2 {
			if compare(v1, v2) {
				result = append(result, v1)
				break
			}
		}
	}
	return result
}

func SliceUnion[T any](slice1, slice2 []T, compare func(T, T) bool) []T {
	result := []T{}
	result = append(result, slice1...)
	for _, v2 := range slice2 {
		exist := false
		for _, v1 := range slice1 {
			if compare(v1, v2) {
				exist = true
				break
			}
		}
		if !exist {
			result = append(result, v2)
		}
	}
	return result
}

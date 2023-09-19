package utility

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hjj28810/go-mod/log"
	"github.com/hjj28810/go-mod/model"
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
		str := string(data)
		fmt.Println(str)
		log.ErrorLogAsync("jsonBodyToObj fail--"+reflect.TypeOf(result).Name(), str, err)
		panic(err)
	}
	return result
}

func XMLBodyToObj[T any](data []byte) (result T) {
	err := xml.Unmarshal(data, &result)
	if nil != err {
		str := string(data)
		fmt.Println(str)
		log.ErrorLogAsync("XMLBodyToObj fail--"+reflect.TypeOf(result).Name(), str, err)
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

func StringToTime(timeValue string) time.Time {
	timestamp, err := time.ParseInLocation(model.TimeLayout, timeValue, time.Local)
	if err != nil {
		fmt.Println("时间转换错误", err)
	}
	return timestamp
}

func SubMonth(t1, t2 time.Time) (month int) {
	y1 := t1.Year()
	y2 := t2.Year()
	m1 := int(t1.Month())
	m2 := int(t2.Month())
	d1 := t1.Day()
	d2 := t2.Day()

	yearInterval := y1 - y2
	// 如果 d1的 月-日 小于 d2的 月-日 那么 yearInterval-- 这样就得到了相差的年数
	if m1 < m2 || m1 == m2 && d1 < d2 {
		yearInterval--
	}
	// 获取月数差值
	monthInterval := (m1 + 12) - m2
	if d1 < d2 {
		monthInterval--
	}
	monthInterval %= 12
	month = yearInterval*12 + monthInterval
	return
}

func RSADecrypt(ciphertext string, privateKeyArr []byte) string {
	// Parse the PEM-encoded private key
	block, _ := pem.Decode(privateKeyArr)
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	// Decrypt the ciphertext using the private key
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(ciphertext))
	if err != nil {
		fmt.Println("Error decrypting ciphertext:", err)
		return ""
	}
	return string(plaintext)
}

func RSAEncrypt(plaintext string, publicKeyArr []byte) string {
	block, _ := pem.Decode(publicKeyArr)
	publicKey, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plaintext))
	if err != nil {
		fmt.Println("Failed to encrypt message:", err)
		return ""
	}
	return string(ciphertext)
}

func ReadFile(path string) []byte {
	pemData, err := os.ReadFile("path/to/file.pem")
	if err != nil {
		fmt.Println("无法读取 PEM 文件:", err)
		return []byte("")
	}
	return pemData
}

package model

import "time"

type ResponseModel[T any] struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Data       T      `json:"data"`
	ServerTime int64  `json:"serverTime"`
}

func (res *ResponseModel[T]) WithMsg(message string) ResponseModel[T] {
	return ResponseModel[T]{
		Code:       res.Code,
		Msg:        message,
		Data:       res.Data,
		ServerTime: res.ServerTime,
	}
}

// 追加响应数据
func (res *ResponseModel[T]) WithData(data T) ResponseModel[T] {
	return ResponseModel[T]{
		Code:       res.Code,
		Msg:        res.Msg,
		Data:       data,
		ServerTime: res.ServerTime,
	}
}

func BaseResponse(code int, msg string) *ResponseModel[interface{}] {
	return &ResponseModel[interface{}]{
		Code:       code,
		Msg:        msg,
		ServerTime: time.Now().Unix(),
	}
}

var (
	ResponseOK  = BaseResponse(200, "ok")  // 通用成功
	ResponseErr = BaseResponse(500, "err") // 通用错误
)

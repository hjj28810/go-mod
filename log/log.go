package log

import (
	"context"

	"github.com/gogf/gf/v2/os/glog"
)

var ctx context.Context
var Path string

func init() {
	ctx = context.TODO()
	glog.SetWriterColorEnable(true)
}

func SetPath(path string) {
	glog.SetPath(path)
}

func DebugLog(title string, content string) {
	glog.Debug(ctx, LogInfo{Title: title, Content: content})
}

func ErrorLog(title string, content string, logErr error) {
	glog.Error(ctx, LogInfo{Title: title, Content: content, LogErr: logErr})
}

func InfoLog(title string, content string) {
	glog.Info(ctx, LogInfo{Title: title, Content: content})
}

func WarningLog(title string, content string) {
	glog.Warning(ctx, LogInfo{Title: title, Content: content})
}

type LogInfo struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	LogErr  error  `json:"logErr"`
}

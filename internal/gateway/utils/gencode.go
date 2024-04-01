package utils

import (
	"bytes"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/pkg/log"
	"io"
	"strconv"
)

func GenCode(c *app.RequestContext, data []byte) {
	c.Response.Reset()
	c.Response.Header.Add("Content-Disposition", "attachment; filename=\"sophie.zip\"")
	c.Response.Header.Add("Content-Length", strconv.Itoa(len(data)))
	c.Response.Header.SetContentType("application/octet-stream; charset=UTF-8")
	buffer := bytes.NewReader(data)
	if _, err := io.Copy(c, buffer); err != nil {
		log.Error("下载文件失败: %s", err.Error())
	}
}

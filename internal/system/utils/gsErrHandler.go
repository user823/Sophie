package utils

import (
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/shutdown"
)

// 将关停过程中的错误打印出来
type GsLogErrHandler struct {
	next shutdown.ErrHandler
}

func (g *GsLogErrHandler) OnError(err error) {
	log.Error(err)
}

func (g *GsLogErrHandler) SetDeliver(eh shutdown.ErrHandler) {
	g.next = eh
}

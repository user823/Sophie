package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/pkg/errors"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/pkg/code"
	"sync"

	"github.com/user823/Sophie/pkg/log"
)

type RPCClient struct {
	*SystemClient
	*JobClient
	*FileClient
	*GenClient
}

var (
	once     sync.Once
	Remoting *RPCClient
)

func Init(generalOpts []client.Option) {
	once.Do(func() {
		// ---使用cfg 初始化各个client组件---
		Remoting = &RPCClient{
			SystemClient: newSystemClient(),
			JobClient:    newJobClient(),
			FileClient:   newFileClient(),
			GenClient:    newGenClient(),
		}

		Remoting.SystemClient.initRPC(generalOpts)
		Remoting.JobClient.initRPC(generalOpts)
		Remoting.FileClient.initRPC(generalOpts)
		Remoting.GenClient.initRPC(generalOpts)
	})
}

func ParseRpcErr(b *v1.BaseResp, err error) error {
	if err != nil {
		log.Errorf("rpc invoke error: %s", err.Error())
		return errors.New("系统内部错误，请重试")
	}

	if b != nil && b.Code == code.ERROR {
		return errors.WithMessage(err, b.Msg)
	}
	return nil
}

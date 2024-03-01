package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/api/thrift/system/v1/systemservice"
	"github.com/user823/Sophie/pkg/log"
)

type SystemClient struct {
	systemservice.Client
}

// opts: 通用的options
func (c *SystemClient) initRPC(opts []client.Option) {
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: v1.ServiceName}))

	cli, err := systemservice.NewClient(v1.ServiceName, opts...)
	if err != nil {
		log.Errorf("%s remoting client init err: %s", v1.ServiceName, err.Error())
		panic(err)
	}
	c.Client = cli
}

func newSystemClient() *SystemClient {
	return &SystemClient{}
}

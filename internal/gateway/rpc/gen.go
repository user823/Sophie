package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	v12 "github.com/user823/Sophie/api/domain/gen/v1"
	v1 "github.com/user823/Sophie/api/thrift/gen/v1"
	"github.com/user823/Sophie/api/thrift/gen/v1/genservice"
	"github.com/user823/Sophie/pkg/log"
)

type GenClient struct {
	genservice.Client
}

func (c *GenClient) initRPC(opts []client.Option) {
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: v12.ServiceName}))

	cli, err := genservice.NewClient(v1.ServiceName, opts...)
	if err != nil {
		log.Errorf("%s remoting client init err: %s", v12.ServiceName, err.Error())
		panic(err)
	}
	c.Client = cli
}

func newGenClient() *GenClient {
	return &GenClient{}
}

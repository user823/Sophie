package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/user823/Sophie/api/domain/file/v1"
	v12 "github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/api/thrift/file/v1/fileservice"
	"github.com/user823/Sophie/pkg/log"
)

type FileClient struct {
	fileservice.Client
}

func (c *FileClient) initRPC(opts []client.Option) {
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: v12.ServiceName}))

	cli, err := fileservice.NewClient(v12.ServiceName, opts...)
	if err != nil {
		log.Errorf("%s remoting client init err: %s", v1.ServiceName, err.Error())
		panic(err)
	}
	c.Client = cli
}

func newFileClient() *FileClient {
	return &FileClient{}
}

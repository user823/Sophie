package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/user823/Sophie/api/domain/job/v1"
	"github.com/user823/Sophie/api/thrift/job/v1/jobservice"
	"github.com/user823/Sophie/pkg/log"
)

type JobClient struct {
	jobservice.Client
}

func (c *JobClient) initRPC(opts []client.Option) {
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: v1.ServiceName}))

	cli, err := jobservice.NewClient(v1.ServiceName, opts...)
	if err != nil {
		log.Errorf("%s remoting client init err: %s", v1.ServiceName, err.Error())
		panic(err)
	}
	c.Client = cli
}

func newJobClient() *JobClient {
	return &JobClient{}
}

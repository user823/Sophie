package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/user823/Sophie/api/thrift/schedule/v1/workerservice"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/pkg/log"
)

type WorkerClient struct {
	workerservice.Client
}

func (w *WorkerClient) initRPC(opts []client.Option) {
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: models.NodesPrefix}))

	cli, err := workerservice.NewClient(models.NodesPrefix, opts...)
	if err != nil {
		log.Errorf("%s remoting client init err: %s", models.NodesPrefix, err.Error())
		panic(err)
	}
	w.Client = cli
}

func newWorkerClient() *WorkerClient {
	return &WorkerClient{}
}

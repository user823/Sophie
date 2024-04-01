package service

import (
	"context"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/schedule/v1"
	v12 "github.com/user823/Sophie/api/thrift/schedule/v1"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/internal/schedule/manager/loadbalance"
	"github.com/user823/Sophie/internal/schedule/manager/rpc"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/internal/schedule/store"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"strconv"
)

type JobSrv interface {
	// 获取定时任务列表
	SelectJobList(ctx context.Context, job *v1.SysJob, opts *api.GetOptions) *v1.JobList
	// 通过调度任务ID查询调度信息
	SelectJobById(ctx context.Context, job *v1.SysJob, opts *api.GetOptions) *v1.SysJob
	// 暂停任务
	PauseJob(ctx context.Context, job *v1.SysJob) error
	// 恢复任务
	ResumeJob(ctx context.Context, job *v1.SysJob) error
	// 删除任务
	DeleteJob(ctx context.Context, job *v1.SysJob, opts *api.DeleteOptions) error
	// 批量删除调度信息
	DeleteJobByIds(ctx context.Context, jobIds []int64, opts *api.DeleteOptions) error
	// 任务调度状态修改
	ChangeStatus(ctx context.Context, job *v1.SysJob, opts *api.UpdateOptions) error
	// 立即执行任务
	Run(ctx context.Context, job *v1.SysJob) bool
	// 新增任务
	InsertJob(ctx context.Context, job *v1.SysJob, opts *api.CreateOptions) error
	// 重置任务，（将同一个任务组的任务重置）
	UpdateJob(ctx context.Context, job *v1.SysJob, opts *api.UpdateOptions) error
}

type JobService struct {
	store    store.Factory
	jobMp    models.JobMap
	nodePool models.NodePool
	lb       loadbalance.LoadBalancer
}

func NewJobs(s store.Factory, jm models.JobMap, np models.NodePool, bl loadbalance.LoadBalancer) JobSrv {
	return &JobService{s, jm, np, bl}
}

func (s *JobService) SelectJobList(ctx context.Context, job *v1.SysJob, opts *api.GetOptions) *v1.JobList {
	result, total, _ := s.store.Jobs().SelectJobList(ctx, job, opts)
	return &v1.JobList{
		ListMeta: api.ListMeta{total},
		Items:    result,
	}
}

func (s *JobService) SelectJobById(ctx context.Context, job *v1.SysJob, opts *api.GetOptions) *v1.SysJob {
	result, err := s.store.Jobs().SelectJobById(ctx, job.JobId, opts)
	if err != nil {
		return nil
	}
	return result
}

func (s *JobService) PauseJob(ctx context.Context, job *v1.SysJob) error {
	// 如果当前任务正在某个实例上运行则暂停任务
	if nodeId, ok := s.jobMp.Exists(job.JobId); ok {
		// 获取任务节点
		node, err := s.nodePool.GetNode(nodeId)

		if err != nil {
			log.Errorf("获取任务节点失败: %s", err.Error())
			s.jobMp.Delete(job.JobId)
		} else {
			resp, err := rpc.Remoting.PauseJobs(ctx, []int64{job.JobId}, callopt.WithHostPort(node.IPAddress))
			if err != nil {
				log.Debugf("rpc调用失败: %s", err.Error())
				return rpc.ErrRPC
			}

			if resp.Code != code.SUCCESS {
				log.Errorf(resp.Msg)
				return errors.New("系统内部错误")
			}
		}
	}

	job.Status = v1.PAUSE
	err := s.store.Jobs().UpdateJob(ctx, job, &api.UpdateOptions{})
	if err != nil {
		log.Debugf("暂停任务失败: %s", err.Error())
		return errors.New("系统内部错误")
	}
	return nil
}

func (s *JobService) ResumeJob(ctx context.Context, job *v1.SysJob) error {
	// 如果当前任务在某个实例上运行则恢复任务
	if nodeId, ok := s.jobMp.Exists(job.JobId); ok {
		// 获取任务节点
		node, err := s.nodePool.GetNode(nodeId)
		// 获取任务节点失败，重新进行调度
		if err != nil {
			log.Errorf("获取任务节点失败: %s", err.Error())
			s.jobMp.Delete(job.JobId)
			return s.Schedule(ctx, job)
		}

		resp, err := rpc.Remoting.ResumeJobs(ctx, []int64{job.JobId}, callopt.WithHostPort(node.IPAddress))
		if err != nil {
			log.Errorf("rpc调用失败: %s", err.Error())
			return rpc.ErrRPC
		}

		if resp.Code != code.SUCCESS {
			return errors.New(resp.Msg)
		}
	} else {
		// 当前任务不在任何一个实例上运行， 则创建调度任务
		if err := s.Schedule(ctx, job); err != nil {
			return err
		}
	}

	job.Status = v1.NORMAL
	err := s.store.Jobs().UpdateJob(ctx, job, &api.UpdateOptions{})
	if err != nil {
		log.Errorf("恢复任务失败: %s", err.Error())
		return errors.New("系统内部错误")
	}
	return nil
}

func (s *JobService) DeleteJob(ctx context.Context, job *v1.SysJob, opts *api.DeleteOptions) error {
	if err := s.store.Jobs().DeleteJobByIds(ctx, []int64{job.JobId}, opts); err != nil {
		return errors.New("系统内部错误")
	}

	// 如果任务在某个实例上运行则删除任务
	if nodeId, ok := s.jobMp.Exists(job.JobId); ok {
		// 获取任务节点
		node, err := s.nodePool.GetNode(nodeId)
		if err != nil {
			return err
		}

		resp, err := rpc.Remoting.RemoveJobs(ctx, []int64{job.JobId}, callopt.WithHostPort(node.IPAddress))
		if err != nil {
			return rpc.ErrRPC
		}

		if resp.Code != code.SUCCESS {
			return errors.New(resp.Msg)
		}

		return s.jobMp.Delete(job.JobId)
	}
	return nil
}

func (s *JobService) DeleteJobByIds(ctx context.Context, jobIds []int64, opts *api.DeleteOptions) error {
	for i := range jobIds {
		if err := s.DeleteJob(ctx, &v1.SysJob{JobId: jobIds[i]}, opts); err != nil {
			log.Debugf("删除任务失败: %s", err.Error())
		}
	}
	return nil
}

func (s *JobService) ChangeStatus(ctx context.Context, job *v1.SysJob, opts *api.UpdateOptions) error {
	if job.Status == v1.PAUSE {
		return s.PauseJob(ctx, job)
	}
	return s.ResumeJob(ctx, job)
}

func (s *JobService) Run(ctx context.Context, job *v1.SysJob) bool {
	// 查看任务信息
	jobinfo, err := s.store.Jobs().SelectJobById(ctx, job.JobId, &api.GetOptions{})
	if err != nil || jobinfo == nil {
		return false
	}

	// 查看任务信息
	if nodeId, ok := s.jobMp.Exists(jobinfo.JobId); ok {
		// 获取任务节点
		node, err := s.nodePool.GetNode(nodeId)
		if err != nil {
			return false
		}

		resp, err := rpc.Remoting.Run(ctx, []int64{jobinfo.JobId}, callopt.WithHostPort(node.IPAddress))
		if err != nil || resp.Code != code.SUCCESS {
			return false
		}
		return true
	}
	return false
}

func (s *JobService) InsertJob(ctx context.Context, job *v1.SysJob, opts *api.CreateOptions) error {
	if opts.Validate {
		if err := job.Validate(); err != nil {
			return err
		}
	}

	// 检查是否存在job
	jobinfo, _ := s.store.Jobs().SelectJobById(ctx, job.JobId, &api.GetOptions{})
	if jobinfo != nil {
		return errors.Errorf("job %d exists", job.JobId)
	}

	job.Status = v1.PAUSE
	// 注册job
	if err := s.store.Jobs().InsertJob(ctx, job, opts); err != nil {
		// 注册失败时要删除任务
		nodeId, _ := s.jobMp.Exists(job.JobId)
		node, _ := s.nodePool.GetNode(nodeId)
		resp, err := rpc.Remoting.RemoveJobs(ctx, []int64{job.JobId}, callopt.WithHostPort(node.IPAddress))
		if err != nil || resp.Code != code.SUCCESS {
			return errors.New("delete job error")
		}
	}

	return nil
}

func (s *JobService) UpdateJob(ctx context.Context, job *v1.SysJob, opts *api.UpdateOptions) error {
	if opts.Validate {
		if err := job.Validate(); err != nil {
			return err
		}
	}

	// 先调度任务后更新
	if nodeid, ok := s.jobMp.Exists(job.JobId); ok {
		node, err := s.nodePool.GetNode(nodeid)
		if err != nil {
			log.Errorf("获取节点失败")
			return err
		}

		resp, err := rpc.Remoting.UpdateJob(ctx, v12.SysJob2JobInfo(job), callopt.WithHostPort(node.IPAddress))
		if err != nil {
			return rpc.ErrRPC
		}

		if resp.Code != code.SUCCESS {
			return errors.New(resp.Msg)
		}
	} else {
		err := s.Schedule(ctx, job)
		if err != nil {
			return err
		}
	}

	// 最后更新任务
	return s.store.Jobs().UpdateJob(ctx, job, opts)
}

func (s *JobService) Schedule(ctx context.Context, job *v1.SysJob) error {
	// 没有实例运行时立刻失败
	nodes := s.nodePool.OnlineNodes()
	if len(nodes) == 0 {
		log.Debugf("当前无可用节点，%d 任务调度失败，进入死亡队列", job.JobId)
		return errors.New("当前无可用节点, 调度失败")
	}

	picker := s.lb.GetPicker(s.nodePool)
	strId := strconv.Itoa(int(job.JobId))
	node := picker.Next(ctx, strId)

	// 调度失败次数
	failTime := 0
	for failTime < 3 {
		// 在目标节点上创建任务
		resp, err := rpc.Remoting.CreateJob(ctx, v12.SysJob2JobInfo(job), callopt.WithHostPort(node.Address().String()))

		if err != nil {
			// rpc调用失败，切换目标节点
			node = picker.Next(ctx, strId)
			failTime++
			continue
		}

		if resp.Code != code.SUCCESS {
			return errors.New(resp.Msg)
		}
	}

	// 任务进入死亡队列
	models.DeadJobQueue.Add(job.JobId, job)
	return nil
}

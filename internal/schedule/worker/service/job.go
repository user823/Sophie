package service

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/internal/schedule/models"
	utils "github.com/user823/Sophie/internal/schedule/utils"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
)

type JobSrv interface {
	CreateJob(ctx context.Context, job *v1.SysJob, opts *api.CreateOptions) error
	RemoveJobs(ctx context.Context, jobIds []int64, opts *api.DeleteOptions) error
	PauseJobs(ctx context.Context, jobIds []int64) error
	ResumeJobs(ctx context.Context, jobIds []int64) error
	Run(ctx context.Context, jobIds []int64) error
	UpdateJob(ctx context.Context, job *v1.SysJob, opts *api.UpdateOptions) error
}

type JobService struct {
	jobMp  models.JobMap
	nodeid string
}

func NewJobs(jobMp models.JobMap, nodeid string) JobSrv {
	return &JobService{jobMp, nodeid}
}

func (s *JobService) CreateJob(ctx context.Context, job *v1.SysJob, opts *api.CreateOptions) error {
	if opts.Validate {
		if err := job.Validate(); err != nil {
			return err
		}
	}

	// 检查调用目标
	if err := utils.CheckTarget(job.InvokeTarget); err != nil {
		return err
	}

	// 检查是否存在job
	if jobinfo := models.GetJobInfo(job.JobId); jobinfo != nil {
		return nil
	}

	if err := s.jobMp.Create(job.JobId, s.nodeid); err != nil {
		return err
	}

	// 注册job
	models.RegisterCronTask(job.CronExpression, func() {
		if err := utils.CallTarget(job.InvokeTarget); err != nil {
			log.Errorf("job invoke error: %s", err.Error())
		}
	}, *job)
	return nil
}

func (s *JobService) RemoveJobs(ctx context.Context, jobIds []int64, opts *api.DeleteOptions) error {
	for i := range jobIds {
		models.DeleteJob(jobIds[i])
	}
	return nil
}

func (s *JobService) PauseJobs(ctx context.Context, jobIds []int64) error {
	for i := range jobIds {
		log.Info(jobIds[i])
		models.PauseJob(jobIds[i])
	}
	return nil
}

func (s *JobService) ResumeJobs(ctx context.Context, jobIds []int64) error {
	for i := range jobIds {
		models.ResumeJob(jobIds[i])
	}
	return nil
}

func (s *JobService) Run(ctx context.Context, jobIds []int64) error {
	var invalidIds []int64
	for i := range jobIds {
		jobInfo := models.GetJobInfo(jobIds[i])
		if jobInfo == nil {
			invalidIds = append(invalidIds, jobIds[i])
			continue
		}
		models.RunJob(jobIds[i])
	}
	return errors.Errorf("不存在目标任务: %v", invalidIds)
}

func (s *JobService) UpdateJob(ctx context.Context, job *v1.SysJob, opts *api.UpdateOptions) error {
	if opts.Validate {
		if err := job.Validate(); err != nil {
			return err
		}
	}

	// 检查调用目标
	if err := utils.CheckTarget(job.InvokeTarget); err != nil {
		return err
	}

	// 存在job则删除
	jobInfo := models.GetJobInfo(job.JobId)
	if jobInfo != nil {
		models.DeleteJob(job.JobId)
	}

	// 重新创建job
	models.RegisterCronTask(job.CronExpression, func() {
		if err := utils.CallTarget(job.InvokeTarget); err != nil {
			log.Errorf("job invoke error: %s", err.Error())
		}
	}, *job)

	return nil
}

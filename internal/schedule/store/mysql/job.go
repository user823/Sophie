package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/internal/schedule/utils"
	"github.com/user823/Sophie/pkg/errors"
	"gorm.io/gorm"
)

type mysqlJobStore struct {
	db *gorm.DB
}

func selectJobVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_job")
}

func (s *mysqlJobStore) SelectJobList(ctx context.Context, job *v1.SysJob, opts *api.GetOptions) ([]*v1.SysJob, int64, error) {
	query := selectJobVo(s.db)
	if job.JobName != "" {
		query = query.Where("job_name like ?", "%"+job.JobName+"%")
	}
	if job.JobGroup != "" {
		query = query.Where("job_group = ?", job.JobGroup)
	}
	if job.Status != "" {
		query = query.Where("status = ?", job.Status)
	}
	if job.InvokeTarget != "" {
		query = query.Where("invoke_taget like ?", "%"+job.InvokeTarget+"%")
	}
	query = opts.SQLCondition(query, "")
	var result []*v1.SysJob
	err := query.Find(&result).Error

	return result, utils.CountQuery(query, opts, ""), err
}

func (s *mysqlJobStore) SelectJobAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysJob, int64, error) {
	return s.SelectJobList(ctx, &v1.SysJob{}, opts)
}

func (s *mysqlJobStore) SelectJobByIds(ctx context.Context, ids []int64, opts *api.GetOptions) ([]*v1.SysJob, int64, error) {
	query := selectJobVo(s.db).Where("job_id in ?", ids)
	var result []*v1.SysJob
	err := query.Find(&result).Error
	if err != nil {
		return []*v1.SysJob{}, 0, err
	}
	return result, int64(len(result)), nil
}

func (s *mysqlJobStore) SelectJobById(ctx context.Context, jobId int64, opts *api.GetOptions) (*v1.SysJob, error) {
	query := selectJobVo(s.db).Where("job_id = ?", jobId)
	var result *v1.SysJob
	err := query.First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *mysqlJobStore) DeleteJobByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("job_id in ?", ids).Delete(&v1.SysJob{})
	return del.Error
}

func (s *mysqlJobStore) UpdateJob(ctx context.Context, job *v1.SysJob, opts *api.UpdateOptions) error {
	if job.JobId == 0 {
		return errors.New("更新job 必须指明id")
	}
	update := opts.SQLCondition(s.db).Model(&v1.SysJob{}).Where("job_id = ?", job.JobId).Updates(job)
	return update.Error
}

func (s *mysqlJobStore) InsertJob(ctx context.Context, job *v1.SysJob, opts *api.CreateOptions) error {
	return opts.SQLCondition(s.db).Create(job).Error
}

package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/internal/schedule/utils"
	"gorm.io/gorm"
)

type mysqlJobLogStore struct {
	db *gorm.DB
}

func selectJobLogVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_job_log")
}

func (s *mysqlJobLogStore) SelectJobLogList(ctx context.Context, jobLog *v1.SysJobLog, opts *api.GetOptions) ([]*v1.SysJobLog, int64, error) {
	query := selectJobVo(s.db)
	if jobLog.JobName != "" {
		query = query.Where(query, "jobName like ?", "%"+jobLog.JobName+"%")
	}
	if jobLog.JobGroup != "" {
		query = query.Where(query, "job_group = ?", jobLog.JobGroup)
	}
	if jobLog.Status != "" {
		query = query.Where(query, "status = ?", jobLog.Status)
	}
	if jobLog.InvokeTarget != "" {
		query = query.Where(query, "invoke_target like ?", "%"+jobLog.InvokeTarget+"%")
	}
	query = opts.SQLCondition(query, "create_time")

	var result []*v1.SysJobLog
	err := query.Find(&result).Error
	return result, utils.CountQuery(query, opts, "create_time"), err
}

func (s *mysqlJobLogStore) SelectJobLogAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysJobLog, int64, error) {
	return s.SelectJobLogList(ctx, &v1.SysJobLog{}, opts)
}

func (s *mysqlJobLogStore) SelectJobLogById(ctx context.Context, jobLogId int64, opts *api.GetOptions) (*v1.SysJobLog, error) {
	query := selectJobLogVo(s.db).Where("job_log_id = ?", jobLogId)
	var result *v1.SysJobLog
	if err := query.First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *mysqlJobLogStore) InsertJobLog(ctx context.Context, jobLog *v1.SysJobLog, opts *api.CreateOptions) error {
	return opts.SQLCondition(s.db).Create(jobLog).Error
}

func (s *mysqlJobLogStore) DeleteJobLogByIds(ctx context.Context, logIds []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("job_log_id in ?", logIds).Delete(&v1.SysJobLog{})
	return del.Error
}

func (s *mysqlJobLogStore) DeleteJobLogById(ctx context.Context, jobId int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("job_log_id = ?", jobId).Delete(&v1.SysJobLog{})
	return del.Error
}

func (s *mysqlJobLogStore) CleanJobLog(ctx context.Context, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("true").Delete(&v1.SysJobLog{})
	return del.Error
}

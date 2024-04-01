namespace go v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/schedule/v1 -service WorkerService -gen-path api/thrift/schedule  worker.thrift

include "job.thrift"

service WorkerService {
    job.BaseResp CreateJob(1:job.JobInfo job)
    job.BaseResp RemoveJobs(1:list<i64> jobIds)
    job.BaseResp PauseJobs(1:list<i64> jobIds)
    job.BaseResp ResumeJobs(1:list<i64> jobIds)
    job.BaseResp Run(1:list<i64> jobIds)
    job.BaseResp UpdateJob(1:job.JobInfo job)
}
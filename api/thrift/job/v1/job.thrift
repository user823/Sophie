namespace go v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/job/v1 -service JobService -gen-path api/thrift/job  job.thrift

service JobService {
    // job 模块
    ListJobsResponse ListJobs(1:ListJobsRequest req)
    ExportJobsResponse ExportJobs(1:ExportJobsRequest req)
    GetJobInfoResponse GetJobInfo(1:i64 jobId)
    BaseResp CreateJob(1:CreateJobRequest req)
    BaseResp UpdateJob(1:UpdateJobRequest req)
    BaseResp ChangeStatus(1:ChangeStatusRequest req)
    BaseResp Run(1:RunRequest req)
    BaseResp RemoveJobs(1:RemoveJobsRequest req)

    // job log 模块
    ListJobLogsResponse ListJobLogs(1:ListJobLogsRequest req)
    ExportJobLogsResponse ExportJobLogs(1:ExportJobLogsRequest req)
    GetJobLogInfoResponse GetJobLogInfo(1:i64 jobLogId)
    BaseResp RemoveJobLogs(1:RemoveJobLogsRequest req)
    BaseResp Clean()
}

struct BaseResp {
    1:i64 code
    2:string msg
}

struct PageInfo {
    1:i64 pageNum
    2:i64 pageSize
    3:string orderByColumn
    4:string isAsc
}

struct DateRange {
    1:i64 beginTime // 使用毫秒
    2:i64 endTime
}

struct JobInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 jobId
    8:string jobName
    9:string jobGroup
    10:string invokeTarget
    11:string cronExpression
    12:string misfirePolicy
    13:string concurrent
    14:string status
}

struct JobLog {
    1:i64 jobLogId
    2:string jobName
    3:string jobGroup
    4:string invokeTarget
    5:string jobMessage
    6:string status
    7:string exceptionInfo
    8:i64 startTime
    9:i64 stopTime
}

struct ListJobsRequest {
    1:PageInfo pageInfo
    2:JobInfo jobInfo
}

struct ListJobsResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<JobInfo> rows
}

struct ExportJobsRequest {
    1:PageInfo pageInfo
    2:JobInfo jobInfo
}

struct ExportJobsResponse {
    1:BaseResp baseResp
    2:list<JobInfo> list
    3:string sheetName
    4:string title
}

struct GetJobInfoResponse {
    1:BaseResp baseResp
    2:JobInfo data
}

struct CreateJobRequest {
    1:JobInfo jobinfo
}

struct UpdateJobRequest {
    1:JobInfo jobInfo
}

struct ChangeStatusRequest {
    1:JobInfo jobInfo
}

struct RunRequest {
    1:JobInfo jobInfo
}

struct RemoveJobsRequest {
    1:list<i64> jobIds
}

struct ListJobLogsRequest {
    1:PageInfo pageInfo
    2:JobLog jobLog
}

struct ListJobLogsResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<JobLog> rows
}

struct ExportJobLogsRequest {
    1:JobLog jobLog
}

struct ExportJobLogsResponse {
    1:BaseResp baseResp
    2:list<JobLog> list
    3:string sheetName
    4:string title
}

struct GetJobLogInfoResponse {
    1:BaseResp baseResp
    2:JobLog data
}

struct RemoveJobLogsRequest {
    1:BaseResp baseResp
    2:list<i64> jobIds
}
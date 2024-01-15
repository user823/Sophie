namespace go v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/job/v1 -service JobService -gen-path api/thrift/job  job.thrift

service JobService {

    // job controll
    ListJobsResponse ListJobs(1:ListJobsRequest req)
    ExportJobResponse Export(1:ExportJobRequest req)
    JobInfoResponse GetJobInfoById(1:i64 id)
    CreateJobResponse CreateJob(1:CreateJobRequest req)
    UpdateJobResponse UpdateJob(1:UpdateJobRequest req)
    ChangeJobStatusResponse ChangeJobStatus(1:ChangeJobStatusRequest req)
    RunJobResponse Run(1:RunJobRequest req)
    DeleteJobResponse DeleteJob(1:DeleteJobRequest req)

    // job log controll
    ListJobLogsResponse ListJobLogs(1:ListJobLogsRequest req)
    ExportJobLogResponse ExportJobLog(1:ExportJobLogRequest req)
    JobLogInfoResponse GetJobLogInfoById(1:i64 id)
    DeleteJobLogResponse DeleteJobLog(1:DeleteJobLogRequest req)
    CleanResponse Clean()
}

struct BaseResp {
    1:i64 code 
    2:string message
}
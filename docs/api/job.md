# 定时任务接口

## 定时任务
GET v1/job/list
GET v1/job/export
GET v1/job/{jobId}
POST v1/job
PUT v1/job
PUT v1/job/changeStatus
PUT v1/job/run
DELETE v1/job/{jobIds}

## 定时任务日志
GET v1/job/log/list
POST v1/job/log/export
GET v1/job/log/{jobLogId}
DELETE v1/job/log/{jobLogIds}
DELETE v1/job/log/clean

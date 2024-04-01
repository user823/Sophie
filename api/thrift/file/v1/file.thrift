namespace go v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/file/v1 -service FileService -gen-path api/thrift/file  file.thrift

service FileService {
    UploadResponse Upload(1:UploadRequest req)
}

struct BaseResp {
    1:i64 code
    2:string msg
}

struct FileInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:string name
    8:string url
}

struct UploadRequest {
    1:binary data
    2:string path
    3:i64 userId
}

struct UploadResponse {
    1:BaseResp baseResp
    2:FileInfo file
}
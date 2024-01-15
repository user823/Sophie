namespace go v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/file/v1 -service FileService -gen-path api/thrift/file  file.thrift

service FileService {
    FileInfoResponse Upload(1:UploadRequest req)
}

struct BaseResp {
    1:i64 code
    2:string message
    3:i64 service_time
}

struct FileInfo {
    1:string name
    2:string url 
}

struct UploadRequest {
    1:binary file_data
}

struct FileInfoResponse {
    1: BaseResp base_resp
    2: FileInfo file_info
}
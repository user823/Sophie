namespace go v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/gen/v1 -service GenService -gen-path api/thrift/gen  gen.thrift

service GenService {
    ListGenTablesResponse ListGenTables(1:ListGenTablesRequest req)
    GetInfoResponse GetInfoById(1:i64 id) // 用于修改代码生成业务
    ListGenTablesResponse DataList(1:DataListRequest req)
    ListGenTablesResponse ListColumnsById(1:i64 id)
    ImportTableSaveResponse ImportTableSave(1:string tables) // 保存表结构
    EditSaveResponse EditSave(1:EditSaveRequest req)
    DeleteTablesResponse DeleteTables(1:DeleteTableRequest req)
    PreviewResponse PreviewById(1:i64 id)
    DownloadResponse DownloadByName(1:string name)
    GenCodeResponse GenCode(1:string name) 
    SynchDBResponse SynchDB(1:string name)
    BatchGenCodeResponse BatchGenCode(1:string tables)
}

struct BaseResp {
    1:i64 code 
    2:string message
}
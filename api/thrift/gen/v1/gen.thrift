namespace go v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/gen/v1 -service GenService -gen-path api/thrift/gen  gen.thrift

service GenService {
    ListGenTablesResponse ListGenTables(1:ListGenTablesRequest req)
    GetInfoResponse GetInfo(1:i64 tableId)
    ListGenTablesResponse DataList(1:DataListRequest req)
    ListGenTablesResponse ColumnList(1:i64 tableId)
    BaseResp ImportTableSave(1:string tables)
    BaseResp EditSave(1:EditSaveRequest req)
    BaseResp Remove(1:RemoveRequest req)
    PreviewResponse Preview(1:i64 tableId)
    DownloadResponse Download(1:string tableName)
    BaseResp GenCode(1:string tableName)
    BaseResp SynchDb(1:string tableName)
    BatchGenCodeResponse BatchGenCode(1:string tables)
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

struct BaseInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
}

struct GenTable {
    1:BaseInfo baseInfo
    2:i64 tableId
    3:optional string tableName
    4:optional string tableComment
    5:optional string subTableName
    6:optional string subTableFkName
    7:optional string className
    8:optional string tplCategory
    9:optional string tplWebType
    10:optional string packageName
    11:optional string moduleName
    12:optional string businessName
    13:optional string functionName
    14:optional string functionAuthor
    15:optional string genType
    16:optional string genPath
    17:optional GenTableColumn pkColumn
    18:optional GenTable subTable
    19:optional list<GenTableColumn> columns
    20:optional string options
    21:optional string treeCode
    22:optional string treeParentCode
    23:optional string treeName
    24:optional string parentMenuId
    25:optional string parentMenuName
}

struct GenTableColumn {
    1:BaseInfo baseInfo
    2:i64 columnId
    3:optional i64 tableId
    4:optional string columnName
    5:optional string columnComment
    6:optional string columnType
    7:optional string GoType
    8:optional string GoField
    9:optional string isPk
    10:optional string isIncrement
    11:optional string isRequired
    12:optional string isInsert
    13:optional string isEdit
    14:optional string isList
    15:optional string isQuery
    16:optional string queryType
    17:optional string htmlType
    18:optional string dictType
    19:optional i64 sort
}

struct ListGenTablesRequest {
    1:PageInfo pageInfo
    2:DateRange dateRange
    3:GenTable genTable
}

struct ListGenTablesResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<GenTable> rows
}

struct GetInfoResponse {
    1:BaseResp baseResp
    2:GenTable info
    3:list<GenTableColumn> rows
    4:list<GenTable> tables
}

struct DataListRequest {
    1:GenTable genTable
}

struct EditSaveRequest {
    1:GenTable genTable
}

struct RemoveRequest {
    1:list<i64> tableIds
}

struct PreviewResponse {
    1:BaseResp baseResp
    2:map<string, string> data
}

struct DownloadResponse {
    1:BaseResp baseResp
    2:binary data
}

struct BatchGenCodeResponse {
    1:BaseResp baseResp
    2:list<binary> data
}
namespace go v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/gen/v1 -service GenService -gen-path api/thrift/gen  gen.thrift

service GenService {
    ListGenTablesResponse ListGenTables(1:ListGenTablesRequest req)
    GetInfoResponse GetInfo(1:i64 tableId)
    ListGenTablesResponse DataList(1:DataListRequest req)
    ListGenTableColumnsResponse ColumnList(1:i64 tableId)
    BaseResp ImportTableSave(1:ImportTableSaveRequest req)
    BaseResp EditSave(1:EditSaveRequest req)
    BaseResp Remove(1:RemoveRequest req)
    PreviewResponse Preview(1:i64 tableId)
    DownloadResponse Download(1:string tableName)
    BaseResp GenCode(1:string tableName)
    BaseResp SynchDb(1:string tableName)
    DownloadResponse BatchGenCode(1:string tables)
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

struct GenTable {
    1:string createBy
    2:string createTime
    3:string updateBy
    4:string updateTime
    5:string remark
    6:map<string,string> params
    7:i64 tableId
    8:string tableName
    9:string tableComment
    10:string subTableName
    11:string subTableFkName
    12:string className
    13:string tplCategory
    14:string tplWebType
    15:string packageName
    16:string moduleName
    17:string businessName
    18:string functionName
    19:string functionAuthor
    20:string genType
    21:string genPath
    22:GenTableColumn pkColumn
    23:GenTable subTable
    24:list<GenTableColumn> columns
    25:string options
    26:string treeCode
    27:string treeParentCode
    28:string treeName
    29:string parentMenuId
    30:string parentMenuName
}

struct GenTableColumn {
    1:string createBy
    2:string createTime
    3:string updateBy
    4:string updateTime
    5:string remark
    6:map<string,string> params
    7:i64 columnId
    8:i64 tableId
    9:string columnName
    10:string columnComment
    11:string columnType
    12:string goType
    13:string goField
    14:string isPk
    15:string isIncrement
    16:string isRequired
    17:string isInsert
    18:string isEdit
    19:string isList
    20:string isQuery
    21:string queryType
    22:string htmlType
    23:string dictType
    24:i64 sort
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

struct ImportTableSaveRequest {
    1:string tables
    2:string operName
}

struct GetInfoResponse {
    1:BaseResp baseResp
    2:GenTable info
    3:list<GenTableColumn> rows
    4:list<GenTable> tables
}

struct DataListRequest {
    1:GenTable genTable
    2:PageInfo pageInfo
    3:DateRange dateRange
}

struct ListGenTableColumnsResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<GenTableColumn> rows
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

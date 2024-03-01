#!/bin/bash

# 建立对应于mysql的elasticsearch 索引

endpoint='https://localhost:9200'
cacert_file="${HOME}/.cert/http_ca.crt"

# sys_oper_log
sys_oper_log=$(cat << EOF
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1
  },
  "mappings": {
    "properties": {
      "operId": {"type": "long"},
      "title": {"type": "keyword"},
      "businessType": {"type": "long"},
      "method": {"type": "keyword"},
      "requestMethod": {"type": "keyword"},
      "operatorType": {"type": "long"},
      "operName": {"type": "keyword"},
      "deptName": {"type": "keyword"},
      "operUrl": {"type": "keyword"},
      "operIp": {"type": "keyword"},
      "operParam": {"type": "keyword"},
      "jsonResult": {"type": "keyword"},
      "status": {"type": "keyword"},
      "errorMsg": {"type": "keyword"},
      "operTime": {"type": "long"},
      "costTime": {"type": "long"}
    }
  }
}
EOF
)

sys_logininfor=$(cat << EOF
{
  "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 1
    },
    "mappings": {
      "properties": {
        "infoId": {"type": "long"},
        "userName": {"type": "keyword"},
        "status": {"type": "keyword"},
        "ipaddr": {"type": "keyword"},
        "msg": {"type": "keyword"},
        "accessTime": {"type": "date"}
      }
    }
}
EOF
)

curl --cacert "${cacert_file}" -u sophie:123456 -XPUT "${endpoint}/sys_oper_log" -H "Content-Type: application/json" -d "$sys_oper_log"
curl --cacert "${cacert_file}" -u sophie:123456 -XPUT "${endpoint}/sys_logininfor" -H "Content-Type: application/json" -d "$sys_sys_logininfor"

# sys_user
#sys_user=$(cat << EOF
#{
#  "settings": {
#    "number_of_shards": 1,
#    "number_of_replicas": 1,
#    "analysis": {
#      "char_filter": {
#        "comma_filter": {
#          "type": "pattern_replace",
#          "pattern": ",",
#          "replacement": " "
#        }
#      },
#      "analyzer": {
#        "comma_analyzer": {
#          "type": "custom",
#          "char_filter": ["comma_filter"],
#          "tokenizer": "standard"
#        }
#      }
#    }
#  },
#  "mappings": {
#    "properties": {
#      "user_id": {"type": "long"},
#      "dept_id": {"type": "long"},
#      "user_name": {"type": "keyword"},
#      "nick_name": {"type": "keyword"},
#      "email": {"type": "keyword"},
#      "phonenumber": {"type": "keyword"},
#      "sex": {"type": "keyword"},
#      "avatar": {"type": "keyword"},
#      "password": {"type": "keyword"},
#      "status": {"type": "keyword"},
#      "del_flag": {"type": "keyword"},
#      "login_ip": {"type": "keyword"},
#      "login_date": {"type": "date"},
#      "create_by": {"type": "keyword"},
#      "create_time": {"type": "date"},
#      "update_by": {"type": "keyword"},
#      "update_time": {"type": "date"},
#      "remark": {"type": "keyword"},
#      "parent_id": {"type": "long"},
#      "dept_name": {"type": "keyword"},
#      "ancestors": {"type": "text", "analyzer": "comma_analyzer"},
#      "order_num": {"type": "long"},
#      "leader": {"type": "keyword"},
#      "dept_status": {"type": "keyword"},
#      "role_id": {"type": "long"},
#      "role_name": {"type": "keyword"},
#      "role_key": {"type": "keyword"},
#      "role_sort": {"type": "keyword"},
#      "data_scope": {"type": "keyword"},
#      "role_status": {"type": "keyword"}
#    }
#  }
#}
#EOF
#)



# 数据库

## table
1. table_data 数据表,存储每张表当前的版本号及真实数据
```
table_name  varchar 255   // 表名
now_version int8    64    // 版本号
data        varchar 65535 // 数据
```

2. name_cmd 指令表
```
version   int4    32     //表版本号，自增
users     varchar 128    //当前版本提交操作的用户
date      varchar 64     //版本生成时间，2019-11-14 10:10:10
cmds      varchar 65535  //操作指令
```

## 表是否存在的判断
用户删除某一张表，则不用显示给前端，但任然保留表数据，退档使用

## tag
```
tagName       tag,string
字段名         字段名为表名，值为版本号，int，只增不减
或者拼接字符串  "tableName：version"  
```


## v1
### 注意事项

1.默认冻结前三行，手动冻结指定列。

2.行列的拖拽，仅在前端展示，不影响实际行列数据。

3.当涉及多人修改同一单元格时，最后的编辑将覆盖之前的编辑数据。

4.文件保存，任意一人点击"保存"将保存从上版本到当前所有人操作指令，生成一条版本记录。数据有修改但无人保存，将在所有客户端关闭时自动保存。

5.在浏览历史版本时，不能修改数据。只有处于编辑状态下才能修改数据。

6.回退版本时，将丢弃当前所有人的编辑指令、数据，实际数据回退到对应版本，强行同步给所有人。

### 功能支持

1.实时显示他人当前编辑的单元格。以不同颜色区分。

2.查看历史版本，仅浏览数据，实际数据不回退，可导出该浏览版本数据。

3.版本回退，可回退到之前任意版本，实际数据将回退到该版本。

4.导出文件，以当前前端展示的数据导出文件，格式支持 csv、xlsx。


## 数据库
### table
1. table_data 数据表,存储每张表当前的版本号及真实数据
```
table_name  varchar 255   // 表名
describe    varchar 255   // 描述
version     int8          // 版本号
date        varchar 64    // 日期
data        varchar 65535 // 数据
```

2. name_cmd 指令表
```
version   int4    32     //表版本号，自增
users     varchar 128    //当前版本提交操作的用户
date      varchar 64     //版本生成时间，2019-11-14 10:10:10
cmds      varchar 65535  //操作指令
```


# 数据库

## table
功能：读取当前表数据，修改数据并添加版本号，回退版本读取历史版本数据

table 只增不删，防止回档出错

1。
表名为文件名，每一行为表的版本数据
数据量大
回退版本直接读取。
```
version   表版本号，int类型自增，
userName  当前版本提交用户名，string
data      表数据，blob
```

2。
以变量的形式来存储当前版本相较于上一版本变化的数据
数据量小
回退版本，反向操作
```
version   表版本号，int类型自增，
userName  当前版本提交用户名，string
cmds      操作指令，blob
```

3。
table_data // 存储每张表当前的版本号及数据
```
table_name // 表名
now_version // 版本号
data   // 数据
```


## 表是否存在的判断
用户删除某一张表，则不用显示给前端，但任然保留表数据，退档使用

## tag
```
tagName       tag,string
字段名         字段名为表名，值为版本号，int，只增不减
或者拼接字符串  "tableName：version"  
```

# 支持操作

## 创建表 http
req  createTable?tableName=xxx&userName=xxx
resp ok= 1/0
     msg= string

## 删除表 http
req  deleteTable?tableName=xxx&userName=xxx
resp ok= 1/0
     msg= string

## 表列表 http
req  getAllTable
resp ok= 1/0
     msg= string
     tables = []{tableName:string,version:int}
        
## 下载表 http
req  downloadTable?tableName=xxx
resp ok= 1/0
     msg= string
     file
        
websocket

## pushCellData 推送变化的单元格数据
```
cmd : "pushCellData"
cellDate:  // []map[string]string 
              []{col:int,row:int,oldValue:string,newValue:string,userName:string}
```

## pushCellSelected
```
cmd : "pushCellSelected"
selected:  // []map[string]string 
              []{col:int,row:int,userName:string}
```

## pushSaveTable 保存表后推送
```
cmd : "pushSaveTable"
version : int
data:  // [][]string
```

## pushAll 推送所有数据
```
cmd : "pushAll"
tableName : "xxx"
version : int
data:  // [][]string
```

## pushError 返回错误信息
```
cmd : "pushError"
doCmd:  string
errMsg: string            
```

## openTable 打开文件
tos
```
cmd : "openTable"
tableName : "xxx.xlsx"
userName: "deng"
```

toc
```
cmd : "openTable"
tableName : "xxx.xlsx"
userName: "deng"
version : int
data: [][]string
```

## 插入行, tos
```
cmd : "insertRow"
index: 2  , int
```

## 删除行, tos
```
cmd : "removeRow"
index: 2  , int
```

## 插入列, tos
```
cmd : "insertCol"
index: 2  , int
```

## 删除列, tos
```
cmd : "removeCol"
index: 2  , int
```

## 修改数据, tos
```
cmd : "setCellValues"
cellValues: []map[string]string
            []{col:int,row:int,oldValue:string,newValue:string}
```

## 选中格子, tos
```
cmd : "cellSelected"
selected: []map[string]string
          []{col:int,row:int}
```

## 保存表，tos
```
cmd : "saveTable"
```


## 查看历史版
tos
```
cmd    : "lookHistory"
version: int
```

toc
```
cmd    : "lookHistory"
version: int
data   :  //[][]string
```

## 回退到历史版本
tos
```
cmd  : "rollback"
goto : int // 回退到版本号
```

toc
```
cmd    : "rollback"
version: int
data   :  //[][]string
```

 
 # 冲突
 
锁定格子，当A用户编辑了一个格子，必须保存后，B用户才能再编辑。
若A用户将A1格子由0写为1，B用户将A1格子1改为2，如果A用户不保存后退出，编辑回退，B用户看到的为0。


## 
1.不能有空行，必须的填有字段
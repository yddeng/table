# http
## 创建表 
req  createTable?tableName=xxx&userName=xxx
resp ok= 1/0
     msg= string

## 删除表 
req  deleteTable?tableName=xxx&userName=xxx
resp ok= 1/0
     msg= string

## 表列表 
req  getAllTable
resp ok= 1/0
     msg= string
     tables = []{tableName:string,version:int}
        
#websocket

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
doCmd:  string   // 执行的指令
errMsg: string   // 错误信息
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

## 插入行, tos,toc
```
cmd : "insertRow"
index: 2  , int
```
## 删除行, tos,toc
```
cmd : "removeRow"
index: 2  , int
```
## 插入列, tos,toc
```
cmd : "insertCol"
index: 2  , int
```
## 删除列, tos ,toc
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
selected: {row:int, col:int, row2:int, col2:int},起始格到结束格 
```
toc
```
cmd : "cellSelected"
userName: string,
selected: {row:int, col:int, row2:int, col2:int},起始格到结束格 
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

## 回到编辑状态
tos
```
cmd : "backEditor"
```

toc
```
cmd    : "backEditor"
version: int
data   :  //[][]string
```

## 回退到历史版本
tos
```
cmd  : "rollback"
version : int      // 回退到版本号
```

toc
```
cmd    : "rollback"
version: int
data   :  //[][]string
```
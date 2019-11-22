# http
## 创建表 post
req  createTable {tableName:string,describe:string,token:string}
resp ok= 1/0
     msg= string

## 删除表 post
req  deleteTable {tableName:string,token:string}
resp ok= 1/0
     msg= string

## 表列表 get
req  getAllTable
resp ok= 1/0
     msg= string
     tables = []{tableName:string,version:int}

## 下载 post
req  downloadTable  {tableName:string,token:string}
resp ok= 1/0
     msg= string
     data = [][]string
     
## 登陆 post
req login {userName:string,password:string}
resp ok  = 1/0
     msg = string
     token = string
        
#websocket
## pushAll 推送所有数据
```
cmd : "pushAll"
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

## 保存表，
tos
```
cmd : "saveTable"
```
toc
```
cmd : "saveTable"
version : int
data:  // [][]string
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

## 版本列表
tos
```
cmd: "versionList"
```

toc
```
cmd    : "versionList"
list: []{version:int,users:[]string,date:string} // 版本历史记录
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

数据库存入cmd
```
cmd    : "rollback"
now    : int  // 当前版本
goto   : int  // 回退到的版本
```
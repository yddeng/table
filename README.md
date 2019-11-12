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
change    变量，blob
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

## 创建表
```
cmd : "createTable"
fileName : "xxx"
```

## 删除表
```
cmd : "deleteTable"
fileName : "xxx"
```

## pushData 推送最新的表数据, toc
```
cmd : "pushData"
data:  data  //[][]string
cellLocked:  // []map[string]string 
                []{col:int,row:int,userName:string}
```

## openFile 打开文件, tos
```
cmd : "openFile"
fileName : "xxx.xlsx"
userName: "deng"
```

## 插入行, tos
```
cmd : "insertRow"
rowIndex: 2  , int
```

## 删除行, tos
```
cmd : "deleteRow"
rowIndex: 2  , int
```

## 插入列, tos
```
cmd : "insertCol"
colIndex: 2  , int
```

## 删除列, tos
```
cmd : "deleteCol"
colIndex: 2  , int
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

## 保存操作，tos
```
cmd : "saveFile"
```
 
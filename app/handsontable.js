//import zhCh from "./dist/zh-CN.min"

var handstable = {};

var StatusEnum = {
    NONE: 1,    // 没有初始化的操作
    EDITOR: 2,  // 编辑状态
    LOOK: 3,    // 查看状态，不能操作
};

handstable.status = StatusEnum.NONE;
handstable.container = null;  // html
handstable.version = null;    // 版本号
handstable.data = null;       // 表数据
handstable.tableName = null;  // 表名
handstable.table = null;      // handsontable 实列

handstable.selected = {};     // 用户选择单元格本地存储
handstable.userColor = {};    //用户颜色

handstable.customBordersPlugin = null;  // 单元格组件
handstable.exportPlugin =null;          // 导出组件

handstable.init = function (container,name) {
    this.container = container;
    this.tableName = name;
};

handstable.new = function(data){
    if (data.length ==0){
        data = [[]]
    }
    handstable.setStatue(StatusEnum.EDITOR);
    handstable.data = data;
    handstable.table =  new Handsontable(handstable.container, {
        data: data,
        width:'100%',
        height:window.innerHeight - 100,
        minCols:26, //最小列数
        //maxCols: 50, //最大列数
        minRows:50 , //最小行数
        //maxRows: 40 , //最大行数
        minSpareCols: 10, //添加空列
        minSpareRows :10,//添加空行
        colHeaders:true, //显示列表头， 默认false， 可取 true/fals/数组 ，当值为数组时，列头为数组的值
        rowHeaders:true,　//显示行表头， 默认 false， 可取 true/fals/数组，当值为数组时，行头为数组的值
        mergeCells: false, //表示允许单元格合并
        fixedRowsTop : 3,    //冻结行（固定顶部开始算起指定行不随垂直滚动条滚动；
        fixedColumnsLeft : 0, //冻结列（固定右边开始算起指定行不随水平滚动条滚动）；
        manualColumnFreeze:true,  //设置true后在单元格中右击出现一个菜单，此菜单会多出一个“Unfteeze this columu”的选项，再次点击就会取消冻结动作。 默认为false
        manualColumnResize: false, //允许拖曳列表头，默认为false
        manualRowResize: false, //允许拖曳行表头，默认为false
        manualColumnMove: false, //true/false 当值为true时，列可拖拽移动到指定列
        manualRowMove: false, //true/false当值为true时，行可拖拽至指定行
        //currentRowClassName:"curRow", //给当前行设置样式名（鼠标点击某个单元格，则整行给样式）
        //currentColClassName:"cur", //给当前行列设置样式名（鼠标点击某个单元格，则整行给样式）
        autoColumnSize: true, //当值为true且列宽未设置时，自适应列大小
        //columnWidth:108,
        //columnSorting:true, //通过点击列表头进行排序（没有图标）
        //colWidths:[10,5,50],
        //rowHeights: [50, 40, 100],
        stretchH:"all", // 可取 last/all/none last：延伸最后一列，all：延伸所有列，none默认不延伸。
        //undo: true,
        //redo: true,
        //language:Handsontable.language.registerLanguageDictionary(zhCh),
        contextMenu:true, //显示表头的下拉菜单默认false 可取 true/false/自定义数组
        contextMenu: {    //内容部分的menu 对功能汉化
            callback: function (key, selection, clickEvent) {
                //console.log(key, selection, clickEvent);
            },
            items: {
                'freeze_column':{
                  name: '冻结该列'
                },
                'unfreeze_column':{
                  name:'取消冻结'
                },
                'freeze_row':{
                    name: '冻结该hang',
                    callback: function () {
                        insertRow(idx)
                    }
                },
              'row_above': {
                name: '上面插入一行',
                  callback: function () {
                    let idx = this.getSelectedLast()[0];
                    insertRow(idx)
                  }
              },
              'row_below': {
                name: '下面插入一行',
                  callback: function () {
                      let idx = this.getSelectedLast()[0]+1;
                      insertRow(idx)
                  }
              },
              'col_left':{
                name: '左侧插入一列',
                  callback: function () {
                      let idx = this.getSelectedLast()[1];
                      insertCol(idx)
                  }
              },
              'col_right':{
                name: '右侧插入一列',
                  callback: function () {
                      let idx = this.getSelectedLast()[1]+1;
                      insertCol(idx)
                  }
              },
                '---------':{},
              'remove_row':{
                name: '移除本行',
                  callback: function () {
                      let idx = this.getSelectedLast()[0];
                      removeRow(idx)
                  }
              },
              'remove_col':{
                name: '移除本列',
                  callback: function () {
                      let idx = this.getSelectedLast()[1];
                      removeCol(idx)
                  }
              },
               /* "colors": { // Own custom option
                    name: 'Colors...',
                    submenu: {
                        // Custom option with submenu of items
                        items: [
                            {
                                // Key must be in the form "parent_key:child_key"
                                key: 'colors:red',
                                name: 'Red',
                                callback: function(key, selection, clickEvent) {
                                    setTimeout(function() {
                                        alert('You clicked red!');
                                    }, 0);
                                }
                            },
                            { key: 'colors:green', name: 'Green' },
                            { key: 'colors:blue', name: 'Blue' }
                        ]
                    }
                },*/
              'separator': Handsontable.plugins.ContextMenu.SEPARATOR,
              'clear_custom': {
                name: '清空所有单元格数据',
                callback: function () {
                  this.clear()
                }
              }
            }
        },
        afterChange:function(changes,source){
            if (changes){
                var cellValues = new Array();
                for (var i = 0;i < changes.length;i++){
                    let item = changes[i];
                    //console.log("change",item)
                    if (item[2] != null && item[2] != ""){
                        cellValues.push( {col:item[1],row:item[0],oldValue:item[2],newValue:item[3]})
                    }else {
                        if (item[3] != null && item[3] != ""){
                            cellValues.push( {col:item[1],row:item[0],oldValue:item[2],newValue:item[3]})
                        }
                    }
                }
                if (cellValues.length >0 ){
                    setCellValues(cellValues)
                }
            }
        },
        afterSelectionEnd: function(row, col, row2, col2) {
            //console.log("afterSelectionEnd",row,col,row2,col2)
            //cellSelected({row:row, col:col, row2:row2, col2:col2})
            handstable.setSelected("名字",{row:row, col:col, row2:row2, col2:col2})
            //handstable.setCellData({col:1,row:1,oldValue:"",newValue:"dd"})
        },
    });
    handstable.customBordersPlugin = handstable.table.getPlugin('customBorders');
    handstable.exportPlugin = handstable.table.getPlugin('exportFile')
};

handstable.setData = function(data){
    handstable.data = data;
    if (handstable.table) {
        handstable.table.loadData(data)
    }
};

handstable.setVersion = function(v){
    handstable.version = v;
    let ver = document.getElementById('n_version');
    ver.innerText = handstable.version
};

handstable.setStatue = function(status){
    handstable.status = status;
    let s = document.getElementById('status');
    s.innerText = util.getStatusString(status)
};

handstable.setSelected = function(name,selected){
    // 将上一次选中清空
    let lastCells = handstable.selected[name];
    if (lastCells){
        handstable.customBordersPlugin.clearBorders(lastCells);
        delete handstable.selected[name];
    }

    let sty = handstable.userColor[name];
    if (!sty){
        let len = Object.keys(handstable.userColor).length;
        sty = util.getColor(len);
        handstable.userColor[name] = sty
    }

    let newCells = [[selected.row,selected.col,selected.row2,selected.col2]];
    //let range = {from: {row: selected.row, col:selected.col},to: {row: selected.row2, col:selected.col2}};
    console.log(handstable.table.getSelectedRange())
    console.log(newCells)
    handstable.customBordersPlugin.setBorders(newCells,sty);
    handstable.selected[name] = newCells;

    let metaRow = handstable.table.getCellMetaAtRow(selected.row)
    console.log(metaRow)
    handstable.table.setCellMeta(selected.row,selected.col,"width","300px")
    let meta = handstable.table.getCellMeta(selected.row,selected.col)
    console.log(meta)
    handstable.table.setDataAtCell(selected.row,selected.col,"sss")
};

handstable.insertRow = function(row){
    handstable.table.alter("insert_row",row)
};
handstable.removeRow = function(row){
    handstable.table.alter("remove_row",row)
};
handstable.insertCol = function(row){
    handstable.table.alter("insert_col",row)
};
handstable.removeCol = function(row){
    handstable.table.alter("remove_col",row)
};

handstable.setCellData = function(msg){
    handstable.table.setDataAtCell(msg.row,msg.col,msg.newValue)
};

handstable.pos2Axis = function(col,row){
    return  util.format("{0}{1}",handstable.table.getColHeader(col),this.table.getRowHeader(row));
};


function downCsv() {
    handstable.exportPlugin.downloadFile('csv', {filename: handstable.tableName});
};

function downExcel() {
    handstable.exportPlugin.downloadFile('csv', {filename: handstable.tableName});
};

function showData() {
    console.log(handstable.table.getData());
}




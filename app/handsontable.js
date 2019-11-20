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

handstable.new = function(data,isReadOnly){
    if (data.length ==0){
        data = [[]]
    }
    handstable.setStatue(StatusEnum.EDITOR);
    handstable.table =  new Handsontable(handstable.container, {
        data: data,
        readOnly:isReadOnly,
        width:'100%',
        height:window.innerHeight - 65,
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
            cellSelected({row:row, col:col, row2:row2, col2:col2})
        },
        afterRender:function (isForced) {
            //console.log("afterRender",isForced);
            $.each(handstable.selected,function (name,item) {
                handstable.showUser(item.selected.row,item.selected.col,item.div)
            });
        }
    });
    handstable.customBordersPlugin = handstable.table.getPlugin('customBorders');
    handstable.exportPlugin = handstable.table.getPlugin('exportFile')
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

handstable.setHistory = function(msg){
    let tmp = `<div class="history-div-nav-li" onclick="lookHistory({0})">
            <div style="width: 100%;height: 5px;"></div>
            &nbsp;&nbsp;<span style="font-weight: bold">第{1}版</span><span class="rbBtn" onclick="rollback({2})">还原</span><br>
            &nbsp;&nbsp;<span style="font-weight: lighter">{3}</span><br>
            {4}
            <div style="width: 100%;height: 5px;"></div>
        </div>`
    let userTmp = `&nbsp;&nbsp;<span style="font-size: 14px">{0}</span><br>`
    let list = document.getElementById('history-div-nav');
    list.innerHTML = "";
    let str = "";
    for(let i = 0;i < msg.length;i++) {
        let item = msg[i];
        let userStr = "";
        for (let j = 0;j < item.users.length;j++){
            userStr += util.format(userTmp,item.users[j])
        }
        str += util.format(tmp,item.version,item.version,item.version,item.date,userStr)
    }
    list.innerHTML = str ;
};

handstable.setSelected = function(name,selected){
    // 将上一次选中清空
    let value = handstable.selected[name];
    if (value){
        let lastSelected = value.selected
        handstable.customBordersPlugin.clearBorders([[lastSelected.row,lastSelected.col,lastSelected.row2,lastSelected.col2]]);

        let cell = handstable.table.getCell(lastSelected.row,lastSelected.col);
        if (cell){
            cell.removeChild(value.div)
        }
        delete handstable.selected[name];
    }

    let colr = handstable.userColor[name];
    if (!colr){
        let len = Object.keys(handstable.userColor).length;
        colr = util.getColor(len);
        handstable.userColor[name] = colr
    }

    let newCells = [[selected.row,selected.col,selected.row2,selected.col2]];
    handstable.customBordersPlugin.setBorders(newCells,util.makeBorderStyle(colr));
    let newDiv = handstable.makeDiv(name,colr);
    handstable.selected[name] = {selected:selected,div:newDiv};
    handstable.showUser(selected.row,selected.col,newDiv)
};

handstable.makeDiv = function(name,color) {
    let divElement = document.createElement("div")
    divElement.innerHTML = name;
    divElement.style.position = "absolute";
    divElement.style.marginTop = "-23px";
    divElement.style.marginLeft = "-5px";
    divElement.style.paddingLeft = "5px";
    divElement.style.paddingRight = "5px";
    divElement.style.zIndex = "1000";
    divElement.style.height = "23px";
    divElement.style.width = "auto";
    divElement.style.color = "black";
    divElement.style.background= color;
    return divElement
};

handstable.showUser =  function(row,column,div_) {
    let cell = handstable.table.getCell(row,column);
    if (cell){
        cell.appendChild(div_)
    }
};

handstable.insertRow = function(idx){
    handstable.table.alter("insert_row",idx)
};
handstable.removeRow = function(idx){
    handstable.table.alter("remove_row",idx)
};
handstable.insertCol = function(idx){
    handstable.table.alter("insert_col",idx)
};
handstable.removeCol = function(idx){
    handstable.table.alter("remove_col",idx)
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
    let ws =  XLSX.utils.aoa_to_sheet(handstable.table.getData());
    let wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Sheet1');
    XLSX.writeFile(wb, util.format("{0}.xlsx",handstable.tableName))
};

function showData() {
    let data = handstable.table.getData();
    let msg = "";
    for (let i = 0;i < data.length;i++){
        let row = data[i];
        for (let j = 0;j < row.length;j++){
            if (row[j] == null || row[j] ==""){
                msg+="  ,"
            }else {
                msg+=row[j]+" ,"
            }
        }
        msg += "\n"
    }
    //console.log(msg);
    showTips(msg,5000)
}


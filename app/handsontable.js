var handstable = {};

handstable.container = null;
handstable.data = null;
handstable.table = null;

handstable.init = function (container) {
    handstable.container = container
};

handstable.new = function(){
    handstable.table =  new Handsontable(handstable.container, {
        data: handstable.data,

        width:'1200px',
        height:'800px',

        minCols:46, //最小列数
        //maxCols: 50, //最大列数
        minRows:50 , //最小行数
        //maxRows: 40 , //最大行数
        minSpareCols: 10, //添加空列
        minSpareRows :10,//添加空行

        colHeaders:true, //显示列表头， 默认false， 可取 true/fals/数组 ，当值为数组时，列头为数组的值
        rowHeaders:true,　//显示行表头， 默认 false， 可取 true/fals/数组，当值为数组时，行头为数组的值

        mergeCells: true, //表示允许单元格合并

        fixedRowsTop : 0,    //冻结行（固定顶部开始算起指定行不随垂直滚动条滚动；
        fixedColumnsLeft : 0, //冻结列（固定右边开始算起指定行不随水平滚动条滚动）；
        manualColumnFreeze:true,  //设置true后在单元格中右击出现一个菜单，此菜单会多出一个“Unfteeze this columu”的选项，再次点击就会取消冻结动作。 默认为false

        manualColumnResize: true, //允许拖曳列表头，默认为false
        manualRowResize: true, //允许拖曳行表头，默认为false

        currentRowClassName:"cur", //给当前行设置样式名（鼠标点击某个单元格，则整行给样式）
        currentColClassName:"cur", //给当前行列设置样式名（鼠标点击某个单元格，则整行给样式）

        autoColumnSize: true, //当值为true且列宽未设置时，自适应列大小
        //columnWidth:50,

        //columnSorting:true, //通过点击列表头进行排序（没有图标）
        //colWidths:[10,5,50],
        //rowHeights: [50, 40, 100],

        stretchH:"last", // 可取 last/all/none last：延伸最后一列，all：延伸所有列，none默认不延伸。

        manualColumnMove: true, //true/false 当值为true时，列可拖拽移动到指定列
        manualRowMove: true, //true/false当值为true时，行可拖拽至指定行

        manualColumnFreeze:true,

        undo: true,
        redo: true,

        contextMenu:true, //显示表头的下拉菜单默认false 可取 true/false/自定义数组
        contextMenu: {    //内容部分的menu 对功能汉化
            items: {
                'row_above': {
                    name: '上方插入一行'
                },
                'row_below': {
                    name: '下方插入一行'
                },
                "col_left": {
                    name: '左方插入列',
                },
                "col_right": {
                    name: '右方插入列'
                },
                "remove_row": {
                    name: '删除行',
                },
                "remove_col": {
                    name: '删除列',
                },
                'separator': Handsontable.plugins.ContextMenu.SEPARATOR,
                'clear_custom': {
                    name: '清除数据',
                    callback: function() {
                        this.clear();
                    }
                }
            }
        },
        afterChange:function(changes,source){
            if (changes){
                var cellValues = new Array();
                for (var i = 0;i < changes.length;i++){
                    let item = changes[i]
                    console.log("change",item)
                    if (item[2] != null && item[2] != ""){
                        cellValues.push( {col:item[1],row:item[0],oldValue:item[2],newValue:item[3]})
                    }else {
                        if (item[3] != null && item[3] != ""){
                            cellValues.push( {col:item[1],row:item[0],oldValue:item[2],newValue:item[3]})
                        }
                    }
                }
                if (cellValues.length >0 ){
                    socket.send({cmd:"setCellValues",cellValues:cellValues})
                }
            }
        },
        afterSelectionEnd: function(row, col, row2, col2) {
            console.log("afterSelectionEnd",row,col,row2,col2)
            var selections = new Array()
            for (var i = row;i <= row2;i++){
                for (var j = col;j <= col2;j++){
                    selections.push({col:j,row:i})
                }
            }
            if (selections.length >0 ){
                socket.send({cmd:"cellSelected",selected:selections})
            }
        },
    });
};

handstable.setData = function(data){
    this.data = data;
    if (this.table) {
        this.table.loadData(data)
    }else {
        this.new()
    }
};

handstable.setSelected = function(selected){
    for (var i  = 0;i < selected.length;i++){
        var item = selected[i]
        //handstable.table.selectRows(item.row)
        let source = this.table.getRowHeader(item.row)
        //this.table.set
        console.log("selected",source)
        console.log("selected",item.row,item.col,this.table.getDataAtCell(item.row,item.col))
        /*let tr = document.getElementsByClassName("current")
        console.log("selected",tr,tr.length)
        for (i = 0; i < tr.length; i++) {
            tr[i].style.backgroundColor='#005EFF';
        }*/
        //handstable.table.show
       //handstable.table.setDataAtCell(item.col,item.row,"ss")

    }
}

handstable.pos2Axis = function(col,row){
    return  String.format("{0}{1}",this.table.getColHeader(col),this.table.getRowHeader(row));
};

handstable.addHook = function(key,callback){
    handstable.table.addHook(key,callback)
}

handstable.addHookOnce = function(key,callback){
    handstable.table.addHookOnce(key,callback)
}


// 字符串格式化
String.format = function(src){
    if (arguments.length == 0) return null;
    var args = Array.prototype.slice.call(arguments, 1);
    return src.replace(/\{(\d+)\}/g, function(m, i){
        return args[i];
    });
};




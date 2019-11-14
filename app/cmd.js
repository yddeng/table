
function openTable(tableName,userName) {
    socket.send({cmd:'openTable',tableName:tableName,userName:userName});
}

function insertRow(idx) {
    socket.send({cmd:'insertRow',index:idx});
}

function removeRow(idx) {
    socket.send({cmd:'removeRow',index:idx});
}

function insertCol(idx) {
    socket.send({cmd:'insertCol',index:idx});
}

function removeCol(idx) {
    socket.send({cmd:'removeCol',index:idx});
}

function setCellValues(data ) {
    socket.send({cmd:'setCellValues',cellValues:data});
}

function cellSelected(data) {
    socket.send({cmd:'cellSelected',selected:data});
}

function saveTable() {
    socket.send({cmd:"saveTable"})
}

function lookHistory() {
    let txt = $("#rbv").val();
    let v = parseInt(txt)
    socket.send({cmd:"lookHistory",version:v})
}

/*****************************************************************************************/
var dispatcher = {};
dispatcher.handler = {};

dispatcher.DispatchMessage = function(msg) {
    var handler = dispatcher.handler[msg.cmd];
    console.log("DispatchMessage",msg);
    if(handler){
        handler(msg);
    }
};

dispatcher.handler["pushErr"] = function(msg) {
    util.alert(util.format("cmd:{0},errMsg:{1}",msg.doCmd,msg.errMsg))
};

dispatcher.handler["openTable"] = function(msg) {
    handstable.new(msg.data)
    handstable.setVersion(msg.version)
};

dispatcher.handler["pushCellSelected"] = function(msg) {
    handstable.setSelected(msg.selected)
};

dispatcher.handler["pushCellData"] = function(msg) {
    //handstable.setData(msg.data)
};

dispatcher.handler["pushSaveTable"] = function(msg) {
    handstable.setData(msg.data)
    handstable.setVersion(msg.version)
    util.alert("文件已保存")
};

dispatcher.handler["pushAll"] = function(msg) {
    handstable.setData(msg.data)
    handstable.setVersion(msg.version)
};

dispatcher.handler["lookHistory"] = function(msg) {
    handstable.setVersion(msg.version)
    handstable.setData(msg.data)
};


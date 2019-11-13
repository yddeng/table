
function openTable(tableName,userName) {
    socket.send({cmd:'openTable',tableName:tableName,userName:userName});
}

function insertRow(idx) {
    socket.send({cmd:'insertRow',rowIndex:idx});
}

function removeRow(idx) {
    socket.send({cmd:'removeRow',rowIndex:idx});
}

function insertCol(idx) {
    socket.send({cmd:'insertCol',colIndex:idx});
}

function removeCol(idx) {
    socket.send({cmd:'removeCol',colIndex:idx});
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

function rollback() {
    let txt = $("#rbv").val();
    let v = parseInt(txt)
    socket.send({cmd:"rollback",now:handstable.version, goto:v})
}

/*****************************************************************************************/
var dispatcher = {};
dispatcher.handler = {};

dispatcher.error = function(msg){
    console.log("error",msg.code,msg.msg)
};

dispatcher.DispatchMessage = function(msg) {
    var handler = dispatcher.handler[msg.cmd];
    console.log("DispatchMessage",msg);
    if(handler){
        handler(msg);
    }
};

dispatcher.handler["openTable"] = function(msg) {
    handstable.new(msg.data)
    handstable.setVersion(msg.version)
};

dispatcher.handler["pushCellSelected"] = function(msg) {
    //handstable.setData(msg.data)
};

dispatcher.handler["pushCellData"] = function(msg) {
    //handstable.setData(msg.data)
};

dispatcher.handler["pushData"] = function(msg) {
    handstable.setData(msg.data)
    handstable.setVersion(msg.version)
};

dispatcher.handler["rollback"] = function(msg) {
    if (msg.ok == 1){
        handstable.setVersion(msg.version)
        handstable.setData(msg.data)
    }
};


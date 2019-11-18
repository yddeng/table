
function openTable(tableName,userName) {
    socket.send({cmd:'openTable',tableName:tableName,userName:userName});
}

function insertRow(idx) {
    if (handstable.status === StatusEnum.EDITOR){
        socket.send({cmd:'insertRow',index:idx});
    }
}

function removeRow(idx) {
    if (handstable.status === StatusEnum.EDITOR) {
        socket.send({cmd: 'removeRow', index: idx});
    }
}

function insertCol(idx) {
    if (handstable.status === StatusEnum.EDITOR) {
        socket.send({cmd: 'insertCol', index: idx});
    }
}

function removeCol(idx) {
    if (handstable.status === StatusEnum.EDITOR) {
        socket.send({cmd: 'removeCol', index: idx});
    }
}

function setCellValues(data ) {
    if (handstable.status === StatusEnum.EDITOR) {
        socket.send({cmd: 'setCellValues', cellValues: data});
    }
}

function cellSelected(data) {
    if (handstable.status === StatusEnum.EDITOR) {
        socket.send({cmd: 'cellSelected', selected: data});
    }
}

function saveTable() {
    if (handstable.status === StatusEnum.EDITOR) {
        socket.send({cmd: "saveTable"})
    }
}

function versionList() {
    if (handstable.status === StatusEnum.EDITOR) {
        socket.send({cmd: "versionList"})
    }
}

function lookHistory(v) {
    socket.send({cmd:"lookHistory",version:v})
}

function backEditor() {
    if (handstable.status === StatusEnum.LOOK) {
        socket.send({cmd: "backEditor"})
    }
}

function rollback(v) {
    var ev=ev||event;
    if(ev && ev.stopPropagation){
        ev.stopPropagation();  //非IE下 它支持W3C的stopPropagation()方法
    } else {
        window.event.cancelBubble = true;  //IE的方式来取消事件冒泡
    }
    socket.send({cmd:"rollback",version:v})
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
    util.alert(util.format("cmd:{0}\nerrMsg:{1}",msg.doCmd,msg.errMsg))
};

dispatcher.handler["pushAll"] = function(msg) {
    if (handstable.status === StatusEnum.EDITOR) {
        handstable.setData(msg.data);
        handstable.setVersion(msg.version)
    }
};

dispatcher.handler["openTable"] = function(msg) {
    handstable.new(msg.data);
    handstable.setVersion(msg.version);
};

dispatcher.handler["cellSelected"] = function(msg) {
    if (handstable.status === StatusEnum.EDITOR) {
        handstable.setSelected(msg.userName, msg.selected)
    }
};

dispatcher.handler["saveTable"] = function(msg) {
    if (handstable.status === StatusEnum.EDITOR) {
        handstable.setData(msg.data)
        handstable.setVersion(msg.version)
        util.alert("文件已保存")
    }
};

dispatcher.handler["insertRow"] = function(msg) {
    if (handstable.status === StatusEnum.EDITOR) {
        handstable.insertRow(msg.index)
    }
};
dispatcher.handler["removeRow"] = function(msg) {
    if (handstable.status === StatusEnum.EDITOR) {
        handstable.removeRow(msg.index)
    }
};
dispatcher.handler["insertCol"] = function(msg) {
    if (handstable.status === StatusEnum.EDITOR) {
        handstable.insertCol(msg.index)
    }
};
dispatcher.handler["removeCol"] = function(msg) {
    if (handstable.status === StatusEnum.EDITOR) {
        handstable.removeCol(msg.index)
    }
};

dispatcher.handler["versionList"] = function(msg) {
    handstable.setHistory(msg.list);
};

dispatcher.handler["lookHistory"] = function(msg) {
    handstable.setVersion(msg.version);
    handstable.setData(msg.data);
    handstable.setStatue(StatusEnum.LOOK);
};

dispatcher.handler["backEditor"] = function(msg) {
    handstable.setVersion(msg.version);
    handstable.setData(msg.data);
    handstable.setStatue(StatusEnum.EDITOR);
};

dispatcher.handler["rollback"] = function(msg) {
    handstable.setVersion(msg.version);
    handstable.setData(msg.data);
    util.alert(util.format("版本已回退到:{0}",msg.version))
};

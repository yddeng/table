
function openTable(tableName,token) {
    socket.send({cmd:'openTable',tableName:tableName,token:token});
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
    let ev=ev||event;
    if(ev && ev.stopPropagation){
        ev.stopPropagation();  //非IE下 它支持W3C的stopPropagation()方法
    } else {
        window.event.cancelBubble = true;  //IE的方式来取消事件冒泡
    }
    socket.send({cmd:"rollback",version:v})
}

function talk() {
    if (event.keyCode == 13){
        event.preventDefault();//屏蔽enter对系统作用。按后增加\r\n等换行
        let inp = document.getElementById("chat-input");
        let txt = inp.value;
        if (txt !== "") {
            inp.value = "";
            socket.send({cmd: "talk",msg:txt})
        }
    }
}

/*****************************************************************************************/
var dispatcher = {};
dispatcher.handler = {};

dispatcher.DispatchMessage = function(msg) {
    let handler = dispatcher.handler[msg.cmd];
    console.log("DispatchMessage",msg);
    if(handler){
        handler(msg);
    }
};

dispatcher.handler["pushErr"] = function(msg) {
    showTips(util.format("cmd:{0}\nerrMsg:{1}",msg.doCmd,msg.errMsg),3000);
};

dispatcher.handler["pushAll"] = function(msg) {
    if (handstable.status === StatusEnum.EDITOR) {
        handstable.table.loadData(msg.data)
        handstable.table.render();
        handstable.setVersion(msg.version)
    }
};

dispatcher.handler["openTable"] = function(msg) {
    handstable.new(msg.data);
    handstable.setVersion(msg.version);

    let txt = util.tokenName(socket.token)+"、";
    txt += util.format("等{0}人在线",Object.keys(handstable.userColor).length+1);
    chatHead(txt)
};

dispatcher.handler["cellSelected"] = function(msg) {
    if (handstable.status === StatusEnum.EDITOR) {
        handstable.setSelected(msg.userName, msg.selected)
    }
};

dispatcher.handler["saveTable"] = function(msg) {
    if (handstable.status === StatusEnum.EDITOR) {
        handstable.table.loadData(msg.data)
        handstable.setVersion(msg.version);
        showTips("文件已保存",3000);
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
    handstable.table = null;
    handstable.new(msg.data,true);
    handstable.setStatue(StatusEnum.LOOK);
};

dispatcher.handler["backEditor"] = function(msg) {
    handstable.setVersion(msg.version);
    handstable.table = null;
    handstable.new(msg.data);
    handstable.setStatue(StatusEnum.EDITOR);
};

dispatcher.handler["rollback"] = function(msg) {
    handstable.setVersion(msg.version);
    handstable.table.loadData(msg.data);
    showTips(util.format("版本已还原至:{0}",msg.version),3000);
};

dispatcher.handler["userEnter"] = function(msg) {
    //showTips(msg.userName,2000);
    // 添加用户
    let len = Object.keys(handstable.userColor).length;
    let colr = util.getColor(len);
    //console.log(len,colr);
    handstable.userColor[msg.userName] = colr;

    let txt = util.tokenName(socket.token)+"、";
    $.each(handstable.userColor, function(name,value){
        txt += name+"、"
    });
    txt += util.format("等{0}人在线",Object.keys(handstable.userColor).length+1);
    chatHead(txt)
};

dispatcher.handler["userLeave"] = function(msg) {
    //showTips(msg.userName,2000);
    let name = msg.userName;
    // 清空用户
    let value = handstable.selected[name];
    if (value){
        let lastSelected = value.selected;
        handstable.customBordersPlugin.clearBorders([[lastSelected.row,lastSelected.col,lastSelected.row2,lastSelected.col2]]);

        let cell = handstable.table.getCell(lastSelected.row,lastSelected.col);
        if (cell){
            cell.removeChild(value.div)
        }
        delete handstable.selected[name];
    }

    let colr = handstable.userColor[name];
    if (colr){
        //console.log(colr)
        delete handstable.userColor[name];
    }

    let txt = util.tokenName(socket.token)+"、";
    $.each(handstable.userColor, function(name,value){
        txt += name+"、"
    });
    txt += util.format("等{0}人在线",Object.keys(handstable.userColor).length+1);
    chatHead(txt)
};

dispatcher.handler["talk"] = function (msg) {
  talkMsg(msg)
};

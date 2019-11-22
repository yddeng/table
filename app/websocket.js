var socket = socket || {};

socket.ws = null;
socket.tableName = null;
socket.token = null;

socket.init = function(tableName,token){
    socket.tableName = tableName;
    socket.token = token;
    socket.connect(wsAddr)
};

socket.connect = function(url) {
    socket.ws = new WebSocket(url);
    socket.ws.onopen = socket.onopen;
    socket.ws.onmessage = socket.onmessage;
    socket.ws.onclose = socket.onclose;
};

socket.onopen = function(){
    console.log("connect ok",socket.tableName,socket.token);
    openTable(socket.tableName,socket.token)
};

socket.onmessage = function(evt){
    let msg = JSON.parse(evt.data);
    dispatcher.DispatchMessage(msg);
};

socket.onclose = function(e) {
    if(socket.ws) {
        let ws = socket.ws;
        socket.ws = null;
        ws.close();
        console.log(e)
    }
};

socket.send = function(msg) {
    if(socket.ws) {
        console.log("send",msg);
        return socket.ws.send(JSON.stringify(msg));
    }else {
        socket.connect(wsAddr)
    }
};
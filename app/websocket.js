var socket = socket || {}

socket.ws = null;
socket.url = null;
socket.fileName = null;
socket.userName = null;

socket.init = function(url,fileName,userName){
    socket.url = url;
    socket.fileName = fileName;
    socket.userName = userName;
    socket.connect(url)
};

socket.connect = function(url) {
    socket.ws = new WebSocket(url);
    socket.ws.onopen = socket.onopen;
    socket.ws.onmessage = socket.onmessage;
    socket.ws.onclose = socket.onclose;
};

socket.onopen = function(){
    console.log("connect ok");
    console.log(socket.url,socket.fileName,socket.userName)
    openTable(socket.fileName,socket.userName)
};

socket.onmessage = function(evt){
    var msg = JSON.parse(evt.data);
    dispatcher.DispatchMessage(msg);
};

socket.onclose = function(e) {
    if(socket.ws) {
        var ws = socket.ws;
        socket.ws = null;
        ws.close();
        console.log(e)
    }
};

socket.send = function(msg) {
    if(socket.ws) {
        return socket.ws.send(JSON.stringify(msg));
    }else {
        socket.connect(this.url)
    }
};
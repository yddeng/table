var socket = socket || {};

socket.ws = null;
socket.url = "ws://10.128.2.123:4545/table";
socket.fileName = null;
socket.userName = null;

socket.init = function(fileName,userName){
    socket.fileName = fileName;
    socket.userName = userName;
    socket.connect(socket.url)
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
        console.log("send",msg)
        return socket.ws.send(JSON.stringify(msg));
    }else {
        socket.connect(this.url)
    }
};
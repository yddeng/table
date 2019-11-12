var message = {};

message.handler = {};

message.error = function(msg){
	console.log("error",msg.code,msg.msg)
};

message.handler["pushData"] = function(event) {
	handstable.setData(event.data)
	var selected = event.cellLocked
	if (selected) {
		handstable.setSelected(selected)
	}
};

message.DispatchMessage = function(msg) {
	var handler = message.handler[msg.cmd];
	console.log("DispatchMessage",msg.cmd);
	if(handler){
		handler(msg);
	}
};



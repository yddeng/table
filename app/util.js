var util = {};

util.alert = function (msg) {
    alert(msg)
};

// 字符串格式化
util.format = function(src){
    if (arguments.length == 0) return null;
    var args = Array.prototype.slice.call(arguments, 1);
    return src.replace(/\{(\d+)\}/g, function(m, i){
        return args[i];
    });
};

// url
util.getUrlParam =  function(name) {
    let reg = new RegExp('(^|&)'+ name + '=([^&]*)(&|$)');
    let r = window.location.search.substr(1).match(reg);
    if(r!=null){
        return unescape(r[2])
    }
    return null
};

var util = {};

util.colors = ["#E74C3C","#9B59B6","#1ABC9C","#F5B041"];
util.cmdString = {};

util.getStatusString = function(status){
    if (status === StatusEnum.NONE){
        return "错误的状态"
    }else  if (status === StatusEnum.EDITOR){
        return "正在编辑"
    }else if (status === StatusEnum.LOOK){
        return  "正在查看历史版本"
    }
};

util.alert = function (msg) {
    alert(msg)
};

// 字符串格式化
util.format = function(src){
    if (arguments.length == 0) return null;
    let args = Array.prototype.slice.call(arguments, 1);
    return src.replace(/\{(\d+)\}/g, function(m, i){
        return args[i];
    });
};

// url
util.getUrlParam =  function(name) {
    let reg = new RegExp('(^|&)'+ name + '=([^&]*)(&|$)');
    let result= window.location.search.substr(1).match(reg);
    return result?decodeURIComponent(result[2]):null;
};

util.getColor = function(idx) {
    return  util.colors[idx%util.colors.length];
};

util.makeBorderStyle = function(colr){
    return {right:{width: 2,color:colr},left:{width: 2,color:colr},
        top:{width: 2,color:colr},bottom:{width: 2,color:colr}};
};


//弹出一个输入框，输入一段文字，可以提交
util.prom = function(msg) {
    let name = prompt(msg, "");
    if (name){
        return name
    }
};

util.localTime = function(){
    let myDate = new Date();
    return myDate.toLocaleString()
};

// 设置cookie的函数  （名字，值，过期时间（天））
util.setCookie = function (cname, cvalue, exdays) {
    let d = new Date();
    d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
    let expires = "expires=" + d.toUTCString();
    document.cookie = cname + "=" + cvalue + "; " + expires;
};

//获取cookie
//取cookie的函数(名字) 取出来的都是字符串类型 子目录可以用根目录的cookie，根目录取不到子目录的 大小4k左右
util.getCookie = function(cname) {
    let name = cname + "=";
    let ca = document.cookie.split(';');
    for(let i=0; i<ca.length; i++)
    {
        let c = ca[i].trim();
        if (c.indexOf(name)===0) return c.substring(name.length,c.length);
    }
    return "";
};

util.tokenName = function(token){
    let name = token.split("@")[1];
    return name
};

util.httpGet = function(url,success,error){
    $.ajax({
        url:url,
        type: "get",
        async: true,
        success: success,
        error: error
    });
};

util.httpPost = function(url,data,success,error){
    $.ajax({
        url:url,
        type: "post",
        async: true,
        dataType: "json",
        data:data,
        success: success,
        error: error
    });
};

//弹出一个询问框，有确定和取消按钮
function firm() {
    //利用对话框返回的值 （true 或者 false）
    if (confirm("你确定提交吗？")) {
        alert("点击了确定");
    }
    else {
        alert("点击了取消");
    }

}

util.exportExcel = function(name,data){
    let ws =  XLSX.utils.aoa_to_sheet(data);
    let wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Sheet1');
    XLSX.writeFile(wb, util.format("{0}.xlsx",name))
};
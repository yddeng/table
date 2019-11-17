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

util.getColor = function(idx) {
    let colr = util.colors[idx%util.colors.length];
    return {right:{width: 2,color:colr},left:{width: 2,color:colr},
        top:{width: 2,color:colr,text:"sssss"},bottom:{width: 2,color:colr}};
}

//弹出一个输入框，输入一段文字，可以提交
function prom() {
    var name = prompt("请输入您的名字", ""); //将输入的内容赋给变量 name ，

    //这里需要注意的是，prompt有两个参数，前面是提示的话，后面是当对话框出来后，在对话框里的默认值
    if (name)//如果返回的有内容
    {
        alert("欢迎您：" + name)
    }

}

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
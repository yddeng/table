       function setCookie(cname, cvalue, exdays) { // 设置cookie的函数  （名字，值，过期时间（天））
        var d = new Date();
        d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
        var expires = "expires=" + d.toUTCString();
        document.cookie = cname + "=" + cvalue + "; " + expires;
      }
      //获取cookie
      function getCookie(cname) { //取cookie的函数(名字) 取出来的都是字符串类型 子目录可以用根目录的cookie，根目录取不到子目录的 大小4k左右
     	let arr = document.cookie.split("; ")
		let val
		for(var i=0;i<arr.length;i++){
			if(arr[i].split('=')[0]==cname){
				 val = arr[i].split('=')[1]
				 break;
			}else{
				val = ''
			}
		}
		return val
      }
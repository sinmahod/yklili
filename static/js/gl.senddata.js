(function ($) {
 /**
 * 提交表单 封装类
 * author: gl
 * time: 2016/11/21 21:07:11
 * email: sinmahod@qq.com
 */
    window.SendData = function () {
    return {
        //Post
        Post: function(controller,data,fn){
            if(data){
                data["IsSendData"] = true;
            }else{
                data = {"IsSendData" : true}
            }
            $.post(controller,data,function(result){
                if (typeof(result["STATUS"]) != "undefined" ) {
                    if (result["STATUS"] == 101) {
                        BootFrame.alert("会话超时请重新登录",null,"超时",true); 
                        return;
                    }
                }
                if(fn){
                    fn(result);
                }
            }).error(function() { 
                BootFrame.alert("服务器发生错误",null,"错误",true); 
            });
        },//<! Post >
        Get: function(controller,data,fn){
            if(data){
                data["IsSendData"] = true;
            }else{
                data = {"IsSendData" : true}
            }
            $.get(controller,data,function(result){
                if (typeof(result["STATUS"]) != "undefined" ) {
                    if (result["STATUS"] == 101) {
                        BootFrame.alert("会话超时请重新登录",null,"超时",true); 
                        return;
                    }
                }
                if(fn){
                    fn(result);
                }
            }).error(function() {
                BootFrame.alert("服务器发生错误",null,"错误",true); 
            });
        }//<! Get >
    }//<! return >
   }();

})(jQuery);
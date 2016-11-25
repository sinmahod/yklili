(function ($) {
 /**
 * 校验表单 封装类
 * author: gl
 * time: 2016/11/25 17:17:18
 * email: sinmahod@qq.com
 */
    window.VerifyForm = function () {
    return {
        //hasError
        hasError: function(formid){
           var s = $("#"+id );
                    if (s.length == 0) {
                        BootFrame.alert("未找到form");
                        return false;
                    }
                    if (s[0].length == 0 ){
                        BootFrame.alert("form中不存在表单");
                        return false;
                    }

                    var a = s[0];
                    for (var i = 0 ; i < a.length; i ++){
                        //判断表单类型
                        var _ver = $(a[i]).attr('verify');
                        if (_ver == 'notnull') {
                            
                        }
                    }
                    return true;
        }//<! hasError >
    }//<! return >
   }();

})(jQuery);
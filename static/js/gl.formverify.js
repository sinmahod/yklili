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
           var s = $("#"+formid );
                    if (s.length == 0) {
                        BootFrame.alert("未找到form");
                        return false;
                    }
                    if (s[0].length == 0 ){
                        BootFrame.alert("form中不存在表单");
                        return false;
                    }

                    if (!s.valid()){
                        BootFrame.alert("请按照规则输入表单");
                        return false;
                    }
                    return true;
        }//<! hasError >
    }//<! return >
   }();

})(jQuery);
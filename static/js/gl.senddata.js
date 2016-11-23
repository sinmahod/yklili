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
            var jstr = null;
            if (data != null){
                jstr = data.toJSON();
            }
            $.post(controller,jstr,function(result){
                if(fn){
                    fn(result);
                }
            });
        },//<! Post >
        Get: function(controller,data,fn){
            var jstr = null;
            if (data != null){
                jstr = data.toJSON();
            }
            $.get(controller,jstr,function(result){
                if(fn){
                    fn(result);
                }
            });
        }//<! Get >
    }//<! return >
   }();

})(jQuery);
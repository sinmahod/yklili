(function ($) {
// 初始化init
/* 使用方式
 * <div>
 *  <ul class="ace-thumbnails clearfix" list="true" init-medhod="/data/image/List">
 *          <li>
 *               <a href="${FilePath}${FileNewName}" data-rel="colorbox">
 *                   <img width="150" height="150" class="lazy" data-original="${FilePath}${FileNewName}" />
 *               </a>
 *
 *               <div class="tools tools-bottom">
 *                   <a href="#">
 *                       <i class="ace-icon fa fa-link"></i>
 *                   </a>
 *
 *                   <a href="#">
 *                       <i class="ace-icon fa fa-pencil"></i>
 *                   </a>
 *
 *                   <a href="#">
 *                       <i class="ace-icon fa fa-times red"></i>
 *                   </a>
 *               </div>
 *           </li>
 *       </ul>
 *
 *   </div>
 *   <div align="center">
 *       <ul class="pagination" id="pagebar"></ul>
 *   </div>
 *
 *  <script type="text/javascript">
 *      jQuery(function($) {
 *           Init({list:listID,pagebar:pagebarID,html:html,fn:function(){}});
 *      });
 *  </script>   
 *  fn每次翻页后的执行函数
 */

var DataListMap = {};

window.Init = function(options){
    var $init,html,fn;

    if (typeof(options) != "undefined"){
        if (typeof(options["id"]) != "undefined" ) {
            if(typeof(options["id"]) == "string"){
                $init = $('#'+options["id"]);
            }else{
                $init = options["id"];
            }
        }
        if (typeof(options["html"]) != "undefined" ) html = options["html"];
        fn = typeof(options["fn"]) != "undefined" ? options["fn"] : null;
        if ($init) start();
    }
    if (!$init) {
         //检查所有带有init-medhod属性的标签，开始解析当前标签（包含子dom）
        $("[init-method]").each(function(){
            $init = $(this);
            start();
        });
    }

    if($init && !html){
         html = $init.html();
    }  
    
    function start(){
        var url = $init.attr("init-method");
        if(!url){
            console.log('Error:没有找到init-method');
            return;
        }

        readData(url,function(result){
            newHtml = renderList(result);
            if (newHtml != null){
                $init.html(newHtml);
                if (fn){
                    fn();
                }
            }
        });
    }

       //渲染列表
    function renderList(result){
        if (result){
            var reg = /\${([^{}]+)}/g;
            var tempHtml = html.replace(reg, function (match, name) {
                var data = result[name];
                return data ? data : "";
            });
            return tempHtml;
        }
        return null;
    }


    //ajax同步方式读取数据
    //url = 请求的url
    //page = 请求的页码
    //size = 每页请求的数据量
    //fn = 后执行函数
    function readData(url,fn){
        $.ajax({
            url : url,
            type : 'post',
            success : function(result){
                if (result != null ){
                    fn(result);
                }else{
                    console.log('Error:服务器发生错误，未获取到数据');
                }
            },
            fail:function(){
                console.log('Error:服务器发生错误，未获取到数据');
            }
        });
    }   
 };

 })(jQuery);
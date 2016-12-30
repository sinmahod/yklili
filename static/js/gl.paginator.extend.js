(function ($) {
// 分页条paginator
/* 使用方式
 * <div>
 *  <ul class="ace-thumbnails clearfix" list="true" dataurl="/data/image/List" page="true" size="1">
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
 *           DataList({list:listID,pagebar:pagebarID,html:html,fn:function(){}});
 *      });
 *  </script>   
 *  fn每次翻页后的执行函数
 */

var DataListMap = {};

window.DataList = function(options){
    var $pagebar,$list,html,url,size,fn; 
    var cpage = 1;
    var total = 0;

    if (typeof(options) != "undefined"){
        if (typeof(options["pagebar"]) != "undefined" ) {
            if(typeof(options["pagebar"]) == "string"){
                $pagebar = $('#'+options["pagebar"]);
            }else{
                $pagebar = options["pagebar"];
            }
        }
        if (typeof(options["list"]) != "undefined" ) {
            if(typeof(options["list"]) == "string"){
                $list = $('#'+options["list"]);
            }else{
                $list = options["list"];
            }
        }
        if (typeof(options["html"]) != "undefined" ) html = options["html"];
        size = typeof(options["size"]) != "undefined" ? options["size"] : 10000;
        fn = typeof(options["fn"]) != "undefined" ? options["fn"] : null;
    }

    if (!$pagebar){
        $("[pagebar='true']").each(function(){
            $pagebar = $(this);
        });
        if (!$pagebar){
            throw new Error( 'not find pagebar' );
        }
    }


    if (!$list){
        //检查所有带有List属性的标签，开始循环当前标签（包含子dom）
        $("[list='true']").each(function(){
            $list = $(this);
            start();
        });
    }else{
        start();
    }

    function start(){
        var id = $list.attr("id");
        if(!id){
            console.log('Error:没有找到list元素的id');
            return;
        }
        if(!html){
             html = $list.html();
        }   
        DataListMap[id] = {template:html,pagebar:$pagebar,size:size,fn:fn};
        execute();
    }

    //第一次的执行方法
    function execute(){
        url = $list.attr("dataurl");
        if (!size){
            size = $list.attr("size");
        }
        var page = $list.attr("page");
        if (!page || page=="false"){
            size = 10000;
        }

        readData(url,1,function(result){
            newHtml = renderList(result);
            if (newHtml != null){
                $list.html(newHtml);
                if (newHtml != ""){
                     $("img.lazy").lazyload();
                    if (fn){
                        fn();
                    }
                    initPagebar($pagebar);
                }
            }
        });
    }
    
    

    //分页条初始化
    function initPagebar(pagebar){
        var n = 5;
        var w = $(document.body).width();
        if (w < 400) n = 3;
        else if (w < 500) n = 4;
        var options = {
            bootstrapMajorVersion: 3, //版本
            alignment: "center",//居中显示
            currentPage: cpage,//当前页码
            totalPages: total,//总页码
            numberOfPages: n,//最多显示几个页码按钮
            itemTexts: function (type, page, current) {
                    switch (type) {
                    case "first":
                        return "首页";
                    case "prev":
                        return "上一页";
                    case "next":
                        return "下一页";
                    case "last":
                        return "尾页";
                    case "page":
                        return page;
                    }
            },
            //点击事件
            onPageClicked: function (event, originalEvent, type, page) {
                readData(url,page,function(result){
                    newHtml = renderList(result);
                    if (newHtml != null){
                        $list.html(newHtml);
                        $("img.lazy").lazyload();
                        if (fn){
                            fn();
                        }
                    }
                });
            }
        };
        pagebar.bootstrapPaginator(options);
    }

    //ajax同步方式读取数据
    //url = 请求的url
    //page = 请求的页码
    //size = 每页请求的数据量
    //fn = 后执行函数
    function readData(url,page,fn){
        $.ajax({
            url : url,
            type : 'post',
            data : {page:page, rows:size},//这里使用json对象
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

    //渲染列表
    function renderList(result){
        if (result["rows"].length == 0){
            return "";
        }
        if (result["records"] > 0){
            total = result["total"];
            cpage = result["page"];
            var reg = /\${([^{}]+)}/g;
            var tempArr = [];
            $.each(result["rows"],function(idx,data){
                var tempHtml = html.replace(reg, function (match, name) {
                    return data[name];
                });
                tempArr.push(tempHtml);
            });
            return tempArr.join(''); 
        }
        return null;
    }

   
 };

 DataList.loadData = function(listid){
    var map = DataListMap[listid];
    DataList({
        list: $('#'+listid),
        pagebar: map["pagebar"],
        html: map["template"],
        size: map["size"],
        fn: map["fn"]
    });
 }

 })(jQuery);
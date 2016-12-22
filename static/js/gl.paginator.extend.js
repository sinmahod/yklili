(function ($) {
//展示图片Colorbox 分页条paginator
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
 *           DataList("pagebar");
 *      });
 *  </script>   
 */

var DataListMap = {};

window.DataList = function(pagebar,list,html){
    var cpage = 1;
    var total = 0;
    var size = 0;
    var url;

    if (!list){
        //检查所有带有List属性的标签，开始循环当前标签（包含子dom）
        $("[list='true']").each(function(){
            list = $(this);
            var id = list.attr("id");
            if(!id){
                BootFrame.alert("检查到页面中有list参数，但是所在元素没有id。",null,"错误",true);
                return;
            }
            html = list.html();
            DataListMap[id] = {template:html,pagebar:pagebar};
            execute();
        });
    }else{
        execute();
    }

    //第一次的执行方法
    function execute(){
        url = list.attr("dataurl");
        size = list.attr("size");
        var page = list.attr("page");
        if (!page || page=="false"){
            size = 10000;
        }

        readData(url,1,function(result){
            newHtml = renderList(result);
            if (newHtml != null){
                list.html(newHtml);
                if (newHtml != ""){
                     $("img.lazy").lazyload();
                    initColorBox();
                    initPagebar($('#'+pagebar));
                }
            }
        });
    }
    
    

    //分页条初始化
    function initPagebar(pagebar){
         var options = {
            bootstrapMajorVersion: 3, //版本
            alignment: "center",//居中显示
            currentPage: cpage,//当前页码
            totalPages: total,//总页码
            numberOfPages: 5,//最多显示几个页码按钮
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
                        list.html(newHtml);
                        $("img.lazy").lazyload();
                        initColorBox();
                    }
                    //initColorBox();
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
                    BootFrame.alert("服务器发生错误，未获取到数据",null,"错误",true);
                }
            },
            fail:function(){
                BootFrame.alert("服务器发生错误，未获取到数据",null,"错误",true); 
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

    //初始化ColorBox
    function initColorBox(){
        var $overflow = '';
        var colorbox_params = {
            rel: 'colorbox',
            reposition:true,
            scalePhotos:true,
            scrolling:false,
            previous:'<i class="ace-icon fa fa-arrow-left"></i>',
            next:'<i class="ace-icon fa fa-arrow-right"></i>',
            close:'&times;',
            current:'{current} of {total}',
            maxWidth:'100%',
            maxHeight:'100%',
            onOpen:function(){
                $overflow = document.body.style.overflow;
                document.body.style.overflow = 'hidden';
            },
            onClosed:function(){
                document.body.style.overflow = $overflow;
            },
            onComplete:function(){
                $.colorbox.resize();
            }
        };

        $('.ace-thumbnails [data-rel="colorbox"]').colorbox(colorbox_params);
        $("#cboxLoadingGraphic").html("<i class='ace-icon fa fa-spinner orange fa-spin'></i>");
        
        $(document).one('ajaxloadstart.page', function(e) {
            $('#colorbox, #cboxOverlay').remove();
        });
    }

 };

 DataList.loadData = function(listid){
    var map = DataListMap[listid];
    DataList(map["pagebar"],$('#'+listid),map["template"]);
 }

 })(jQuery);
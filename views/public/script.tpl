
        <script type="text/javascript">
            if('ontouchstart' in document.documentElement) document.write("<script src='/static/js/jquery.mobile.custom.min.js'>"+"<"+"/script>");
        </script>

        <!-- progressbar -->
        <script src="/static/js/jquery-ui.min.js"></script>
        <script src="/static/js/jquery.ui.touch-punch.min.js"></script>


        <!-- core -->
        <script src="/static/js/bootstrap.min.js"></script>

        {{str2html "<!--[if lte IE 8]>"}}
          <script src="/static/js/excanvas.min.js"></script>
        {{str2html "<![endif]-->"}}
        
        <!-- bootstrapDialog -->
        <script src="/static/js/bootstrap-dialog.min.js"></script>
        <script src="/static/js/jquery.gritter.min.js"></script>
        <script src="/static/js/gl.bootstrap-dialog.extend.js"></script>

        <!-- formverify -->
        <script src="/static/js/bootstrapValidator.js"></script>
        <script src="/static/js/gl.validate.js"></script>

        <!-- pagebar -->
        <script src="/static/js/bootstrap-paginator.min.js"></script>
        <script src="/static/js/gl.paginator.extend.js"></script>

        <!-- init -->
        <script src="/static/js/gl.init.js"></script>

        <!-- webUploader -->
        <script src="/static/webuploader/webuploader.nolog.min.js"></script>

        <!-- image view -->
        <script src="/static/js/jquery.lazyload.js"></script>
        <script src="/static/js/jquery.colorbox.js"></script>

        <!-- colorpicker -->
        <script src="/static/js/bootstrap-colorpicker.min.js"></script>
        <!-- ace scripts -->
        <script src="/static/js/ace-elements.min.js"></script>
        <script src="/static/js/ace.min.js"></script>

        <!-- pjax -->
        <script src="/static/js/jquery.pjax.js"></script>

        <!-- cookie -->
        <script src="/static/js/jquery.cookie.js"></script>

        <!-- ajax -->
        <script src="/static/js/gl.senddata.js"></script>

        <!-- markdown -->
        <script src="/static/editor.md/lib/marked.min.js"></script>
        <script src="/static/editor.md/lib/prettify.min.js"></script>
        <script src="/static/editor.md/lib/underscore.min.js"></script>
        <script src="/static/editor.md/editormd.js"></script>

        <!-- get form data -->
        <script src="/static/js/jquery.formHelp.js"></script>

        <script src="/static/js/echarts.min.js"></script>

        <script type="text/javascript">
            $(function(){
              $(document).pjax('.menu-pjax-a', '#main')

              $(".menu-pjax-a").click(function(){
                    $("#platform-menu").find("li").removeClass("active");
                    $("#platform-menu").find(".hover-show").removeClass("hover-show");
                    $(this).parent("li").siblings(".parent-menu").find(".submenu").slideUp("fast");
                    $(this).parent("li").addClass("active").siblings().removeClass("active");
                    $(this).parent("li").siblings(".open").removeClass("open");
                    $(this).parents(".parent-menu").siblings(".parent-menu").find(".submenu").slideUp("fast");
                    $(this).parents(".parent-menu").addClass("active").siblings().removeClass("active");
                    $(this).parents(".open").siblings(".open").removeClass("open");
              });

              $(document).click(function(event){
                  var _con = $('#sidebar');
                  if(!_con.is(event.target) && _con.has(event.target).length === 0){ 
                        _con.removeClass("display");
                        $('#menu-toggler').removeClass("display");
                  }
              });

              $('.colorpick-btn').click(function(){
                    var c = $(this).attr('data-color');
                    if(c && c == '#438EB9'){
                        $.cookie('skin',null);
                    }
                    if(c && c == '#222A2D'){
                        $.cookie('skin','1');
                    }
                    if(c && c == '#C6487E'){
                        $.cookie('skin','2');
                    }
                    if(c && c == '#D0D0D0'){
                        $.cookie('skin','3');
                    }
              });

            })
            
            $('[data-rel=tooltip]').tooltip({container:'body'});
        </script>

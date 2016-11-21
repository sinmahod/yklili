       
        <script type="text/javascript">
            if('ontouchstart' in document.documentElement) document.write("<script src='../static/js/jquery.mobile.custom.min.js'>"+"<"+"/script>");
        </script>
        <script src="../static/js/bootstrap.min.js"></script>

        <!--[if lte IE 8]>
          <script src="../static/js/excanvas.min.js"></script>
        <![endif]-->
        
        <!-- bootstrapDialog -->
        <script src="../static/js/bootstrap-dialog.min.js"></script>
        <script src="../static/js/bootstrap-dialog.extend.js"></script>

        <!-- ace scripts -->
        <script src="../static/js/ace-elements.min.js"></script>
        <script src="../static/js/ace.min.js"></script>

        <!-- pjax -->
        <script src="../static/js/jquery.pjax.js"></script>

        <!-- ajax -->
        <script src="../static/js/jquery.senddata.js"></script>
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
            })
            
            $('[data-rel=tooltip]').tooltip({container:'body'});
        </script>

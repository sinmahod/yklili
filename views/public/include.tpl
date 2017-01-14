        <!-- bootstrap & fontawesome -->
        <link rel="stylesheet" href="/static/css/bootstrap.min.css" />
        <link rel="stylesheet" href="/static/font-awesome/4.5.0/css/font-awesome.min.css" />

        <!-- page specific plugin styles -->
        <link rel="stylesheet" href="/static/css/jquery-ui.min.css" />
        <link rel="stylesheet" href="/static/css/bootstrap-datepicker3.min.css" />
        <link rel="stylesheet" href="/static/css/ui.jqgrid.min.css" />

        <!-- text fonts -->
        <link rel="stylesheet" href="/static/css/fonts.googleapis.com.css" />

        <!-- jquery gritter -->
        <link rel="stylesheet" href="/static/css/jquery.gritter.min.css" />

        <!-- image view -->
        <link rel="stylesheet" href="/static/css/colorbox.min.css" />

        <!-- ace styles -->
        <link rel="stylesheet" href="/static/css/ace.min.css" class="ace-main-stylesheet" id="main-ace-style" />

        {{str2html "<!--[if lte IE 9]>"}}
            <link rel="stylesheet" href="/static/css/ace-part2.min.css" class="ace-main-stylesheet" />
        {{str2html "<![endif]-->"}}
        <link rel="stylesheet" href="/static/css/ace-skins.min.css" />
        <link rel="stylesheet" href="/static/css/ace-rtl.min.css" />

        <!-- bootstrapDialog -->
        <link rel="stylesheet" href="/static/css/bootstrap-dialog.min.css" />

        <!-- webUploader -->
        <link rel="stylesheet" href="/static/webuploader/webuploader.css" />

        {{str2html "<!--[if lte IE 9]>"}}
          <link rel="stylesheet" href="/static/css/ace-ie.min.css" />
        {{str2html "<![endif]-->"}}

        <!-- gl style -->
        <link rel="stylesheet" href="/static/css/gl_style.css" />        

        <!-- inline styles related to this page -->

        <!-- ace settings handler -->
        <script src="/static/js/ace-extra.min.js"></script>

        {{str2html "<!-- HTML5shiv and Respond.js for IE8 to support HTML5 elements and media queries -->"}}

        {{str2html "<!--[if lte IE 8]>"}}
        <script src="/static/js/html5shiv.min.js"></script>
        <script src="/static/js/respond.min.js"></script>
        {{str2html "<<![endif]-->"}}

        {{str2html "<!--[if !IE]> -->"}}
        <script src="/static/js/jquery-2.1.4.min.js"></script>
        {{str2html "<!-- <![endif]-->"}}
        <!-- <![
        {{str2html "<!--[if IE]>"}}
        <script src="/static/js/jquery-1.11.3.min.js"></script>
        {{str2html "<![endif]-->"}}

        <!-- page specific plugin scripts -->
        <script src="/static/js/bootstrap-datepicker.min.js"></script>
        <script src="/static/js/jquery.jqGrid.min.js"></script>
        <script src="/static/js/grid.locale-en.js"></script>
        <script src="/static/js/gl.jqGrid.extend.js" ></script>
        <script type="text/javascript">
            //配置jqGrid交替颜色
            jQuery.extend(jQuery.jgrid.defaults, { altRows:true });
        </script>


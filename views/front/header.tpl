<div class="header" style="background-image: url('{{"Banner" | SiteConfig}}')">
    <div class="top-header">
        <div class="navigation">
            <div class="top-menu">
            <span class="menu"><img src="static/front/images/menu.png" alt=""> </span>
                <ul id="list2" dataurl="/catalog/List?Status=1" page="false">
                </ul>
            </div>
            <!--script for menu-->
            <script>
                $("span.menu").click(function(){
                    $(".top-menu ul").slideToggle(500, function(){
                    });
                });

            </script>
            <!--script for menu-->
        </div>
        <div class="search">
             <!-- start search-->
            <div class="search-box">
                <div id="sb-search" class="sb-search">
                    <form action="/search">
                        <input class="sb-search-input" placeholder="请输入您要查询的内容" type="search" name="q" id="q">
                        <input class="sb-search-submit" type="submit" value="">
                        <span class="sb-icon-search"> </span>
                    </form>
                </div>
            </div>
            <!-- search-scripts -->
            <script src="static/front/js/classie.js"></script>
            <script src="static/front/js/uisearch.js"></script>
                <script>
                    new UISearch( document.getElementById( 'sb-search' ) );
                </script>
            <!-- //search-scripts -->
        </div>
        <div class="clearfix"></div>
    </div>
    <div class="logo logotext text-center">
        <h3>{{"Subtitle" | SiteConfig}}</h3>
        <h1>{{"Name" | SiteConfig}}</h1>
    </div>
</div>
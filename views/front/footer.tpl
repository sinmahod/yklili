    <div class="footer">
        <div class="footer-top">
            <h2>{{"Name" | SiteConfig}}</h2> 
            <p>{{"Desc" | SiteConfig}}</p>
            <div class="clearfix"></div>
        </div>
        <div class="footer-bottom">
            <div class="bottom-menu">
                 <ul id="list3" dataurl="/catalog/List?Status=1" page="false">
                </ul>
            </div>
            <div class="copy-rights">
                <p>{{"Copyright" | SiteConfig}} | Proudly published with  <a href="https://yklili.com" target="target_blank">YKlili</a></p>
            </div>
            <div class="clearfix"></div>
        </div>
    </div>
    <script type="text/javascript">
        jQuery(function($) {
            DataList({
                    list: "list3",
                    size: 10,
                    html: '<li><a href="${Link}">${CatalogName}</a></li>'
            });
        });
    </script>

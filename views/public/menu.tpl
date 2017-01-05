            <div id="sidebar" class="sidebar responsive ace-save-state">
                <script type="text/javascript">
                    try{ace.settings.loadState('sidebar')}catch(e){}
                </script>

                <div class="sidebar-shortcuts" id="sidebar-shortcuts">
                    <div class="sidebar-shortcuts-large" id="sidebar-shortcuts-large">
                        <button class="btn btn-success">
                            <i class="ace-icon fa fa-signal"></i>
                        </button>

                        <button class="btn btn-info">
                            <i class="ace-icon fa fa-pencil"></i>
                        </button>

                        <button class="btn btn-warning">
                            <i class="ace-icon fa fa-users"></i>
                        </button>

                        <button class="btn btn-danger">
                            <i class="ace-icon fa fa-cogs"></i>
                        </button>
                    </div>     

                    <div class="sidebar-shortcuts-mini" id="sidebar-shortcuts-mini">
                        <span class="btn btn-success"></span>

                        <span class="btn btn-info"></span>

                        <span class="btn btn-warning"></span>

                        <span class="btn btn-danger"></span>
                    </div>
                </div><!-- /.sidebar-shortcuts -->

                <ul class="nav nav-list" id="platform-menu">
                    {{range .Menus}}
                    <li class="{{if .Checked}}active{{if .ChildNode}} open{{end}}{{end}} parent-menu">
                        <a href="{{.Link}}" class="{{if .ChildNode}}dropdown-toggle {{else}}menu-pjax-a{{end}}">
                            <i class="menu-icon fa {{.Icon}}"></i>
                            <span class="menu-text"> 
                                {{.MenuName}}
                            </span>
                            <b class="arrow{{if .ChildNode}} fa fa-angle-down{{end}}"></b>
                        </a>
                        {{if .ChildNode}}
                        <b class="arrow"></b>
                        <ul class="submenu">
                            {{range .ChildNode}}
                            <li class="{{if .Checked}}active{{end}}">
                                <a href="{{.Link}}" class="menu-pjax-a">
                                    <i class="menu-icon fa fa-caret-right"></i>
                                    {{.MenuName}}
                                </a>
                                <b class="arrow"></b>
                            </li>
                            {{end}} 
                        </ul>
                        {{end}}
                    </li>
                    {{end}} 
                    
                </ul><!-- /.nav-list -->

                <div class="sidebar-toggle sidebar-collapse" id="sidebar-collapse">
                    <i id="sidebar-toggle-icon" class="ace-icon fa fa-angle-double-left ace-save-state" data-icon1="ace-icon fa fa-angle-double-left" data-icon2="ace-icon fa fa-angle-double-right"></i>
                </div>
            </div>
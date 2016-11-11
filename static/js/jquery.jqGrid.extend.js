(function ($) {
/**
 * jqGrid 封装类
 */
 window.DataGrid = function(options){
	//url 数据请求连接  必选
	//rows每页展示/请求数量  可选，默认10000
	//pkcolumn 主键字段
	//istree  是否树形展示
	//addsave  添加数据或保存数据的url
	//<table id="dataGrid" url="test" rows="10" pkcolumn="Id" istree="true"></table>
	var dataGrid =  $('#'+options.tableName);
	var dataGridPage = $('#'+options.pageName);

	if (!dataGrid){
		console.log('Error:没有设置tableName');
		return;
	}
	if (!dataGridPage){
		console.log('Error:没有设置pageName');
		return;
	}
	var rows = dataGrid.attr('rows') ? dataGrid.attr('rows') : parseInt(($(window).height() - 320)/39);

	var parent_column = dataGrid.closest('[class*="col-"]');
                //resize to fit page size
                $(window).on('resize.jqGrid', function () {
                    dataGrid.jqGrid( 'setGridWidth', parent_column.width() );
                })
                
                //resize on sidebar collapse/expand
                $(document).on('settings.ace.jqGrid' , function(ev, event_name, collapsed) {
                    if( event_name === 'sidebar_collapsed' || event_name === 'main_container_fixed' ) {
                        //setTimeout is for webkit only to give time for DOM changes and then redraw!!!
                        setTimeout(function() {
                            dataGrid.jqGrid( 'setGridWidth', parent_column.width() );
                        }, 20);
                    }
                })

	options = $.extend({
		url: dataGrid.attr('url'),
		mtype: "GET", 
		datatype: "json",

		
		caption: options.title,  	//表格名称
		rowNum: rows, 		// 显示行数，-1为显示全部,默认自适应
		rowList: [rows,30,50],    //每页最大数量支持
		rownumWidth: 30,	// 序号列宽
		multiboxonly: true,	// 单击复选框时在多选
		altRows: true, 		// 斑马线样式，交替行altclass
		
		// 设置列模型
		columnModel: [], 
		colNames: [], 
		colModel: options.columnModel,

		// 列表参数
		dataId: dataGrid.attr('pkcolumn') ? dataGrid.attr('pkcolumn') : 'id', 		// 指定数据主键
		showRownum: true,	// 是否显示行号
		showCheckbox: true,// 是否显示复选框
		multiselect: true,      //是否可以多选
    		multiboxonly: true,     //只有选中checkbox才有效
		sortable: false,	// 列表是否允许支持
		
		// 树结构表格
		treeGrid:  dataGrid.attr('istree') && dataGrid.attr('istree')=='true'? true : false, 	// 启用树结构表格
		treeGridModel: 'adjacency',					// 启用简单结构模式
		ExpandColClick: true,						// 单击列可展开
		defaultExpandLevel: 0,						// 默认展开的层次
		initExpandLevel: options.defaultExpandLevel,			// 保存初始化是设置的展开层次
		treeReader: {	// 自定义树表格JSON读取参数
			level_field: "Level",  
			parent_id_field: "Pid",  
			leaf_field: "IsLeaf",  
			expanded_field: "Expanded" 
		},
		ExpandColumn: options.treeColumn,	//属性结构列的列明

		//设置宽度为0px，不显示滚动条
		scrollOffset: 0,  
		viewrecords : true, //是否显示总数量
		height: $(window).height() - 320,     //表格高度，宽度默认自动填充

		pager : dataGridPage,     //翻页导航

		closeAfterAdd: true,    //添加数据后关闭窗口
		closeAfterEdit:true,     //修改数据后关闭窗口
		//页码的文字
		pgtext: '转到 <input class="ui-pg-input ui-corner-all" type="text" size="2" maxlength="7" value="0" role="textbox"> 页，共<span id="sp_1_grid-pager"></span>页',
		loadComplete : function() {
                       		var table = this;
                        		setTimeout(function(){
                            		styleCheckbox(table);
		                          updateActionIcons(table);
                            		updatePagerIcons(table);
                            		enableTooltips(table);
                        		}, 0);
                    	},
                    	editurl: dataGrid.attr('addsave'), //添加修改数据请求的url
                        
      	}, options);
	
	// 获取列标题
	options.colNames = [];
	for(var i=0; i<options.colModel.length; i++){
		options.colNames.push(options.colModel[i].header);
		// 如果是树结构表格
		if (options.treeGrid || !options.sortable){
			options.colModel[i].sortable = false; // 是否排序列
		}
	}
	
	// 如果是树结构表格
	if (options.treeGrid){
		options.showRownum = false;	// 是否显示行号
		options.showCheckbox = false; // 是否显示复选框
	}
	
	// 显示序号
	if (options.showRownum){
		options.rownumbers = true;	// 显示序号
	}
	
	// 显示多选复选框
	if (options.showCheckbox){
		options.multiselect = true;	// 显示多选复选框
	}
	
	// 如果设置了多级表头或冻结列
	if (options.groupHeaders || options.frozenCols){
		options.shrinkToFit = false;// 不按百分比自适应列宽
	}else{
		options.shrinkToFit = true;	// 按百分比自适应列宽
	}
	
	// 是否显示合计行
	if (options.showFooter){
		options.footerrow = true; 		// 显示底部合计行
		options.userDataOnFooter = true; // 使用json数据作为合计行数据
	}
	
	// 初始化jqGrid
	dataGrid.jqGrid(options);
	
	// 是否冻结列
	if (options.frozenCols){
		dataGrid.jqGrid('setFrozenColumns'); // 冻结列，在colModel指定frozen: true
	}	
	
    //postdata=提交的数据
    var fn_editSubmit=function(response,postdata){ 
        var json=response.responseText; 
        alert(json);//显示返回值 
        $('#cData').trigger('click');//执行关闭按钮的点击事件
    } 

	// 自动调整表格大小
	$(window).triggerHandler('resize.jqGrid');
	dataGrid.jqGrid('navGrid','#'+dataGridPage.attr('id'),
            	{   //navbar options
                        edit: true,
                        editicon : 'ace-icon fa fa-pencil blue',
                        add: true,
                        addicon : 'ace-icon fa fa-plus-circle purple',
                        del: true,
                        delicon : 'ace-icon fa fa-trash-o red',
                        search: true,
                        searchicon : 'ace-icon fa fa-search orange',
                        refresh: true,
                        refreshicon : 'ace-icon fa fa-refresh green',
                        view: true,
                        viewicon : 'ace-icon fa fa-search-plus grey',
             },
             {
                        //edit record form
                        //closeAfterEdit: true,
                        //width: 700,
                        closeOnEscape: true,    //开启ESC关闭对话框功能
                        afterSubmit: fn_editSubmit, //提交后执行的函数

                        bSubmit: "保存",
                        bCancel: "关闭",
                       // closeAfterAdd: true,    //添加数据后关闭窗口
                       // closeAfterEdit:true,     //修改数据后关闭窗口
                        recreateForm: true,
                        beforeShowForm : function(e) {
                            var form = $(e[0]);
                            form.closest('.ui-jqdialog').find('.ui-jqdialog-titlebar').wrapInner('<div class="widget-header" />')
                            style_edit_form(form);
                        }
             },
             {
                        //new record form
                        closeAfterAdd: true,
                        recreateForm: true,
                        viewPagerButtons: false,
                        beforeShowForm : function(e) {
                            var form = $(e[0]);
                            form.closest('.ui-jqdialog').find('.ui-jqdialog-titlebar')
                            wrapInner('<div class="widget-header" />')
                            style_edit_form(form);
                        }
             },
             {
                        //delete record form
                        recreateForm: true,
                        beforeShowForm : function(e) {
                            var form = $(e[0]);
                            if(form.data('styled')) return false;
                            
                            form.closest('.ui-jqdialog').find('.ui-jqdialog-titlebar').wrapInner('<div class="widget-header" />')
                            style_delete_form(form);
                            
                            form.data('styled', true);
                        },
                        onClick : function(e) {
                        }
             },
             {
                        recreateForm: true,
                        afterShowSearch: function(e){
                            var form = $(e[0]);
                            form.closest('.ui-jqdialog').find('.ui-jqdialog-title').wrap('<div class="widget-header" />')
                            style_search_form(form);
                        },
                        afterRedraw: function(){
                            style_search_filters($(this));
                        }
                        ,
                        multipleSearch: true,
             },
             {
                        recreateForm: true,
                        beforeShowForm: function(e){
                            var form = $(e[0]);
                            form.closest('.ui-jqdialog').find('.ui-jqdialog-title').wrap('<div class="widget-header" />')
                        }
             });

	        $(document).one('ajaxloadstart.page', function(e) {
                    $.jgrid.gridDestroy(grid_selector);
                    $('.ui-jqdialog').remove();
             });
};


//封装常用方法

/**
*   得到当前选中行ID，返回数值型数组
*   @Param  tableid  table的ID
*   @Return  ids[] 选中行字段值
*/
DataGrid.getSelectRowIds = function(tableid){
	var tableGrid =  $('#'+tableid);
	var pkcolumn = tableGrid.attr('pkcolumn');
	var rowIds = tableGrid.jqGrid('getGridParam','selarrrow');
	if (rowIds.length == 0){
		var rowIds = new Array();
		 id =  tableGrid.jqGrid('getGridParam','selrow');
		 if (id!=null){
			 rowIds[0] = id	
		 }
	}
	var ids = new Array();
	for (var i = 0 ; i < rowIds.length ; i ++){
		ids[i] = $("#datatable").jqGrid('getCell',rowIds[i],pkcolumn);	
	}
	return  ids;
}



//文本替换
	/**
	 * jqGrid Chinese Translation
	 * 咖啡兔 yanhonglei@gmail.com
	 * http://www.kafeitu.me 
	 * Dual licensed under the MIT and GPL licenses:
	 * http://www.opensource.org/licenses/mit-license.php
	 * http://www.gnu.org/licenses/gpl.html
	**/
	$.jgrid = $.jgrid || {};
	$.extend($.jgrid,{
	    defaults : {
	        recordtext: "{0} - {1}\u3000共 {2} 条", // 共字前是全角空格
	        emptyrecords: "无数据显示",
	        loadtext: "正在加载...",
	        pgtext : " {0} 共 {1} 页"
	    },
	    search : {
	        caption: "搜索...",
	        Find: "查找",
	        Reset: "重置",
	        odata: [{ oper:'eq', text:'等于\u3000\u3000'},{ oper:'ne', text:'不等\u3000\u3000'},{ oper:'lt', text:'小于\u3000\u3000'},{ oper:'le', text:'小于等于'},{ oper:'gt', text:'大于\u3000\u3000'},{ oper:'ge', text:'大于等于'},{ oper:'bw', text:'开始于'},{ oper:'bn', text:'不开始于'},{ oper:'in', text:'属于\u3000\u3000'},{ oper:'ni', text:'不属于'},{ oper:'ew', text:'结束于'},{ oper:'en', text:'不结束于'},{ oper:'cn', text:'包含\u3000\u3000'},{ oper:'nc', text:'不包含'},{ oper:'nu', text:'不存在'},{ oper:'nn', text:'存在'}],
	        groupOps: [ { op: "AND", text: "所有" },    { op: "OR",  text: "任一" } ],
			operandTitle : "Click to select search operation.",
			resetTitle : "Reset Search Value"
	    },
	    edit : {
	        addCaption: "添加记录",
	        editCaption: "编辑记录",
	        bSubmit: "提交",
	        bCancel: "取消",
	        bClose: "关闭",
	        saveData: "数据已改变，是否保存？",
	        bYes : "是",
	        bNo : "否",
	        bExit : "取消",
	        msg: {
	            required:"此字段必需",
	            number:"请输入有效数字",
	            minValue:"输值必须大于等于 ",
	            maxValue:"输值必须小于等于 ",
	            email: "这不是有效的e-mail地址",
	            integer: "请输入有效整数",
	            date: "请输入有效时间",
	            url: "无效网址。前缀必须为 ('http://' 或 'https://')",
	            nodefined : " 未定义！",
	            novalue : " 需要返回值！",
	            customarray : "自定义函数需要返回数组！",
	            customfcheck : "必须有自定义函数!"
	        }
	    },
	    view : {
	        caption: "查看记录",
	        bClose: "关闭"
	    },
	    del : {
	        caption: "删除",
	        msg: "删除所选记录？",
	        bSubmit: "删除",
	        bCancel: "取消"
	    },
	    nav : {
	        edittext: "",
	        edittitle: "编辑所选记录",
	        addtext:"",
	        addtitle: "添加新记录",
	        deltext: "",
	        deltitle: "删除所选记录",
	        searchtext: "",
	        searchtitle: "查找",
	        refreshtext: "",
	        refreshtitle: "刷新表格",
	        alertcap: "注意",
	        alerttext: "请选择记录",
	        viewtext: "",
	        viewtitle: "查看所选记录"
	    },
	    col : {
	        caption: "选择列",
	        bSubmit: "确定",
	        bCancel: "取消"
	    },
	    errors : {
	        errcap : "错误",
	        nourl : "没有设置url",
	        norecords: "没有要处理的记录",
	        model : "colNames 和 colModel 长度不等！"
	    },
	    formatter : {
	        integer : {thousandsSeparator: ",", defaultValue: '0'},
	        number : {decimalSeparator:".", thousandsSeparator: ",", decimalPlaces: 2, defaultValue: '0.00'},
	        currency : {decimalSeparator:".", thousandsSeparator: ",", decimalPlaces: 2, prefix: "", suffix:"", defaultValue: '0.00'},
	        date : {
	            dayNames:   [
	                "日", "一", "二", "三", "四", "五", "六",
	                "星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六",
	            ],
	            monthNames: [
	                "一", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二",
	                "一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"
	            ],
	            AmPm : ["am","pm","上午","下午"],
	            S: function (j) {return j < 11 || j > 13 ? ['st', 'nd', 'rd', 'th'][Math.min((j - 1) % 10, 3)] : 'th';},
	            srcformat: 'Y-m-d',
	            newformat: 'Y-m-d',
	            parseRe : /[#%\\\/:_;.,\t\s-]/,
	            masks : {
	                // see http://php.net/manual/en/function.date.php for PHP format used in jqGrid
	                // and see http://docs.jquery.com/UI/Datepicker/formatDate
	                // and https://github.com/jquery/globalize#dates for alternative formats used frequently
	                // one can find on https://github.com/jquery/globalize/tree/master/lib/cultures many
	                // information about date, time, numbers and currency formats used in different countries
	                // one should just convert the information in PHP format
	                ISO8601Long:"Y-m-d H:i:s",
	                ISO8601Short:"Y-m-d",
	                // short date:
	                //    n - Numeric representation of a month, without leading zeros
	                //    j - Day of the month without leading zeros
	                //    Y - A full numeric representation of a year, 4 digits
	                // example: 3/1/2012 which means 1 March 2012
	                ShortDate: "n/j/Y", // in jQuery UI Datepicker: "M/d/yyyy"
	                // long date:
	                //    l - A full textual representation of the day of the week
	                //    F - A full textual representation of a month
	                //    d - Day of the month, 2 digits with leading zeros
	                //    Y - A full numeric representation of a year, 4 digits
	                LongDate: "l, F d, Y", // in jQuery UI Datepicker: "dddd, MMMM dd, yyyy"
	                // long date with long time:
	                //    l - A full textual representation of the day of the week
	                //    F - A full textual representation of a month
	                //    d - Day of the month, 2 digits with leading zeros
	                //    Y - A full numeric representation of a year, 4 digits
	                //    g - 12-hour format of an hour without leading zeros
	                //    i - Minutes with leading zeros
	                //    s - Seconds, with leading zeros
	                //    A - Uppercase Ante meridiem and Post meridiem (AM or PM)
	                FullDateTime: "l, F d, Y g:i:s A", // in jQuery UI Datepicker: "dddd, MMMM dd, yyyy h:mm:ss tt"
	                // month day:
	                //    F - A full textual representation of a month
	                //    d - Day of the month, 2 digits with leading zeros
	                MonthDay: "F d", // in jQuery UI Datepicker: "MMMM dd"
	                // short time (without seconds)
	                //    g - 12-hour format of an hour without leading zeros
	                //    i - Minutes with leading zeros
	                //    A - Uppercase Ante meridiem and Post meridiem (AM or PM)
	                ShortTime: "g:i A", // in jQuery UI Datepicker: "h:mm tt"
	                // long time (with seconds)
	                //    g - 12-hour format of an hour without leading zeros
	                //    i - Minutes with leading zeros
	                //    s - Seconds, with leading zeros
	                //    A - Uppercase Ante meridiem and Post meridiem (AM or PM)
	                LongTime: "g:i:s A", // in jQuery UI Datepicker: "h:mm:ss tt"
	                SortableDateTime: "Y-m-d\\TH:i:s",
	                UniversalSortableDateTime: "Y-m-d H:i:sO",
	                // month with year
	                //    Y - A full numeric representation of a year, 4 digits
	                //    F - A full textual representation of a month
	                YearMonth: "F, Y" // in jQuery UI Datepicker: "MMMM, yyyy"
	            },
	            reformatAfterEdit : false
	        },
	        baseLinkUrl: '',
	        showAction: '',
	        target: '',
	        checkbox : {disabled:true},
	        idName : 'id'
	    }
	});
})(jQuery);

/**  ace funcation  **/
//switch element when editing inline
function aceSwitch( cellvalue, options, cell ) {
    setTimeout(function(){
        $(cell) .find('input[type=checkbox]')
            .addClass('ace ace-switch ace-switch-5')
            .after('<span class="lbl"></span>');
    }, 0);
}
//enable datepicker
function pickDate( cellvalue, options, cell ) {
    setTimeout(function(){
        $(cell) .find('input[type=text]')
            .datepicker({format:'yyyy-mm-dd' , autoclose:true}); 
    }, 0);
}

//navButtons
function style_edit_form(form) {
    form.find('input[name=sdate]').datepicker({format:'yyyy-mm-dd' , autoclose:true})
    
    form.find('input[name=stock]').addClass('ace ace-switch ace-switch-5').after('<span class="lbl"></span>');
    var buttons = form.next().find('.EditButton .fm-button');
    buttons.addClass('btn btn-sm').find('[class*="-icon"]').hide();//ui-icon, s-icon
    buttons.eq(0).addClass('btn-primary').prepend('<i class="ace-icon fa fa-check"></i>');
    buttons.eq(1).prepend('<i class="ace-icon fa fa-times"></i>')
    
    buttons = form.next().find('.navButton a');
    buttons.find('.ui-icon').hide();
    buttons.eq(0).append('<i class="ace-icon fa fa-chevron-left"></i>');
    buttons.eq(1).append('<i class="ace-icon fa fa-chevron-right"></i>');       
}

function style_delete_form(form) {
    var buttons = form.next().find('.EditButton .fm-button');
    buttons.addClass('btn btn-sm btn-white btn-round').find('[class*="-icon"]').hide();//ui-icon, s-icon
    buttons.eq(0).addClass('btn-danger').prepend('<i class="ace-icon fa fa-trash-o"></i>');
    buttons.eq(1).addClass('btn-default').prepend('<i class="ace-icon fa fa-times"></i>')
}

function style_search_filters(form) {
    form.find('.delete-rule').val('X');
    form.find('.add-rule').addClass('btn btn-xs btn-primary');
    form.find('.add-group').addClass('btn btn-xs btn-success');
    form.find('.delete-group').addClass('btn btn-xs btn-danger');
}
function style_search_form(form) {
    var dialog = form.closest('.ui-jqdialog');
    var buttons = dialog.find('.EditTable')
    buttons.find('.EditButton a[id*="_reset"]').addClass('btn btn-sm btn-info').find('.ui-icon').attr('class', 'ace-icon fa fa-retweet');
    buttons.find('.EditButton a[id*="_query"]').addClass('btn btn-sm btn-inverse').find('.ui-icon').attr('class', 'ace-icon fa fa-comment-o');
    buttons.find('.EditButton a[id*="_search"]').addClass('btn btn-sm btn-purple').find('.ui-icon').attr('class', 'ace-icon fa fa-search');
}

function beforeDeleteCallback(e) {
    var form = $(e[0]);
    if(form.data('styled')) return false;
    
    form.closest('.ui-jqdialog').find('.ui-jqdialog-titlebar').wrapInner('<div class="widget-header" />')
    style_delete_form(form);
    
    form.data('styled', true);
}

function beforeEditCallback(e) {
    var form = $(e[0]);
    form.closest('.ui-jqdialog').find('.ui-jqdialog-titlebar').wrapInner('<div class="widget-header" />')
    style_edit_form(form);
}

function styleCheckbox(table) {           
}


function updateActionIcons(table) {
}

function updatePagerIcons(table) {
    var replacement = {
        'ui-icon-seek-first' : 'ace-icon fa fa-angle-double-left bigger-140',
        'ui-icon-seek-prev' : 'ace-icon fa fa-angle-left bigger-140',
        'ui-icon-seek-next' : 'ace-icon fa fa-angle-right bigger-140',
        'ui-icon-seek-end' : 'ace-icon fa fa-angle-double-right bigger-140'
    };
    $('.ui-pg-table:not(.navtable) > tbody > tr > .ui-pg-button > .ui-icon').each(function(){
        var icon = $(this);
        var $class = $.trim(icon.attr('class').replace('ui-icon', ''));
        
        if($class in replacement) icon.attr('class', 'ui-icon '+replacement[$class]);
    })
}

function enableTooltips(table) {
    $('.navtable .ui-pg-button').tooltip({container:'body'});
    $(table).find('.ui-pg-div').tooltip({container:'body'});
}
            

                
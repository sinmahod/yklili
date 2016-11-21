(function ($) {
 /**
 * BootstrapDialog 封装类
 * author: gl
 * time: 2016/11/18 22:33:07
 * email: sinmahod@qq.com
 */
    window.BootFrame = function () {
	return {
		//弹出框
		alert: function (msg,fn,tle) {
	        	BootstrapDialog.alert({
		            title: tle ? tle : '提示',
		            message: msg,
		            type: BootstrapDialog.TYPE_PRIMARY, // <-- Default value is BootstrapDialog.TYPE_WARNING
		            closable: true, // <-- Default value is false
		            draggable: true, // <-- Default value is false
		            buttonLabel: '确定', // <-- Default value is 'OK',
		            callback: function(result) {
		               if (fn){
		               		fn();
		               }
		            }
		        });
		},
		//选择框
		confirm: function (message,truefn,falsefn){
			BootstrapDialog.confirm(message, function(result){
            				if(result) {
			             	if(truefn){
			             		truefn();
			             	}
			            	}else {
			             	if(falsefn){
			             		falsefn();
			             	}
			            }
			});
		},
		//Dialog
		dialog: function(){
			var t,m;
			var b = {};
			var hideclose = false;
			var w = 600;
			var h = 150;
			var diaid;
			var id;
			return{
				id : function(id){
					diaid = id;
				},
				title : function(title){
					t = title;
				},
				message :function(message){
					m = message;
				},
				url: function(u){
					m = $('<div></div>').load(u);
				},
				width:function(width){
					w = width / 5 * 5;
				},
				height:function(height){
					h = height / 5 * 5 ; //框头框尾130px
				},
				addButton: function(name,fn,css){
					if(!css){
						css = 'btn-primary';
					}
					if (jQuery.isEmptyObject(b) ) {
						b = [{
							label:name,
							cssClass:css,
							action:function(){
								fn();
							}
						}]
					 }else{
						var bt = {
							label:name,
							cssClass:css,
							action:function(){
								fn();
							}
						}
						b.push(bt)
					}
				},
				hideClose: function(){
					hideclose = true;
				},
				show:function(){
					if (!hideclose){
						var bt = {
							label: '取消',
		                    	action: function(dialogItself){
		                        		dialogItself.close();
		                    	}
						}
						b.push(bt)	
					}
				  	var dialog = BootstrapDialog.show({
				  					id: diaid,
					                title: t,
					                message: m,
					                cssClass: 'dialog-'+w+' dialog-h-'+h,
					                draggable: true,
					                buttons: b
				            });
				  	id = dialog.$modal[0].id;
				},
				getFormData: function(){
					var s = $("#"+id + " form");
					if (s.length == 0) {
						BootFrame.alert("未找到form标签");
						return;
					}
					if (s[0].length == 0 ){
						BootFrame.alert("form中不存在表单");
						return;
					}
					var a = s[0];
					var tempArr = [];
					tempArr.push('{');
					for (var i = 0 ; i < a.length; i ++){
						if (i != 0){
							tempArr.push(',');
						}
						tempArr.push('"'+a[i].id + '":"'+a[i].value+'"');
					}
					tempArr.push('}');
					return $.parseJSON(tempArr.join(''));
				}
			}//<! return >
		}//<! dialog >
	}//<! return >
   }();

})(jQuery);
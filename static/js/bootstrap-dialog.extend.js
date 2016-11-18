(function ($) {
 /**
 * BootstrapDialog 封装类
 */
    window.BootFrame = function () {
	return {
		//弹出框
		alert: function (message,fn) {
		      	BootstrapDialog.alert(message, function(){
		      		if (fn){
		      			fn();
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
			return{
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
				  	BootstrapDialog.show({
					                title: t,
					                message: m,
					                cssClass: 'dialog-'+w+' dialog-h-'+h,
					                draggable: true,
					                buttons: b
				            });
				}
			}//<! return >
		}//<! dialog >
	}//<! return >
   }();

})(jQuery);
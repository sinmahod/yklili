(function ($) {
 /**
 * BootstrapDialog 封装类
 * author: gl
 * time: 2016/11/18 22:33:07
 * email: sinmahod@qq.com
 */
    window.BootFrame = function (e) {
	return {
		//弹出框
		alert: function (msg,fn,tle,iswarning,h,w) {
			h = h ? h : 60;
			w = w ? w : 350;
	        		var dobj = BootstrapDialog.alert({
		            		title: tle ? tle : '提示',
		            		message: msg,
		            		type: iswarning ? BootstrapDialog.TYPE_WARNING : BootstrapDialog.TYPE_PRIMARY, // <-- Default value is BootstrapDialog.TYPE_WARNING
		            		cssClass: 'dialog-'+w+' dialog-h-'+h,
		            		closable: true, // <-- Default value is false
		            		draggable: true, // <-- Default value is false
		            		buttonLabel: '确定', // <-- Default value is 'OK',
		            		callback: function(result) {
	                 			if (fn && fn instanceof Function) {
		               			fn();
	               			}
		           		}
		        	});
		        	dobj.getButton(dobj.getButtons()[0].id).addClass('btn-sm');
		        	dobj.getModalHeader().css('padding','10px 10px 10px 15px');
		  	dobj.getModalFooter().css('padding','10px 15px');
		},
		//选择框
		confirm: function (message,fn,falsefn){
			var dobj = BootstrapDialog.confirm({
			            title: '确认操作',
			            message: message,
			            type: BootstrapDialog.TYPE_SUCCESS, // <-- Default value is BootstrapDialog.TYPE_PRIMARY
			            closable: true, // <-- Default value is false
			            draggable: true, // <-- Default value is false
			            btnCancelLabel: '取消', // <-- Default value is 'Cancel',
			            btnOKLabel: '确定', // <-- Default value is 'OK',
			            btnOKClass: 'btn-success', // <-- If you didn't specify it, dialog type will be used,
			            callback: function(result) {
					if(result) {
					   	if (fn && fn instanceof Function) {
			               			fn();
		               			}
					}else {
					    	if (falsefn && falsefn instanceof Function) {
			               			falsefn();
		               			}
					}
			             }
		        	});
			dobj.getButton(dobj.getButtons()[0].id).addClass('btn-sm');
			dobj.getButton(dobj.getButtons()[1].id).addClass('btn-sm');
			dobj.getModalHeader().css('padding','10px 10px 10px 15px');
	  		dobj.getModalFooter().css('padding','10px 15px');
		},
		//Dialog
		dialog: function(){
			var t,m;
			var b = [];
			var hideclose = false;
			var w = 600;
			var h = 150;
			var diaid;
			var id;
			var dobj;
			var p;
			return{
				id : function(id){
					diaid = id;
				},
				title : function(title){
					t = title;
				},
				type : function(type){
					p = type;
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
				addButton: function(name,fn,css,keycode){
					var bf = this;  //获得dialog实体

					if(!css){
						css = 'btn-primary';
					}
					if (jQuery.isEmptyObject(b) ) {
						b = [{
							label:name,
							cssClass:css + ' btn-sm',
							hotkey:keycode?keycode:0,
							action:function(){
								if (fn && fn instanceof Function) {
						               		fn(bf,this);  //这里的this指的是按钮
					               		}
							}
						}]
					 }else{
						var bt = {
							label:name,
							cssClass:css + ' btn-sm',
							hotkey:keycode?keycode:0,
							action:function(){
								if (fn && fn instanceof Function) {
						               		fn(bf,this);  //这里的this指的是按钮
					               		}
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
							cssClass: 'btn-success btn-sm',
					                    	action: function(dialogItself){
					                        		dialogItself.close();
					                    	}
						}
						b.push(bt);
					}
				  	dobj = BootstrapDialog.show({
			  			   id: diaid,
					                title: t,
					                message: m,
					                cssClass: 'dialog-'+w+' dialog-h-'+h,
					                type: p,
					                closeByBackdrop: false,   //点击空白位置关闭窗口失效
					                draggable: true,
					                buttons: b,
					                onshown:function(){  //增加verify属性检查(注册校验)
					                	BValidate({dialogId:this.id});
					                }
				           });
				  	dobj.getModalHeader().css('padding','10px 10px 10px 15px');
				  	dobj.getModalFooter().css('padding','10px 15px');
				  	id = dobj.$modal[0].id;
				},
				verifyForm: function(warn){
					var s = $("#"+id + " form");
					if (s.length == 0) {
						BootFrame.alert("未找到form标签",null,"警告",true);
						return false;
					}
					if (s[0].length == 0 ){
						BootFrame.alert("form中不存在表单",null,"警告",true);
						return false;
					}
					s.data('bootstrapValidator').validate();
					if (!s.data('bootstrapValidator').isValid()){
						if(warn){
							BootFrame.alert("请按照规则输入表单",null,"校验错误",true);
						}
						return false;
					}
					return true;
				},
				getFormData: function(){
					var s = $("#"+id + " form");
					if (s.length == 0) {
						BootFrame.alert("未找到form标签",null,"警告",true);
						return;
					}
					if (s[0].length == 0 ){
						BootFrame.alert("form中不存在表单",null,"警告",true);
						return;
					}
					var a = s[0];
					var $json = {};
					for (var i = 0 ; i < a.length; i ++){
						//判断表单类型
						var elementid = a[i].id;
						if (elementid == "") {
							elementid = a[i].name;
						}
						if (a[i].type == "radio"){ 
							if(a[i].checked == true){
								$json[elementid] = a[i].value;
							}
						}else{
							$json[elementid] = a[i].value;
						}
					}
					return $json;
				},
				close: function(){
					dobj.close();
				}
			}//<! return >
		},//<! dialog >
		//Dialog
		progressbar: function(taskid){
	 		var dialog = new BootstrapDialog({
		 		message: function(dialogRef){
		                		var $message = $('<div id="progressbar"><div class="progress-label">加载中...</div></div>');
		                		return $message;
		            		},
		            		title:"任务进行中，请耐心等待",
		            		draggable: true,
		            		closable: true,
		            		closeByBackdrop: false, 
		            		cssClass: 'dialog-450 dialog-h-60',
		            		buttons: [{
					label: '确定',
					cssClass: 'btn-primary btn-sm',
			                    	action: function(dialogItself){
			                        		dialogItself.close();
			                    	}
				}],
		            		onshown:function(dialogRef){
		            			$( "#progressbar" ).progressbar({
						value: 0,
						create: function( event, ui ) {
							var ele = $(this);
							ele.addClass('progress progress-striped active');
							$(ele.children(0)[1]).addClass('progress-bar progress-bar-success');
						},
						change: function() {
							$( ".progress-label" ).text( $( "#progressbar" ).progressbar( "value" ) + "%" );
					     	},
				      		complete: function() {
					        		$( ".progress-label" ).text( "完成！" );
					      	}
					});
					var interval = setInterval(function(){
						var data = {"taskId" : taskid}
			            			$.post("/platform/prog",data,function(result){
			            				if (result != null){
								$( "#progressbar" ).progressbar( "value", result.Perc);
								if (result.Msg){
									dialog.setTitle(result.Msg);
								}
								if (result.Perc >= 100){
									dialog.setTitle('任务完成');
									dialog.getModalFooter().show();
									clearInterval(interval);
								}
			            				}else{
			            					dialog.setTitle('任务完成');
			            					dialog.getModalFooter().show();
			            					$( "#progressbar" ).progressbar( "value", 100);
			            					clearInterval(interval);	
			            				}
						}).error(function() { 
							clearInterval(interval);
							dialog.getModalFooter().show();
					    		BootFrame.alert("服务器发生错误，未获取到数据",null,"错误",true); 
						});
					},1000);
		            		}
			});
			dialog.realize();
			dialog.getModalHeader().css('padding','10px 10px 10px 15px');
			dialog.getModalFooter().css('padding','10px 15px');
			dialog.getModalFooter().hide();
			dialog.open();
		},//<! dialog >
		gritter: function (msg,title,time) {
			$.gritter.add({
				title: title?title:'<i class="ace-icon fa fa-bell bigger-120 blue"></i>&nbsp<font color="orange">提示</font>',
				text: msg,
				class_name: 'gritter-info gritter-radius',
				time:time?time:2000
			});
		},//<! gritter >
		uploader: function(){
			var p,u,l;
			var t = '';
			var t2 = '';
			var s = 10;
			var b = [];
			var hideclose = false;
			var w = 600;
			var h = 100;
			var dobj;
			return{
				url: function(url){
					l = url;
				},
				fileSize: function(size){
					s = size;
				},
				fileType: function(type){
					if (type=='image'){
						t = 'gif,jpg,jpeg,bmp,png';
						t2 = 'image/*';
					}
				},
				show: function(){
					if (!hideclose){
						var bt = {
							label: '取消',
							cssClass: 'btn-success btn-sm',
					                    	action: function(dialogItself){
				                        		dialogItself.close();
					                    	}
						};
						var sc = {
							label: '开始上传',
							cssClass: 'btn-primary btn-sm',
					                    	action: function(dialogItself){
					                    		var s = u.getFiles();
					                    		for (var i= 0 ; i < s.length ; i++){
					                    			if (s[i].getStatus() == 'inited' ){
					                    				u.upload();	
					                    				return;	
					                    			}
					                    		}
					                    		BootFrame.alert("请选择要上传文件的文件",null,null,true,null,200);
					                    	}
						};
						b.push(sc);
						b.push(bt);
					}
				  	dobj = BootstrapDialog.show({
					                title: '文件上传',
					                message: '<div id="uploader" class="uploader-dialog">'+
				                					 	'<div class="btns">'+
											        '<div id="picker" class="webuploader-container">选择文件</div>'+
										    	'</div>'+
										    	'<div id="thelist" class="uploader-list"></div>'+
										'</div>',
					                cssClass: 'dialog-'+w+' dialog-h-'+h,
					                closeByBackdrop: false,   //点击空白位置关闭窗口失效
					                draggable: true,
					                buttons: b,
					                onshown:function(){
					                	 u = WebUploader.create({
									        // swf文件路径
									        swf: '/static/webuploader/Uploader.swf',
									        // 文件接收服务端。
									        server: '/data/image/Upload',
									        fileVal :'fileupload' ,
									        // 选择文件的按钮。可选。
									        // 内部根据当前运行是创建，可能是input元素，也可能是flash.
									        // mutiple 是否允许多选文件上传
									        pick: {id:'#picker',multiple: false},
									        //可上传的文件数量
									        fileNumLimit: s,
									        // 不压缩image, 默认如果是jpeg，文件上传前会压缩一把再上传！
									        resize: false,
									        // 只允许选择的文件，可选。
									        accept: {
									            title: '文件',
									            extensions: t,
									            mimeTypes: t2
									        }
									    });
				                				// 当有文件被添加进队列的时候
									    u.on( 'fileQueued', function( file ) {
									          var icon = 'fa-picture-o file-image';
									          if ( file.type.indexOf('image') == -1 ) {
									                icon = 'fa-file';
									          }
									          $(".uploader-list").append('<div id="' + file.id + '" class="ace-file-input item">'+
									            '<span id="span' + file.id + '" class="ace-file-container selected" data-title="待上传">'+
									                '<div class="progress progress-striped active" style="height:28px;background: #ffffff;">'+
									                '<div class="progress-bar" role="progressbar" style="width: 0%">' +
									                '<span class="ace-file-name" style="display:inherit;position: absolute;" data-title="' + file.name + '">'+
									                    '<i class=" ace-icon fa ' + icon + '"></i>'+
									                '</span>'+
									                '</div>'+
									                '</div>'+
									            '</span>'+
									            '<a id="a' + file.id + '" class="remove" href="#;">'+
									                '<i class=" ace-icon fa fa-times"></i>'+
									            '</a>'+
									        '</div>');
									          $('#a'+file.id).click(function(){
									       		 u.removeFile(file.id,true);
									   	 });
									    });

									    // 文件被移除
									    u.on( 'fileDequeued',function( file ){
									        $('div#'+file.id).remove();
									    });

									    // 文件上传过程中创建进度条实时显示。
									    u.on( 'uploadProgress', function( file, percentage ) {
									        var $li = $( '#span'+file.id ),
									            $percent = $li.find('.progress .progress-bar');
									        // 避免重复创建
									        if ( !$percent.length ) {
									            $percent = $('<div class="progress progress-striped active">' +
									              '<div class="progress-bar" role="progressbar" style="width: 0%">' +
									              '</div>' +
									            '</div>').appendTo( $li ).find('.progress-bar');
									        }
									        $( '#span'+file.id ).attr('data-title','上传中');
									        $percent.css( 'width', percentage * 100 + '%' );
									    });

									    // 文件上传成功
									    u.on( 'uploadSuccess', function( file ) {
									        $( '#span'+file.id ).attr('data-title','已上传');
									        $( '#a'+file.id ).addClass( "success" );
									        $( '#a'+file.id ).unbind( "click" );
									        $( '#a'+file.id ).children("i").removeClass('fa-times').addClass('fa-check');
									    });

									    // 文件上传失败
									    u.on( 'uploadError', function( file ) {
									        $( '#span'+file.id ).attr('data-title','有错误');
									    });

									    // 文件上传完成（不管成功失败）
									    u.on( 'uploadComplete', function( file ) {
									        //$( '#'+file.id ).find('.progress').fadeOut();
									    });									    
					                }
				           });
				  	dobj.getModalHeader().css('padding','10px 10px 10px 15px');
				  	dobj.getModalFooter().css('padding','10px 15px');
				},
				del: function(id){
					u.removeFile(id,true);
				}
			}
		}
	}//<! return >
   }();

})(jQuery);
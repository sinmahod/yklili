/*!
 * Extend BootstrapValidator (http://yklili.com)
 *
 * @BootstrapValidator          v0.5.2
 * @jQuery                           v.2.1.4
 * @author     		  gl
 * @time			  2016/12/02
 * @email			  sinmahod@qq.com
 * @example                         verify="notEmpty;phone"  //非空+手机号码
 */
 (function ($) {
 
/*!
 * 引用博客园loogn的替换{0},{1}...{n}代码
 * @example  "aaa,{0},{1}".format("bbb","ccc")  输出 aaa,bbb,ccc
 */
 String.prototype.format = function(args) {
    var result = this;
    if (arguments.length > 0) {    
        if (arguments.length == 1 && typeof (args) == "object") {
            for (var key in args) {
                if(args[key]!=undefined){
                    var reg = new RegExp("({" + key + "})", "g");
                    result = result.replace(reg, args[key]);
                }
            }
        }
        else {
            for (var i = 0; i < arguments.length; i++) {
                if (arguments[i] != undefined) {
      　　　　var reg= new RegExp("({)" + i + "(})", "g");
                    result = result.replace(reg, arguments[i]);
                }
            }
        }
    }
    return result;
}

 window.BValidate = function(options){
        var form ,f;
        var fields = {};
        var varMsg = {
                notEmpty : "此项不可为空",
                identical : "两次输入的密码不一致",   
                lessThan :  "输入的数值必须小于{0}",
                lessThans : "输入的数值必须小于等于{0}",
                greaterThan : "输入的数值必须大于{0}",
                greaterThans : "输入的数值必须大于等于{0}",
                zipCode : "请输入正确的邮编",
                uri : "请输入正确的网址",
                emailAddress : "请输入正确的邮箱地址",
                stringLength : "输入的字符长度至少{0}位，最多{1}位",
                phone : "请输入正确的手机号码",
                regexp : "输入不合法",
                remote : "远程校验不合法"
          }
          
          if(options == undefined || options["dialogId"] == undefined){
                form = "form";
          }else{
                form = "#" + options.dialogId + " form";
          }

            $(form).each(function(){
		f = $(this).attr("id");  					                    //form id
		if (f != undefined){
			$(this).find("[verify]").each(function(){
				var v = $(this).attr("verify"); 	                               //校验值  verify="notEmpty"  返回 notEmpty
				var name = $(this).attr("name");
				var id = $(this).attr("id");
				if (id == undefined && name == undefined){
					return true;					//相当于for循环的continue
				}
				if (name == undefined && id != undefined){	//填充id，保证id与name两个属性都有
					$(this).attr("name",id);
					name = id;
				}
				if (name != undefined && id == undefined){	//填充name，保证id与name两个属性都有
					$(this).attr("id",name);
					id = name;
				}
				//开始添加校验参数
                                         fields[name] = verifyToJson(v);
			});
		}	
	});

	options = $.extend({
            message: '表单验证',
            feedbackIcons: {
                valid: 'glyphicon glyphicon-ok',
                invalid: 'glyphicon glyphicon-remove',
                validating: 'glyphicon glyphicon-refresh'
            },
            fields: fields
            }, options);

          //整体规则参数拆分
          function verifyToJson(verifystr){
                    var verifyjson = {};
                    var validators = {};
                    var vs = verifystr.split(";");
                    for ( var i in vs){
                        var t,p,m;
                        var vm = vs[i].split("|");
                        if (vm.length == 2){
                            m = vm[1];
                        }
                        var v = /(.*)\((.*)\)/.exec(vm[0]);     //stringLength(6-30) = ["stringLength(6-30)","stringLength","6-30"]
                        if (v != null && v.length == 3){
                            t = v[1];
                            p = v[2];
                        }else{
                            t = vm[0];
                        }

                         switch (t) {
                             case 'remote':
                                var urth = p.split(",");
                                validators["threshold"] = urth[1];
                                validators[t] = verifyRule(t,urth[0],m);
                                break;
                             case 'phone':
                                validators["regexp"] = verifyRule(t,p,m);
                                break;
                             default:
                                validators[t] = verifyRule(t,p,m);    
                                break;
                         }
                    }
                    verifyjson["validators"] = validators;
                    return verifyjson;
          }

         

          //单个规则构造
          function verifyRule(type,param,msg){
                var verifyType = {};
                
                if (varMsg[type] == undefined ){
                    return verifyType;
                }

                if (msg){
                        verifyType["message"] = msg;
                } 
                if (param){        //包含()的
                        switch (type) {
                             case 'remote':
                                verifyType["url"] = param;
                                verifyType["delay"] = 2000;  //每输入一个字符，就发ajax请求，服务器压力还是太大，设置2秒发送一次ajax（默认输入一个字符，提交一次，服务器压力太大）
                                verifyType["type"] = "POST";
                                if (verifyType["message"]  == undefined){
                                    verifyType["message"] = varMsg[type];
                                }
                                break;
                             case 'regexp':
                                verifyType["regexp"] = eval(param) ;
                                if (verifyType["message"]  == undefined){
                                    verifyType["message"] = varMsg[type];
                                }
                                break;
                            case 'identical':
                                verifyType["field"] = param ;
                                 if (verifyType["message"]  == undefined){
                                    verifyType["message"] = varMsg[type];
                                }
                                break;
                            case 'stringLength':
                                var n = param.split(",");
                                if (n != 2){
                                    return false;
                                }
                                verifyType["min"] = n[0] ;
                                verifyType["max"] = n[1] ;
                                if (verifyType["message"]  == undefined){
                                    verifyType["message"] = varMsg[type].format(n[0],n[1]);
                                }
                                break;
                            default:
                                verifyType["value"] = param;
                                verifyType["inclusive"] = false;
                                if (type == 'lessThans' || type == 'greaterThans'){
                                    verifyType["inclusive"] = true;
                                }
                                if (verifyType["message"]  == undefined){
                                    verifyType["message"] = varMsg[type].format(param);
                                }
                                break;
                        }
                }else{
                    if (type == 'phone'){
                        verifyType["regexp"] = /^1[3|5|8]{1}[0-9]{9}$/ ;
                    }
                    if (verifyType["message"]  == undefined){
                        verifyType["message"] = varMsg[type];
                    }
                }
                return verifyType;
          }
          if (f){
              $('#'+f).bootstrapValidator(options);  
          }
 };
 })(jQuery);
$(function () {

window.BValidate = function(options){

            $("ul[list]").each(function(){
                var v = $(this).attr("verify");                                    //校验值  verify="notEmpty"  返回 notEmpty
                var name = $(this).attr("name");
                var id = $(this).attr("id");
            });

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
          };

            options = $.extend({
            message: '表单验证',
            feedbackIcons: {
                valid: 'ace-icon glyphicon glyphicon-ok',
                invalid: 'ace-icon glyphicon glyphicon-remove',
                validating: 'ace-icon glyphicon glyphicon-refresh'
            },
            fields: fields
            }, options);

 };

 })(jQuery);
                    var options = {
                        bootstrapMajorVersion: 2, //版本
                        currentPage: 1, //当前页数
                        totalPages: 2, //总页数
                        itemTexts: function (type, page, current) {
                            switch (type) {
                                case "first":
                                    return "首页";
                                case "prev":
                                    return "上一页";
                                case "next":
                                    return "下一页";
                                case "last":
                                    return "末页";
                                case "page":
                                    return page;
                            }
                        },//点击事件，用于通过Ajax来刷新整个list列表
                        onPageClicked: function (event, originalEvent, type, page) {
                            $.ajax({
                                url: "/OA/Setting/GetDate?id=" + page,
                                type: "Post",
                                data: "page=" + page,
                                success: function (data1) {
                                    if (data1 != null) {
                                        $.each(eval("(" + data + ")").list, function (index, item) { //遍历返回的json
                                            $("#list").append('<table id="data_table" class="table table-striped">');
                                        });
                                    }
                                }
                            });
                        }
                    };
                    $('#example').bootstrapPaginator(options);
                }
            }
        });
    })
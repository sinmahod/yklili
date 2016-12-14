/*!
 * Extend Dropzone (http://yklili.com)
 *
 * @BootstrapValidator          v0.5.2
 * @jQuery                           v.2.1.4
 * @author            gl
 * @time              2016/12/14
 * @email             sinmahod@qq.com
 * @example                         <upload>
 */
 (function ($) {

 window.UploadEx = function(divid){
    $('body').find('upload').each(function(){
        alert($(this).attr("id"))
    });  

    var options = {
            url: "/data/image/Upload",
            previewTemplate: $('#preview-template').html(),
            paramName: "filetest", //上传域name
            thumbnailHeight: 120,
            thumbnailWidth: 120,
            maxFilesize: 0.5,   //上传文件不能超过N(m)
            addRemoveLinks : true,  //添加移除文件
            autoProcessQueue: false,    //取消自动上传
            acceptedFiles: ".jpg,.png,.jpeg,.gif,.bmp" ,  //允许上传的类型
            parallelUploads:5,  //每次最多上传数量
            maxFiles: 5,    //最多可以添加数量
            uploadMultiple: true,  //一次提交多个文件
            init: function() {
                var submitButton = document.querySelector("#submit-all")
                myDropzone = this; // closure

                submitButton.addEventListener("click", function() {
                    myDropzone.processQueue(); // 手动上传
                });

                // 文件上传之前执行，添加其他表单参数
                this.on("sendingmultiple", function(files, xhr, formData) {
                    formData.append("id", "test");
                });
                this.on("successmultiple", function(files, response) {
                    alert("1"+response)
                });
                this.on("errormultiple", function(files, response) {
                    alert("2"+response)
                }); 
            },
            //addRemoveLinks : true,
            //dictRemoveFile: 'Remove',
            dictDefaultMessage: '<span class="bigger-150 bolder"><i class="ace-icon fa fa-caret-right red"></i> 点击此处或将要上传的图片拖放到此处<br /> \
            <i class="upload-icon ace-icon fa fa-cloud-upload blue fa-3x"></i>',
            dictFallbackMessage: 'Fallback 情况下的提示文本。', 
            dictInvalidInputType: '文件类型被拒绝 ',
            dictFileTooBig: '文件大小过大 ',
            dictCancelUpload: '取消上传 ',
            dictCancelUploadConfirmation: '取消上传确认 ',
            dictRemoveFile: '移除文件 ',
            dictMaxFilesExceeded: '超过最大文件数 ',

            thumbnail: function(file, dataUrl) {
                if (file.previewElement) {
                    $(file.previewElement).removeClass("dz-file-preview");
                    var images = $(file.previewElement).find("[data-dz-thumbnail]").each(function() {
                        var thumbnailElement = this;
                        thumbnailElement.alt = file.name;
                        thumbnailElement.src = dataUrl;
                    });
                    setTimeout(function() {
                        $(file.previewElement).addClass("dz-image-preview");
                    }, 1);
                }
            }
        }
    //new Dropzone('#dropzone', options);
 };
 })(jQuery);
(function ($) {
 
  window.Ewin = function () {
    var html = '<div id="[Id]" class="modal fade" role="dialog" aria-labelledby="modalLabel">' +
               '<div class="modal-dialog modal-sm">' +
                 '<div class="modal-content">' +
                   '<div class="modal-header">' +
                     '<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">×</span><span class="sr-only">Close</span></button>' +
                     '<h4 class="modal-title" id="modalLabel">[Title]</h4>' +
                   '</div>' +
                   '<div class="modal-body">' +
                   '<p>[Message]</p>' +
                   '</div>' +
                    '<div class="modal-footer">' +
    '<button type="button" class="btn btn-default cancel" data-dismiss="modal">[BtnCancel]</button>' +
    '<button type="button" class="btn btn-primary ok" data-dismiss="modal">[BtnOk]</button>' +
  '</div>' +
                 '</div>' +
               '</div>' +
             '</div>';
 
 
    var dialogdHtml = '<div id="[Id]" class="modal fade" role="dialog" aria-labelledby="modalLabel">' +
               '<div class="modal-dialog">' +
                 '<div class="modal-content">' +
                   '<div class="modal-header">' +
                     '<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">×</span><span class="sr-only">Close</span></button>' +
                     '<h4 class="modal-title" id="modalLabel">[Title]</h4>' +
                   '</div>' +
                   '<div class="modal-body">' +
                   '</div>' +
                   '<div class="modal-footer">' +
                     '<button type="button" class="btn btn-default cancel" data-dismiss="modal">[BtnCancel]</button>' +
                     '<button type="button" class="btn btn-primary no" data-dismiss="modal" style="display:none">[BtnNo]</button>' +
                     '<button type="button" class="btn btn-primary en" style="display:none" data-dismiss="modal">[BtnEn]</button>' +
                     '<button type="button" class="btn btn-primary ok" data-dismiss="modal">[BtnOk]</button>' +
                   '</div>' +
                 '</div>' +
               '</div>' +
             '</div>';
    var reg = new RegExp("\\[([^\\[\\]]*?)\\]", 'igm');
    var generateId = function () {
      var date = new Date();
      return 'mdl' + date.valueOf();
    }
    var init = function (options) {
      options = $.extend({}, {
        title: "&nbsp;",
        message: "提示内容",
        btnok: "确定",
        btncl: "取消",
        width: 200,
        auto: false
      }, options || {});
      var modalId = generateId();
      var content = html.replace(reg, function (node, key) {
        return {
          Id: modalId,
          Title: options.title,
          Message: '<center style="font-size: 18px;">'+options.message+'</center>',
          BtnOk: options.btnok,
          BtnCancel: options.btncl,
          BtnNo: options.btnno
        }[key];
      });
      $('body').append(content);
      $('#' + modalId).modal({
        width: options.width,
        backdrop: 'static'
      });
      $('#' + modalId).on('hide.bs.modal', function (e) {
        $('body').find('#' + modalId).remove();
      });
      return modalId;
    }
 
    return {
      alert: function (options) {
        if (typeof options == 'string') {
          options = {
            message: options
          };
        }
        var id = init(options);
        var modal = $('#' + id);
        modal.find('.ok').removeClass('btn-success').addClass('btn-primary');
        modal.find('.cancel').hide();
        return {
          id: id,
          on: function (callback) {
            if (callback && callback instanceof Function) {
              modal.find('.ok').click(function () { callback(true); });
            }
          },
          hide: function (callback) {
            if (callback && callback instanceof Function) {
              modal.on('hide.bs.modal', function (e) {
                callback(e);
              });
            }
          }
        };
      },
      confirm: function (options) {
        var id = init(options);
        var modal = $('#' + id);
        modal.find('.ok').removeClass('btn-primary').addClass('btn-success');
        modal.find('.cancel').show();
        return {
          id: id,
          on: function (callback) {
            if (callback && callback instanceof Function) {
              modal.find('.ok').click(function () { callback(true); });
            }
          },
          off: function (callback) {
              if (callback && callback instanceof Function) {
                modal.find('.cancel').click(function () { callback(true); });
              }
          },
          hide: function (callback) {
            if (callback && callback instanceof Function) {
              modal.on('hide.bs.modal', function (e) {
                callback(e);
              });
            }
          }
        };
      },
      dialog: function (options) {
        options = $.extend({}, {
          title: 'title',
          url: '',
          width: 800,
          height: 200,
          btnok: "保存",
          btncl: "取消",
          btnno: "不选专家",
          btnen: "确定",
          onReady: function () { },
          onShown: function (e) { }
        }, options || {});
        var modalId = generateId();
        var content = dialogdHtml.replace(reg, function (node, key) {
          return {
            Id: modalId,
            Title: options.title,
            BtnOk: options.btnok,
            BtnCancel: options.btncl,
            BtnNo: options.btnno,
            BtnEn: options.btnen
          }[key];
        });
        $(".modal.fade").remove();
        $('body').append(content);
        var target = $('#' + modalId);
        target.find('.modal-body').html('<iframe id="myframes" name="myframes" src="'+options.url+'" width="100%" height="100%" frameborder="no" border="0" marginwidth="0" marginheight="0" scrolling="yes" ></iframe>')
        if(options.height){
          target.find('.modal-body').css("height",options.height);
        }
        if(options.width&&options.width!='800'){
          target.find('.modal-content').css("width",options.width);
        }
      /*  debugger;
        if(options.width&&options.width=='510'){
          target.find('.no').show();
        }*/
        if(options.title&&options.title=='选择专家'){
          target.find('.no').show();
        }
        if(options.width&&options.title=='同约人员'){
          target.find('.cancel').hide();
          target.find('.ok').hide();
          target.find('.en').show();
        }
        if(options.title&&options.title=='查看回复'){
          target.find('.cancel').hide();
          target.find('.ok').hide();
        }
        if(options.hidden){
          if(options.hidden=='footer'){
            target.find('.modal-footer').hide();
          }
        }
        if (options.onReady())
          options.onReady.call(target);
          target.modal();
          target.on('shown.bs.modal', function (e) {
            if (options.onReady(e))
              options.onReady.call(target, e);
          });
         /* target.on('hide.bs.modal', function (e) {
            $('body').find(target).remove();
          });*/
        return {
            id: modalId,
            on: function (callback) {
              if (callback && callback instanceof Function) {
                target.find('.ok').click(function () { 
                  if(callback(true)){
                    $('body').find(target).remove();
                  }
              });
              }
            },
            no:function(callback){
               if (callback && callback instanceof Function) {
                  target.find('.no').click(function () { callback(true); });
                 }
            },
            en:function(callback){
             if (callback && callback instanceof Function) {
                  target.find('.en').click(function () { callback(true); });
                }
           },
            off: function (callback) {
                if (callback && callback instanceof Function) {
                  target.find('.cancel').click(function () { callback(true); });
                }
            },
            hide: function (callback) {
              if (callback && callback instanceof Function) {
                target.on('hide.bs.modal', function (e) {
                  callback(e);
                });
              }
            }
          };
      }
    }
  }();
})(jQuery);
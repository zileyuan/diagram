webix.i18n.locales["zh-CN"].calendar.clear = "清除";
webix.i18n.setLocale("zh-CN");


function ajax(req_method, req_url, req_data, done_func, fail_func, always_func) {
    var setting = {
        url: req_url,
        type: req_method,
        dataType: 'json',
    };
    if (req_data) {
        setting.data = req_data;
    }
    return $.ajax(setting)
        .done(function (data) {
            if (done_func) {
                done_func(data);
            }
        })
        .fail(function () {
            if (fail_func) {
                fail_func();
            }
        })
        .always(function () {
            if (always_func) {
                always_func();
            }
        });
}

//获取url参数
function getUrlParam(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
    var r = window.location.search.substr(1).match(reg);  //匹配目标参数
    console.log(r);
    if (r != null) return unescape(r[2]);
    return null; //返回参数值
}
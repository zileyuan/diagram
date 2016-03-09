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
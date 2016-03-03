global_tag_zh = "zh";
global_tag_namespace_storage = "audit-embed";
global_tag_space = "space";
global_tag_topest = "topest";
global_tag_inner = "inner";
global_tag_uid = "uid";
global_tag_inserted = "add";
global_tag_deleted = "delete";
global_tag_update = "update";
global_tag_nav_label = "NL-";
global_tag_nav_item = "N-";
global_tag_title = "T-";
global_tag_prefix = "I-";
global_tag_bind = "bind";
global_tag_folder = "folder";
global_tag_amount = "amount";
global_tag_reveal = "reveal";
global_tag_data_reveal = "data-reveal";
global_tag_data_open = "data-open";
global_tag_data_sticky_container = "data-sticky-container";
global_tag_accordion = "accordion";
global_tag_check = "check-";
global_tag_radio = "radio-";
global_tag_tmpId = "tmp_685ab04c-7c12-49bd-9a35-2bb37e490671";


function appInit() {
    //window.Foundation.Magellan.defaults.barOffset = 20;
    $(document).foundation();
    Foundation.IHearYou();
    webix.i18n.setLocale(global_tag_zh);
    do_resize();
}

function show_tip(message) {
    webix.message(message);
}

function lower_trim(text) {
    return text.toLowerCase().replace(/(^\s*)|(\s*$)/g, '');
}

function include(sample, filter) {
    if (sample && filter) {
        sample = sample.toString().toLowerCase();
        filter = lower_trim(filter.toString());
        return sample.indexOf(filter) !== -1;
    } else {
        return false;
    }
}

function do_resize() {
    var last_li_height = $("#Content").children("ul").last().children("li").last().height();
    var window_height = $(window).height();
    var space_height = window_height - last_li_height;
    if (space_height > 0) {
        $("#" + global_tag_space).css({
                height: space_height,
            }
        );
    }

    $("#nav").height(parseInt($(window).height()) - 58);
}

function get_local_storage() {
    var ns = $.initNamespaceStorage(global_tag_namespace_storage);
    var storage = ns.localStorage;
    return storage;
}

function store_position() {
    var topest = $("." + global_tag_topest);
    if (topest) {
        var uid = topest.attr(global_tag_uid);
        var storage = get_local_storage();
        var top = $(document).scrollTop();
        if (uid && top) {
            storage.set(uid, top);
        }
    }
}

function lazy_load() {
    var top = $(document).scrollTop();
    var height = $(document).height() - $(window).height();
    ;
    if (top > parseInt(height * 0.9)) {
        var topest = $("." + global_tag_topest);
        if (topest) {
            var uid = topest.attr(global_tag_uid);
            uid = global_folder_list.next_folder(uid);
            if (uid) {
                var radio = $("#" + global_tag_radio + uid);
                if (radio) {
                    radio.click();
                }
            }
        }
        //do it;
    }
}

function do_bottom() {
    var top = $(document).scrollTop();
    var height = $(document).height() - $(window).height();
    if (top > parseInt(height * 0.9)) {
        setTimeout(lazy_load, 2000);
    }
}

function sync_position() {
    var topest = $("." + global_tag_topest);
    if (topest) {
        var uid = topest.attr(global_tag_uid);
        var storage = get_local_storage();
        var top = storage.get(uid);
        if (top) {
            $(document).scrollTop(top);
        }
    }

    $(window).scroll(function () {
        store_position();
        //do_bottom();
    });
}

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

function clearJson(json) {
    for (var i in json) {
        if (i.substr(0, 1) == '$' || i == 'id') {
            delete json[i];
        }
    }
    return json;
}

$(document).ready(function () {
    do_resize();
    $(window).resize(function () {
        do_resize();
    });
});

function getMajorIcon(name) {
    if (include(name, "电气")) {
        return "icon-dianqi";
    } else if (include(name, "建筑")) {
        return "icon-jianzhu";
    } else if (include(name, "结构")) {
        return "icon-jiegou";
    } else if (include(name, "暖通")) {
        return "icon-nuantong";
    } else if (include(name, "排水")) {
        return "icon-paishui";
    } else if (include(name, "岩土")) {
        return "icon-yantu";
    }  else if (include(name, "岩土")) {
        return "icon-yantu";
    }  else if (include(name, "新增")) {
        return "fa fa-thumb-tack";
    } else if (include(name, "收藏")) {
        return "fa fa-star-o";
    } else if (include(name, "删除")) {
        return "fa fa-trash-o";
    } else if (include(name, "取消")) {
        return "fa fa-trash-o";
    } else if (include(name, "修改")) {
        return "fa fa-pencil-square-o";
    } else if (include(name, "复制")) {
        return "fa fa-clone";
    } else if (include(name, "选用")) {
        return "fa fa-share-square-o";
    } else if (include(name, "引用")) {
        return "fa fa-share-alt";
    } else if (include(name, "定位")) {
        return "fa fa-map-marker";
    } else if (include(name, "查看")) {
        return "fa fa-commenting-o";
    } else if (include(name, "回复")) {
        return "fa fa-reply";
    } else {
        return "fa fa-building";
    }
}

function update_status(status, event) {
    for (var i in global_grid_list.grids) {
        var grid = global_grid_list.grids[i];
        var item = grid.grid.find(function (obj) {
            return (obj.Id == event.fileid);
        }, true);
        if (item && !Array.isArray(item)) {
            item[status] = event.filestate;
            grid.grid.update_item(item.id, item);
            return;
        }
    }
}


function do_ajax() {
    var data = {};
    ajax("POST", "/overview", data, function (result) {
        webix.message(JSON.stringify(result));
    })
}
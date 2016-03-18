$(document).ready(function() {
	var data = {}
	ajax("POST", "/custtype", data, function (result) {
		var list = $$('custtype').getPopup().getList();
		list.clearAll();
		list.parse(result.data);
	})
	ajax("POST", "/storecode", data, function (result) {
		var list = $$('storecode').getPopup().getList();
		list.clearAll();
		list.parse(result.data);
	})
})

function do_cust() {
    //var data = {
    //    CustName: $$('custname').getValue(),
		//Mobile: $$('mobile').getValue(),
    //    CustType: $$('custtype').getValue(),
    //    StoreCode: $$('storecode').getValue(),
    //};

    $$('win2').show()
}

webix.ui({
    view: "toolbar",
    container: "toolbar",
    elements: [
       {
            view: "combo",
            id:"storecode",
            width:210,
            label: '所属门店',
            labelWidth:80,
			options: []
       },
	   {
            view: "combo",
            id:"custtype",
            width:205,
            label: '客户类型',
            labelWidth:80,
			options: []
       },
       {
            view:"text",
            id: "custname",
			width:200,
            label:"客户名称",
			labelWidth: 80
        },
		{
            view:"text",
            id: "mobile",
			width:185,
            label:"手机号",
			labelWidth: 60
        },
		{
            view: "button",
            id:"btn",
			width:80,
			label: '查询',
			click: do_cust
        },
		{
            view: "button",
            id:"btn",
			width:80,
			label: '新增客户',
			click: do_cust
        },
    ]
});

webix.ui.datafilter.pageAvgColumn = webix.extend({
    refresh:function(master, node, value){
        var result = 0;
        var page = master.getPage();
        var pager = master.getPager();
        var cc = pager.data.size;
        var iid = master.getIdByIndex(page * cc);
        master.mapCells(iid, value.columnId, cc, 1, function(value){
            value = value*1;
            if (!isNaN(value))
                result+=value;
            return value;
        });

        node.firstChild.innerHTML = Math.round(result / cc);
    }
}, webix.ui.datafilter.summColumn);

webix.ui.datafilter.pageSummColumn = webix.extend({
    refresh:function(master, node, value){
        var result = 0;
        var page = master.getPage();
        var pager = master.getPager();
        var cc = pager.data.size;
        var iid = master.getIdByIndex(page * cc);
        master.mapCells(iid, value.columnId, cc, 1, function(value){
            value = value*1;
            if (!isNaN(value))
                result+=value;
            return value;
        });

        node.firstChild.innerHTML = Math.round(result);
    }
}, webix.ui.datafilter.summColumn);

webix.ui.datafilter.avgColumn = webix.extend({
    refresh:function(master, node, value){
        var result = 0;
        master.mapCells(null, value.columnId, null, 1, function(value){
            value = value*1;
            if (!isNaN(value))
                result+=value;
            return value;
        });

        node.firstChild.innerHTML = Math.round(result / master.count());
    }
}, webix.ui.datafilter.summColumn);

grida = webix.ui({
    container:"testA",
    view:"datatable",
    id: "dtable",
    columns:[
        { id:"crname",	header:"客户名称",width:190},
		{ id:"crtitle",	header:"称呼",width:80},
        { id:"store",header:"所属门店",width:147},
		{ id:"kehlx",header:"客户类型",width:120},
		{ id:"crsex",header:"性别",width:55},
        { id:"mobile", header:"手机号",width:120},
		{ id:"crbirthday", header:"生日",width:105},
		{ id:"",template:"<input class='detail' type='button' value='详情'><input class='other' type='button' value='积分兑换'>",
            css:"padding_less",width:100},
    ],
    select:"cell",
    autowidth:true,
    //footer:true,
    pager:{
        template:"{common.first()} {common.prev()} {common.pages()} {common.next()} {common.last()}",
        container:"paging_here",
        size:100,
        group:5
    },
});

grida.on_click.detail=function(e, id, trg){
    var o = this.getItem(id.row)

    $$('win2').show()

    return false;
};

var form = {
    view:"form",
    borderless:true,
    elements: [
        {
            rows:[
                {
                    cols:[
                        { view:"combo", label:'Store', id:"Store" },
                        { view:"text", label:'Crname', id:"Crname" },
                    ]}
            ]
        },

        { view:"text", label:'Crtitle', id:"Crtitle" },
        { view:"combo", label:'Kehlx', id:"Kehlx" },
        { view:"text", label:'Crsex', name:"Crsex" },
        { view:"text", label:'Mobile', name:"Mobile" },
        { view:"datepicker", label:'Crbirthday', name:"Crbirthday" },
        { view:"button", value: "Submit", click:function(){


            var data = {}

            ajax("POST", "/customer", data, function (result) {
                if (result.success) {
                    var data = result.data;
                    $$('dtable').clearAll();
                    $$('dtable').parse(data);
                } else {
                    webix.message(result.data);
                }
            })
        }},
        { view:"button", value: "Cancel", click:function(){
                this.getTopParentView().hide(); //hide window
        }}
    ],
    elementsConfig:{
        labelPosition:"top",
    }
};

webix.ui({
    view:"window",
    id:"win2",
    width:600,
    position:"center",
    modal:true,
    head:"User's data",
    body:webix.copy(form)
});
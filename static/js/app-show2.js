var kehlx_data,store_data;

$(document).ready(function() {
	var data = {}
	ajax("POST", "/custtype", data, function (result) {
		var list = $$('custtype').getPopup().getList();
		list.clearAll();
		list.parse(result.data);
		
		list = $$('Kehlx').getPopup().getList();
		list.clearAll();
		list.parse(result.data);
		
		kehlx_data = result.data
	})
	ajax("POST", "/storecode", data, function (result) {
		var list = $$('storecode').getPopup().getList();
		list.clearAll();
		list.parse(result.data);
		
		store_data = result.data
		
//		list = $$('Store').getPopup().getList();
//		list.clearAll();
//		list.parse(result.data);
	})
})

function do_cust() {
    var data = {
        CustName: $$('custname').getValue(),
		Mobile: $$('mobile').getValue(),
        CustType: $$('custtype').getValue(),
        StoreCode: $$('storecode').getValue(),
    };
    ajax("POST", "/customer", data, function (result) {
		if (result.success) {
        var data = result.data;
        $$('dtable').clearAll();
        $$('dtable').parse(data);
		} else {
			webix.message(result.data);
		}
    })
}

function do_addcust() {
    //var data = {
    //    CustName: $$('custname').getValue(),
		//Mobile: $$('mobile').getValue(),
    //    CustType: $$('custtype').getValue(),
    //    StoreCode: $$('storecode').getValue(),
    //};
    $$('uuid').setValue("")
	$$('Uid').setValue("")
	$$('Crname').setValue("")
	$$('Crqcode').setValue("")
	//$$('Kehlx').setValue("")
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
            id:"btnfind",
			width:80,
			label: '查询',
			click: do_cust
        },
		{
            view: "button",
            id:"btnadd",
			width:80,
			label: '新增客户',
			click: do_addcust
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
//		{ id:"",template:"<input class='detail' type='button' value='详情'><input class='other' type='button' value='积分兑换'>",
//            css:"padding_less",width:100},
        { id:"",template:"<input class='detail' type='button' value='详情'>"},
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
	$$('uuid').setValue(o.id)
	$$('Uid').setValue(o.uid)
	$$('Crname').setValue(o.crname)
	$$('Crqcode').setValue(o.crqcode)
	$$('Kehlx').setValue(o.kehlxid)
	//$$('Kehlx').disable();
    $$('win2').show()

    return false;
};

var form = {
    view:"form",
    borderless:true,
    elements: [
		{ view:"text", label:'uuid', id:"uuid", hidden:true },
        { view:"text", label:'Uid', id:"Uid", hidden:true },
        {
            rows:[
                {
                    cols:[
                        { view:"combo", label:'Kehlx', id:"Kehlx" ,options: []},
                        { view:"text", label:'Crname', id:"Crname" },
                    ]}
            ]
        },
        { view:"text", label:'Crqcode', id:"Crqcode" },
        { view:"text", label:'Crtitle', id:"Crtitle" },
        { view:"text", label:'Crsex', name:"Crsex" },
        { view:"text", label:'Mobile', name:"Mobile" },
        { view:"datepicker", label:'Crbirthday', name:"Crbirthday" },
        { view:"button", value: "Submit", click:function(){

			var b = this;
			
            var data = {
				uid: $$('Uid').getValue(),
				crname: $$('Crname').getValue(),
				kehlxid: $$('Kehlx').getValue(),
				crqcode: $$('Crqcode').getValue(),
				storeid: getUrlParam("store")
			}
			for (v in store_data) {
				if (store_data[v].id == data.storeid) {
					data.store = store_data[v].value;
					break;
				}
			}
			for (v in kehlx_data) {
				if (kehlx_data[v].id == data.kehlxid) {
					data.kehlx = kehlx_data[v].value;
					break;
				}
			}

            ajax("POST", "/updcust", data, function (result) {
                if (result.success) {
					//var o = grida.getItem($('uuid').getValue())
					if ($$('uuid').getValue()=="") {
						data.uid = result.data;
					  	grida.add(data, 0)	
					} else {
					  grida.updateItem($$('uuid').getValue(), data);	
					}
					b.getTopParentView().hide();
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
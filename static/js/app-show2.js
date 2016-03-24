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
	$$('Crtitle').setValue("")
	$$('Crsex').setValue("")
	$$('Mobile').setValue("")
	$$('Crbirthday').setValue("")
	$$('Crmarriage').setValue("")
	$$('Crzip').setValue("")
	$$('Cridentity').setValue("")
	$$('Crhobby').setValue("")
	$$('Crmemo').setValue("")
	//$$('Kehlx').setValue("")
	$$('btnsave').enable()
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
        { id:"crname",	header:"客户名称",width:190, template: function(obj) {
            return "<a target='_blank' href = '/cardinfo?custid=" + obj.uid + "'>" + obj.crname + "</a>";}},
		{ id:"crqcode", header:"速查码",width:90},
		{ id:"crtitle",	header:"称呼",width:80},
        { id:"store",header:"所属门店",width:147},
		{ id:"kehlx",header:"客户类型",width:120},
		{ id:"crsex",header:"性别",width:55},
        { id:"mobile", header:"手机号",width:115},
		{ id:"crbirthday", header:"生日",width:100},
		//{ id:"",template:"<input class='detail' type='button' value='详情'>",  
          //  css:"padding_less",width:55},
		{ id:"",width:70,template:'<div class="webix_view webix_control webix_el_button" style="margin-left: 0px; width: 55px;"><div class="webix_el_box" style="padding:0px"><button type="button" class="webixtype_base detail">详情</button></div></div>'},
    ],
	//select:"row",
    select:"cell",
    //autowidth:true,
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
	$$('Crtitle').setValue(o.crtitle)
	$$('Kehlx').setValue(o.kehlxid)
	if (o.storeid != getUrlParam("store")) {
	  $$('btnsave').disable()	
	} else {
	  $$('btnsave').enable()	
	}
	$$('Crsex').setValue(o.crsex)
	$$('Mobile').setValue(o.mobile)
	$$('Crbirthday').setValue(o.crbirthday)
	$$('Crmarriage').setValue(o.crmarriage)
	$$('Crzip').setValue(o.crzip)
	$$('Cridentity').setValue(o.cridentity)
	$$('Crhobby').setValue(o.crhobby)
	$$('Crmemo').setValue(o.crmemo)
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
                        { view:"combo", label:'客户类型', id:"Kehlx" ,options: []},
                        { view:"text", label:'姓名', id:"Crname" },
                    ]}
            ]
        },
        {
            rows:[
                {
                    cols:[
                        { view:"text", label:'速查码', id:"Crqcode" },
                        { view:"text", label:'称呼', id:"Crtitle" },
                    ]}
            ]
        },
		{
            rows:[
                {
                    cols:[
                        { view:"combo", label:'性别', id:"Crsex",options:[
                          { id:"男", value:"男"},
                          { id:"女", value:"女"},
						  { id:"不详", value:"不详"}
                          ]},
                        { view:"datepicker", label:'生日', id:"Crbirthday", stringResult: true },
                    ]}
            ]
        },
		{
            rows:[
                {
                    cols:[
                        { view:"text", label:'手机', id:"Mobile" },
                        { view:"text", label:'联系地址', id:"Crzip" },
                    ]}
            ]
        },
		{
            rows:[
                {
                    cols:[
                        { view:"combo", label:'婚否', id:"Crmarriage",options:[
                          { id:"未婚", value:"未婚"},
                          { id:"已婚", value:"已婚"},
						  { id:"不详", value:"不详"}
                          ]},
                        { view:"text", label:'证件号', id:"Cridentity" },
                    ]}
            ]
        },
		{
            rows:[
                {
                    cols:[
                        { view:"text", label:'口味爱好', id:"Crhobby" },
                        { view:"text", label:'备注', id:"Crmemo" },
                    ]}
            ]
        },
        { view:"button", value: "保存",id:"btnsave",click:function(){
			var b = this;
            var data = {
				uid: $$('Uid').getValue(),
				crname: $$('Crname').getValue(),
				kehlxid: $$('Kehlx').getValue(),
				crqcode: $$('Crqcode').getValue(),
				crtitle: $$('Crtitle').getValue(),
				storeid: getUrlParam("store"),
				crsex: $$('Crsex').getValue(),
				mobile: $$('Mobile').getValue(),
				crbirthday: $$('Crbirthday').getValue(),
				crmarriage: $$('Crmarriage').getValue(),
				crzip: $$('Crzip').getValue(),
				cridentity: $$('Cridentity').getValue(),
				crhobby: $$('Crhobby').getValue(),
				crmemo: $$('Crmemo').getValue()
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
        	if (data.kehlxid == "") {
	        	webix.message("请选择客户类型！");
	        	return
	        }
			if (data.crname == "") {
	        	webix.message("请输入客户姓名！");
	        	return
	        }
            ajax("POST", "/updcust", data, function (result) {
                if (result.success) {
					if ($$('uuid').getValue()=="") {
						data.uid = result.data;
						data.crbirthday = data.uid.split(",")[2] 
						data.crqcode = data.uid.split(",")[1] 
						data.uid = data.uid.split(",")[0]
					  	grida.add(data, 0)	
					} else {
					  data.crbirthday = result.data;	
					  grida.updateItem($$('uuid').getValue(), data)
					}
					b.getTopParentView().hide();
                } else {
                    webix.message(result.data);
                }
            })
        }},
        { view:"button", value: "退出", click:function(){
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
    head:"客户详情",
    body:webix.copy(form)
});
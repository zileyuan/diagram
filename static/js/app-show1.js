function do_card() {
    var data = {
        StartDate: $$('sd1').getValue(),
        FinishDate: $$('fd1').getValue(),
        KH: $$('kh').getValue(),
        Show0: $$('show0').getValue(),
        CardType: $$('cardtype').getValue(),
        CardPoint: $$('cardpoint').getValue(),
    };
    ajax("POST", "/cardtotal", data, function (result) {
        var data = result.data;
        $$('dtable').clearAll();
        $$('dtable').parse(data);
    })
}

function do_ajax() {
    var data = {
        StartDate: $$('sd').getValue(),
        FinishDate: $$('fd').getValue(),
    };
    ajax("POST", "/overview", data, function (result) {
        //webix.message(JSON.stringify(result));
        // 基于准备好的dom，初始化echarts实例
        var data = result.data;
        xAxisData = [];
        ss = [];
        xs = [];
        zk = [];
        zs = [];
        zj = [];
        rs = [];
        rj = [];
        for (i in data) {
            obj = data[i];
            xAxisData.push(obj.csname);
            ss.push(obj.samt)
            xs.push(obj.xsamt)
            zk.push(obj.rbtamt)
            zs.push(obj.zhuos)
            zj.push(obj.zhuoj)
            rs.push(obj.rens)
            rj.push(obj.renj)
        }
        var myChart = echarts.init(document.getElementById('main'));

// 指定图表的配置项和数据
        option = {
            tooltip: {
                trigger: 'axis'
            },
            toolbox: {
                show: true,
                feature: {
                    mark: {show: true},
                    dataView: {show: true, readOnly: false},
                    magicType: {show: true, type: ['line', 'bar']},
                    restore: {show: true},
                    saveAsImage: {show: true}
                }
            },
            calculable: true,
            legend: {
                data: ["实收金额", "销售金额", "折扣金额", "桌数", "桌均", "人数", "人均"],
            },
            xAxis: [
                {
                    type: 'category',
                    data: xAxisData,
                }
            ],
            yAxis: [
                {
                    type: 'value',
                    name: '金额',
                    min: 0,
                    max: 2500000,
                    interval: 500000,
                    axisLabel: {
                        formatter: '{value} 元'
                    }
                },
                {
                    type: 'value',
                    name: '人数',
                    min: 0,
                    max: 25000,
                    interval: 5000,
                    axisLabel: {
                        formatter: '{value}'
                    }
                }
            ],
            series: [

                {
                    name: '实收金额',
                    type: 'bar',
                    data: ss
                },
                {
                    name: '销售金额',
                    type: 'line',
                    yAxisIndex: 0,
                    data: xs
                },
                {
                    name: '折扣金额',
                    type: 'line',
                    yAxisIndex: 0,
                    data: zk
                },
                {
                    name: '桌数',
                    type: 'line',
                    yAxisIndex: 1,
                    data: zs
                },
                {
                    name: '桌均',
                    type: 'line',
                    yAxisIndex: 1,
                    data: zj
                },
                {
                    name: '人数',
                    type: 'line',
                    yAxisIndex: 1,
                    data: rs
                },
                {
                    name: '人均',
                    type: 'line',
                    yAxisIndex: 1,
                    data: rj
                }
            ]
        };

// 使用刚指定的配置项和数据显示图表。
        myChart.setOption(option);

    })
}


webix.ui({
    view: "toolbar",
    container: "toolbar",
    elements: [
        {
            view: "datepicker",
            id: "sd",
            align: "right",
            label: 'Start Date',
            labelWidth: 150,
            stringResult: true

        },
        {
            view: "datepicker",
            id: "fd",
            align: "right",
            label: 'Finish Date',
            labelWidth: 150,
            stringResult: true
        }
    ]
});

webix.ui({
    view: "toolbar",
    container: "toolbar2",
    elements: [
        {
            view: "datepicker",
            id: "sd1",
            align: "right",
            label: 'Start Date',
            labelWidth: 150,
            stringResult: true

        },
        {
            view: "datepicker",
            id: "fd1",
            align: "right",
            label: 'Finish Date',
            labelWidth: 150,
            stringResult: true
        },
        {
            view:"text",
            id: "kh",
            label:"卡号"
        },
        {
            view:"checkbox",
            label:"不显示全零金额",
            id: "show0",
            uncheckValue: 0,
            checkValue: 1,
        },
        {
            view:"multiselect",
            label:"卡类型",
            id: "cardtype",
            labelWidth:100,
            options:[
                { id:1, value:"金卡" },
                { id:2, value:"银卡" },
                { id:3, value:"通卡" },
                { id:4, value:"贴卡" }
            ],
        },
        {
            view: "combo",
            id:"cardpoint",
            width:400,
            label: '发卡店',
            labelWidth:220,
            options:[
                { id:"001", value:"一店"   },
                { id:"002", value:"二店"   },
                { id:"003", value:"三店" }
            ]
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
        { id:"huiyk_id",	header:"卡号", footer:[
            { height:40, text:"xiaoji", colspan:3 },
            { height:40, text:"xiaojiavg", colspan:3 },
            { height:40, text:"heji", colspan:3 },
            { height:40, text:"hejiavg", colspan:3 }
        ]},
        { id:"huiyk_germc",	header:"持卡人"},
        { id:"name",	header:"卡类型"},
        { id:"chuz",	header:"本期储值", footer:[
            { content:"pageSummColumn" },
            { content:"pageAvgColumn" },
            { content:"summColumn" },
            { content:"avgColumn" },
        ]},
        { id:"xiaof",	header:"本期消费"},
        { id:"zxye",	header:"期初金额"},
        { id:"zdye",	header:"本期余额"},
        { id:"zzye",	header:"最终余额"},
        { id:"crmobile",	header:"手机号"},
        { id:"huiyk_zhuangt",	header:"卡状态"},
        { id:"huiyk_fakrq",	header:"发卡日期"},
        { id:"huiyk_jiezrq",	header:"截止日期"}
    ],
    select:"cell",
    autowidth:true,
    footer:true,
    pager:{
        template:"{common.first()} {common.prev()} {common.pages()} {common.next()} {common.last()}",
        container:"paging_here",
        size:100,
        group:5
    },
});
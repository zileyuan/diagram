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


webix.ui(
    {
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


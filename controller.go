package main

import (
	"fmt"

	"github.com/axgle/mahonia"
	"gopkg.in/macaron.v1"
)

type ContextResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func DoIndex(ctx *macaron.Context) {
	ctx.HTML(200, "show1")
}

func DoOverview(ctx *macaron.Context) {
	sql := `select cscode,csname,a.xsamt,-c.rbtamt rbtamt,b.samt,d.zhuos,d.rens,a.xsamt/d.zhuos zhuoj,a.xsamt/d.rens renj from cmstore inner join (select cb.bill_storecode,sum(isnull(ci.bci_debit,0)) xsamt from cmbill cb left join cmcashitem ci on cb.bill_uid=ci.bci_billuid left join bizaccount on bizaccount.bizaccount_huid=bci_acchuid where bill_bizdate between '20120501' and '20120620' and bizaccount.actype='104' group by cb.bill_storecode)a on cscode=a.bill_storecode left join(select cb.bill_storecode,sum(isnull(ci.bci_credit,0)) samt from cmbill cb left join cmcashitem ci on cb.bill_uid=ci.bci_billuid where bill_bizdate between '20120501' and '20120620' group by cb.bill_storecode )b on  cscode=b.bill_storecode left join(select cb.bill_storecode,case when isnull(sum(ci.bci_debit),0)=0 then sum(isnull(ci.bci_credit,0)) else isnull(sum(ci.bci_debit),0) end as rbtamt from cmbill cb inner join cmcashitem ci on cb.bill_uid=ci.bci_billuid where bill_bizdate between '20120501' and '20120620' and bci_accuid in (00000015,00000027,00000032) group by cb.bill_storecode)c on cscode=c.bill_storecode left join (SELECT bill_storecode,SUM(bill_tablecount) zhuos,SUM(bill_personcount) rens FROM cmBill where bill_amount>0 and bill_bizdate between '20120501' and '20120620' GROUP BY bill_storecode)d on cscode=d.bill_storecode where 1=1 order by cscode`
	rows, err := AppDB.Query(sql)
	defer rows.Close()
	ovs := []Overview{}
	for rows.Next() {
		var ov Overview
		err = rows.ScanStructByName(&ov)
		if err == nil {
			ov.Csname = GBKToUtf8(ov.Csname)
			ovs = append(ovs, ov)
		}
	}
	fmt.Printf("999999999 %#v", ovs)

	if err == nil {
		ctx.JSON(200, &ContextResult{
			Success: true,
			Data:    ovs,
		})
	} else {
		ctx.JSON(200, &ContextResult{
			Success: false,
		})
	}
}

//gbk转utf-8
func GBKToUtf8(str string) string {
	//字符集转换
	enc := mahonia.NewDecoder("gbk")
	return enc.ConvertString(str)
}

package app

import (
	"fmt"
	"time"

	"github.com/axgle/mahonia"
	"gopkg.in/macaron.v1"
)

const (
	TimeLayout = "2006-01-02 15:04"
	DateLayout = "2006-01-02"
)

type ContextResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func DoIndex(ctx *macaron.Context) {
	ctx.HTML(200, "show1")
}

func DoOverview(ctx *macaron.Context) {
	startDateStr := ctx.Query("StartDate")
	finishDateStr := ctx.Query("FinishDate")
	fmt.Println(startDateStr, finishDateStr)

	startDate := String2Time(TimeLayout, startDateStr)
	finishDate := String2Time(TimeLayout, finishDateStr)

	startDateStr = fmt.Sprintf("%04d%02d%02d", startDate.Year(), startDate.Month(), startDate.Day())
	finishDateStr = fmt.Sprintf("%04d%02d%02d", finishDate.Year(), finishDate.Month(), finishDate.Day())
	fmt.Println(startDateStr, finishDateStr)

	sql := `select cscode,csname,a.xsamt,-c.rbtamt rbtamt,b.samt,d.zhuos,d.rens,a.xsamt/d.zhuos zhuoj,a.xsamt/d.rens renj from cmstore inner join (select cb.bill_storecode,sum(isnull(ci.bci_debit,0)) xsamt from cmbill cb left join cmcashitem ci on cb.bill_uid=ci.bci_billuid left join bizaccount on bizaccount.bizaccount_huid=bci_acchuid where bill_bizdate between '%s' and '%s' and bizaccount.actype='104' group by cb.bill_storecode)a on cscode=a.bill_storecode left join(select cb.bill_storecode,sum(isnull(ci.bci_credit,0)) samt from cmbill cb left join cmcashitem ci on cb.bill_uid=ci.bci_billuid where bill_bizdate between '%s' and '%s' group by cb.bill_storecode )b on  cscode=b.bill_storecode left join(select cb.bill_storecode,case when isnull(sum(ci.bci_debit),0)=0 then sum(isnull(ci.bci_credit,0)) else isnull(sum(ci.bci_debit),0) end as rbtamt from cmbill cb inner join cmcashitem ci on cb.bill_uid=ci.bci_billuid where bill_bizdate between '%s' and '%s' and bci_accuid in (00000015,00000027,00000032) group by cb.bill_storecode)c on cscode=c.bill_storecode left join (SELECT bill_storecode,SUM(bill_tablecount) zhuos,SUM(bill_personcount) rens FROM cmBill where bill_amount>0 and bill_bizdate between '%s' and '%s' GROUP BY bill_storecode)d on cscode=d.bill_storecode where 1=1 order by cscode`
	sql = fmt.Sprintf(sql, startDateStr, finishDateStr, startDateStr, finishDateStr, startDateStr, finishDateStr, startDateStr, finishDateStr)
	fmt.Println(sql)
	rows, err := AppDB.Query(sql)
	fmt.Printf("Query %v", err)
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

//字符串转换成时间
func String2Time(layout string, timeStr string) time.Time {
	tm, _ := time.Parse(layout, timeStr)
	return tm
}
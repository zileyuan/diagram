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

func DoCardTotal(ctx *macaron.Context) {
	startDateStr := ctx.Query("StartDate")
	finishDateStr := ctx.Query("FinishDate")
	KH := ctx.Query("KH")
	Show0 := ctx.QueryInt("Show0")
	CardType := ctx.Query("CardType")
	CardPoint := ctx.Query("CardPoint")
	fmt.Printf("%v %v %v %v %v %v", startDateStr, finishDateStr, KH, Show0, CardType, CardPoint)

	sql := `select huiyk_id,case when (crname='' or crname is null) then huiyk_germc else crname end khmc,sscate.name,isnull(x.chuz,0) chuz,isnull(y.xiaof,0) xiaof,zxye,zdye,zzye,crmobile
       ,case huiyk_zhuangt when '01' then '启用' when '04' then '挂失' when '03' then '作废' end as 卡状态,huiyk_fakrq,huiyk_jiezrq
from huiyk left join bizfolio on bizfolio.account=huiyk_danwid
left join crcustomer on huiyk_gerid=crcustomer.uid
left join sscate on huiyk_leixid=sscate.uid
left join
(
select bizfolio.account,isnull(sum(Credit),0) as chuz
from bizfolio left join huiyk on huiyk_danwid=bizfolio.account
where bizfolio.Type='NCCard' and CONVERT(char(10), bizfolio.sysdate, 112)  between '@bdate@' and '@edate@'
IIF{xians=0,"and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{xians=1,"and huiyk_zhuangt='01'"}
IIF{xians=2,"and huiyk_zhuangt='04'"}
IIF{xians=3,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)>1) and huiyk_zhuangt='03'"}
IIF{xians=4,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)=1) and huiyk_zhuangt='03'"}
IIF{xians=5,"and huiyk_jiezrq<@edate@ and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{storeid, "and left(huiyk_danwmc,4)='@storeid@'"}
IIF{caozy,"and opid='@caozy@'"}
group by bizfolio.Account ) x on bizfolio.account=x.account
left join
(
select bizfolio.account,isnull(sum(Debit),0) as xiaof
from bizfolio  left join huiyk on huiyk_danwid=bizfolio.account
where bizfolio.Type='NCCard' and CONVERT(char(10), bizfolio.sysdate, 112) between '@bdate@' and '@edate@'
IIF{xians=0,"and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{xians=1,"and huiyk_zhuangt='01'"}
IIF{xians=2,"and huiyk_zhuangt='04'"}
IIF{xians=3,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)>1) and huiyk_zhuangt='03'"}
IIF{xians=4,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)=1) and huiyk_zhuangt='03'"}
IIF{xians=5,"and huiyk_jiezrq<@edate@ and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{storeid, "and left(huiyk_danwmc,4)='@storeid@'"}
IIF{caozy,"and opid='@caozy@'"}
group by bizfolio.account
)y on bizfolio.account=y.account
left join
(select sum(zxye) zxye,sum(zdye) zdye,sum(zzye) zzye,account from
(
select sum(balance) zxye,0 zdye,0 zzye,account from bizfolio left join huiyk on bizfolio.account=huiyk_danwid
where bizfolio.Type='NCCard' and bizfolio.uid in (select max(uid) as zxuid from bizfolio where CONVERT(char(10),bizfolio.sysdate,112)<'@bdate@' and Type='NCCard' group by account)
IIF{xians=0,"and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{xians=1,"and huiyk_zhuangt='01'"}
IIF{xians=2,"and huiyk_zhuangt='04'"}
IIF{xians=3,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)>1) and huiyk_zhuangt='03'"}
IIF{xians=4,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)=1) and huiyk_zhuangt='03'"}
IIF{xians=5,"and huiyk_jiezrq<@edate@ and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{storeid, "and left(huiyk_danwmc,4)='@storeid@'"}
group by account
union
select 0 zxye,sum(balance) zdye,0 zzye,account
from bizfolio left join huiyk on bizfolio.account=huiyk_danwid
where Type='NCCard' and bizfolio.uid in (select max(uid) as zduid from bizfolio where CONVERT(char(10),bizfolio.sysdate,112)<='@edate@' and Type='NCCard' group by account)
IIF{xians=0,"and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{xians=1,"and huiyk_zhuangt='01'"}
IIF{xians=2,"and huiyk_zhuangt='04'"}
IIF{xians=3,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)>1) and huiyk_zhuangt='03'"}
IIF{xians=4,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)=1) and huiyk_zhuangt='03'"}
IIF{xians=5,"and huiyk_jiezrq<@edate@ and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{storeid, "and left(huiyk_danwmc,4)='@storeid@'"}
group by account
union
select 0 zxye,0 zdye,sum(balance) zzye,account
from bizfolio left join huiyk on bizfolio.account=huiyk_danwid
where Type='NCCard' and bizfolio.uid in (select max(uid) as zduid from bizfolio where Type='NCCard'  group by account)
IIF{xians=0,"and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{xians=1,"and huiyk_zhuangt='01'"}
IIF{xians=2,"and huiyk_zhuangt='04'"}
IIF{xians=3,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)>1) and huiyk_zhuangt='03'"}
IIF{xians=4,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)=1) and huiyk_zhuangt='03'"}
IIF{xians=5,"and huiyk_jiezrq<@edate@ and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{storeid, "and left(huiyk_danwmc,4)='@storeid@'"}
group by account
)a
group by account
)n on bizfolio.account=n.account
where 1=1
IIF{xdate,"and huiyk_fakrq>='@xdate@'"}
IIF{ydate,"and huiyk_fakrq<='@ydate@'"}
IIF{kehjl, "and crcustomer.crmail3='@kehjl@'"}
IIF{cardid,"and huiyk.huiyk_id like '%@cardid@%'"}
IIF{cardman,"and huiyk_germc like '%@cardman@%'"}
IIF{klx,"and huiyk_leixid='@klx@'"}
IIF{storeid, "and left(huiyk_danwmc,4)='@storeid@'"}
IIF{itemid=0,"and (zxye<>0 or zdye<>0 or zzye<>0 or x.chuz<>0 or y.xiaof<>0)"}
IIF{xians=0,"and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
IIF{xians=1,"and huiyk_zhuangt='01'"}
IIF{xians=2,"and huiyk_zhuangt='04'"}
IIF{xians=3,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)>1) and huiyk_zhuangt='03'"}
IIF{xians=4,"and huiyk_danwid in (select huiyk_danwid from huiyk group by huiyk_danwid having count(huiyk_danwid)=1) and huiyk_zhuangt='03'"}
IIF{xians=5,"and huiyk_jiezrq<@edate@ and (huiyk_zhuangt='01' or huiyk_zhuangt='04')"}
group by huiyk_id,crname,huiyk_germc,x.chuz,y.xiaof,zxye,zdye,zzye,sscate.name,sscate.name,huiyk_zhuangt,crmobile,huiyk_fakrq,huiyk_jiezrq
order by huiyk_id`
	fmt.Println(sql)
	rows, err := AppDB.Query(sql)
	if err != nil {
		cts := []CardTotal{}
		for i := 0; i < 10000; i++ {
			ct := CardTotal{
				Huiyk_id:      "001",
				Huiyk_germc:   "咱梦" + fmt.Sprintf("%05d", i),
				Name:          "储值卡",
				Chuz:          100.0,
				Xiaof:         120.0,
				Zxye:          100.0,
				Zdye:          200.0,
				Zzye:          300.0,
				Crmobile:      "1390000",
				Huiyk_zhuangt: "启用",
				Huiyk_fakrq:   "20160202",
				Huiyk_jiezrq:  "20170202",
			}
			cts = append(cts, ct)
		}
		ctx.JSON(200, &ContextResult{
			Success: true,
			Data:    cts,
		})
		return
	}
	fmt.Printf("Query %v", err)
	defer rows.Close()
	cts := []CardTotal{}
	for rows.Next() {
		var ct CardTotal
		err = rows.ScanStructByName(&ct)
		if err == nil {
			ct.Huiyk_germc = GBKToUtf8(ct.Huiyk_germc)
			ct.Name = GBKToUtf8(ct.Name)
			ct.Huiyk_zhuangt = GBKToUtf8(ct.Huiyk_zhuangt)
			cts = append(cts, ct)
		}
	}
	fmt.Printf("999999999 %#v", cts)

	if err == nil {
		ctx.JSON(200, &ContextResult{
			Success: true,
			Data:    cts,
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

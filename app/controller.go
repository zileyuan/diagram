package app

import (
	"fmt"
	"strings"
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

func DoCardType(ctx *macaron.Context) {
	sql := `select uid as id,name as value from sscate where type='CardType'`
	fmt.Println(sql)
	rows, err := AppDB.Query(sql)
	defer rows.Close()
	cts := []Cardtype{}
	for rows.Next() {
		var ct Cardtype
		//fmt.Println("------------111111111111------------------------")
		err = rows.ScanStructByName(&ct)
		//fmt.Println("------------222222222222222------------------------")
		if err == nil {
			ct.Value = GBKToUtf8(ct.Value)
			cts = append(cts, ct)
		}
	}
	ctx.JSON(200, &ContextResult{
		Success: true,
		Data:    cts,
	})
}

func DoCardTotal(ctx *macaron.Context) {
	startDateStr := ctx.Query("StartDate")
	finishDateStr := ctx.Query("FinishDate")
	KH := ctx.Query("KH") + "%"
	Show0 := ctx.QueryInt("Show0")
	CardType := ctx.Query("CardType")
	strs := strings.Split(CardType, ",")
	var sqltext string
	for _, j := range strs {
		sqltext = sqltext + ` or huiyk_leixid='` + j + `'`
	}
	sqltext = strings.Replace(sqltext, `or`, `where (`, 1) + ")"
	fmt.Println(sqltext)

	//startDate := String2Time(TimeLayout, startDateStr)
	//finishDate := String2Time(TimeLayout, finishDateStr)

	//startDateStr = fmt.Sprintf("%04d-%02d-%02d", startDate.Year(), startDate.Month(), startDate.Day())
	//finishDateStr = fmt.Sprintf("%04d-%02d-%02d", finishDate.Year(), finishDate.Month(), finishDate.Day())
	//fmt.Println(startDateStr, finishDateStr)
	CardPoint := ctx.Query("CardPoint")
	fmt.Printf("%v %v %v %v %v %v", startDateStr, finishDateStr, KH, Show0, CardType, CardPoint)
	sql := `select huiyk_id as huiykid, crname,cardtype,credit,debit,balance,acbalance,
	crmobile,huiyk_zhuangt as huiykzhuangt,huiyk_fakrq as huiykfakrq,
	huiyk_jiezrq as huiykjiezrq from (SELECT HUIYK.huiyk_id,crCustomer.crName,
	ssCate.Name as cardtype,SUM(BizFolio.Credit) AS credit,
	SUM(BizFolio.Debit) AS debit,crAccount.acBalance,crCustomer.crMobile,
	HUIYK.huiyk_zhuangt,HUIYK.huiyk_fakrq, HUIYK.huiyk_jiezrq,crAccount.UID 
	FROM HUIYK INNER JOIN crAccount ON HUIYK.huiyk_danwid = crAccount.UID 
	INNER JOIN BizFolio ON crAccount.UID = BizFolio.Account INNER JOIN ssCate 
	ON HUIYK.huiyk_leixid = ssCate.UID LEFT OUTER JOIN crCustomer ON 
	HUIYK.huiyk_Gerid = crCustomer.UID %s and sysdate between '%s' and '%s' and huiyk_id like '%s' GROUP BY HUIYK.huiyk_id, crCustomer.crName,
	ssCate.Name,crAccount.acBalance, crCustomer.crMobile, HUIYK.huiyk_zhuangt,
	HUIYK.huiyk_fakrq, HUIYK.huiyk_jiezrq,crAccount.UID) a,(SELECT BizFolio.Account,
	BizFolio.Balance FROM BizFolio INNER JOIN (SELECT MAX(UID) AS uid, Account 
	FROM BizFolio GROUP BY Account) m ON BizFolio.UID = m.uid) b where a.uid=b.account`
	//sql := `select huiyk_id`
	//sql := `select huiyk_id as huiykid from huiyk`
	sql = fmt.Sprintf(sql, sqltext, startDateStr, finishDateStr, KH)
	fmt.Println(sql)
	rows, err := AppDB.Query(sql)
	if err != nil {
		//		cts := []CardTotal{}
		//		for i := 0; i < 10000; i++ {
		//			ct := CardTotal{
		//				Huiykid: "001",
		//				//Crname:        "咱梦" + fmt.Sprintf("%05d", i),
		//				//Cardtype:  "储值卡",
		//				//Credit:    100.0,
		//				//Debit:     120.0,
		//				//Balance:   200.0,
		//				//Acbalance: 300.0,
		//				//Crmobile:      "1390000",
		//				//Huiyk_zhuangt: "启用",
		//				//Huiyk_fakrq:   "20160202",
		//				//Huiyk_jiezrq:  "20170202",
		//			}
		//			cts = append(cts, ct)
		//		}
		//		ctx.JSON(200, &ContextResult{
		//			Success: true,
		//			Data:    cts,
		//		})
		//		return
	}
	fmt.Printf("Query %v", err)
	fmt.Println("------------------------------------")
	defer rows.Close()
	cts := []CardTotal{}
	for rows.Next() {
		var ct CardTotal
		//fmt.Println("------------111111111111------------------------")
		err = rows.ScanStructByName(&ct)
		//fmt.Println("------------222222222222222------------------------")

		if err == nil {
			ct.Huiykid = GBKToUtf8(ct.Huiykid)
			ct.Crname = GBKToUtf8(ct.Crname)
			ct.Cardtype = GBKToUtf8(ct.Cardtype)
			ct.Huiykzhuangt = GBKToUtf8(ct.Huiykzhuangt)
			cts = append(cts, ct)
		} else {
			fmt.Printf("**********%v*********", err)
		}
	}
	fmt.Printf("999999999 %#v", cts)
	fmt.Printf("err %#v", err)
	if true {
		ctx.JSON(200, &ContextResult{
			Success: true,
			Data:    cts,
		})
	} else {
		ctx.JSON(200, &ContextResult{
			Success: false,
			Data:    "error",
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

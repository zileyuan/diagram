package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/mozillazg/go-pinyin"
)

const (
	TimeLayout = "2006-01-02 15:04"
	DateLayout = "2006-01-02"
)

type ContextResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func NeedSignedIn(ctx *Context) {
	if !ctx.IsSigned {
		ctx.Redirect("/")
	} else {
		ctx.Next()
	}
}

func DoIndex(ctx *Context) {
	ctx.HTML(200, "login")
}

func DoPage1(ctx *Context) {
	ctx.HTML(200, "show1")
}

func DoPage2(ctx *Context) {
	ctx.HTML(200, "show2")
}

func OnLogin(ctx *Context) {
	username := ctx.Query("username")
	userpass := ctx.Query("userpass")
	fmt.Println(username, userpass)
	if username == "root" && userpass == "pass" {
		fmt.Println("login ok!")
		ctx.Session.Set(SessionKey, username)
		ctx.JSON(200, &ContextResult{
			Success: true,
			Data: "/NetCust",
		})
		return
	}
	ctx.JSON(200, &ContextResult{
		Success: false,
	})
}

func OnOverview(ctx *Context) {
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

func OnCustType(ctx *Context) {
	sql := `select yudlx_id as id,yudlx_mingc as value from yudlx where itype=1`
	fmt.Println(sql)
	rows, err := AppDB.Query(sql)
	defer rows.Close()
	cts := []Custtype{}
	for rows.Next() {
		var ct Custtype
		err = rows.ScanStructByName(&ct)
		if err == nil {
			ct.Id = GBKToUtf8(ct.Id)
			ct.Value = GBKToUtf8(ct.Value)
			cts = append(cts, ct)
		}
	}
	ctx.JSON(200, &ContextResult{
		Success: true,
		Data:    cts,
	})
}

func OnStoreCode(ctx *Context) {
	sql := `select csCode as id,csName as value from cmstore order by csCode`
	fmt.Println(sql)
	rows, err := AppDB.Query(sql)
	defer rows.Close()
	cts := []Storecode{}
	for rows.Next() {
		var ct Storecode
		//fmt.Println("------------111111111111------------------------")
		err = rows.ScanStructByName(&ct)
		//fmt.Println("------------222222222222222------------------------")
		if err == nil {
			ct.Id = GBKToUtf8(ct.Id)
			ct.Value = GBKToUtf8(ct.Value)
			cts = append(cts, ct)
		}
	}
	ctx.JSON(200, &ContextResult{
		Success: true,
		Data:    cts,
	})
}

func OnCardType(ctx *Context) {
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

func OnCustomer(ctx *Context) {
	StoreCode := ctx.Query("StoreCode")
	CustType := ctx.Query("CustType")
	CustName := ctx.Query("CustName")
	Mobile := ctx.Query("Mobile")
	var sqltext, sqltext2, sqltext3, sqltext4 string
	if StoreCode == "" {
		sqltext = "where 1=1"
	} else {
		sqltext = `where crID='` + StoreCode + `'`
	}
	if CustType != "" {
		sqltext2 = `and crKehlx='` + CustType + `'`
	}
	if CustName != "" {
		sqltext3 = `and (crName like '%` + CustName + `%' or crQCode like '` + CustName + `%')`
	}
	if Mobile != "" {
		sqltext4 = `and crMobile like '` + Mobile + `%'`
	}
	fmt.Printf("%v %v %v %v", StoreCode, CustType, CustName, Mobile)
	sql := `SELECT crCustomer.uid,isnull(crID,'') as storeid,isnull(csName,'') as store,
	isnull(crname,'') as crname,isnull(crqcode,'') as crqcode,isnull(crtitle,'') as crtitle,
	isnull(crKehlx,'') as kehlxid,isnull(yudlx_mingc,'') as kehlx,isnull(crsex,'') as crsex,
	isnull(crmobile,'') as mobile,isnull(crbirthday,'') as crbirthday,
	isnull(crZip,'') as crzip,isnull(crIdentity,'') as cridentity,
	isnull(crMarriage,'') as crmarriage,isnull(crMemo,'') as crmemo,
	isnull(crHobby,'') as crhobby 
	FROM crCustomer INNER JOIN crPerson ON crCustomer.UID = crPerson.UID 
	LEFT OUTER JOIN YUDLX ON crCustomer.crKehlx = YUDLX.YUDLX_ID 
	LEFT OUTER JOIN cmStore ON crCustomer.crID = cmStore.csCode %s %s %s %s`
	sql = fmt.Sprintf(sql, sqltext, sqltext2, sqltext3, sqltext4)
	fmt.Println(sql)
	rows, err := AppDB.Query(sql)
	if err != nil {
		fmt.Printf("Query %v", err)
	}
	defer rows.Close()
	cts := []Customer{}
	for rows.Next() {
		var ct Customer
		err = rows.ScanStructByName(&ct)

		if err == nil {
			ct.Uid = GBKToUtf8(ct.Uid)
			ct.Store = GBKToUtf8(ct.Store)
			ct.Crname = GBKToUtf8(ct.Crname)
			ct.Crtitle = GBKToUtf8(ct.Crtitle)
			ct.Kehlx = GBKToUtf8(ct.Kehlx)
			ct.Crsex = GBKToUtf8(ct.Crsex)
			ct.Mobile = GBKToUtf8(ct.Mobile)
			ct.Crbirthday = GBKToUtf8(ct.Crbirthday)
			ct.Crmarriage = GBKToUtf8(ct.Crmarriage)
			ct.Crzip = GBKToUtf8(ct.Crzip)
			ct.Cridentity = GBKToUtf8(ct.Cridentity)
			ct.Crhobby = GBKToUtf8(ct.Crhobby)
			ct.Crmemo = GBKToUtf8(ct.Crmemo)
			cts = append(cts, ct)
		} else {
			fmt.Printf("**********%v*********", err)
		}
	}
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

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func OnUpdCust(ctx *Context) {
	Uid := ctx.Query("uid")
	Crname := ctx.Query("crname")
	Kehlxid := ctx.Query("kehlxid")
	Crqcode := ctx.Query("crqcode")
	Crtitle := ctx.Query("crtitle")
	Storeid := ctx.Query("storeid")
	Crsex := ctx.Query("crsex")
	Mobile := ctx.Query("mobile")
	Crbirthday := ctx.Query("crbirthday")
	Crbirthday = Substr(Crbirthday, 0, 10)
	Crmarriage := ctx.Query("crmarriage")
	Crzip := ctx.Query("crzip")
	Cridentity := ctx.Query("cridentity")
	Crhobby := ctx.Query("crhobby")
	Crmemo := ctx.Query("crmemo")
	//fmt.Printf("birthday %v", Crbirthday)
	//myDate, _ := time.Parse("2006-01-02", Crbirthday)
	//Crbirthday = myDate.Format("2006-01-02")
	//myDate := String2Time(TimeLayout, Crbirthday)
	//Crbirthday = fmt.Sprintf("%04d-%02d-%02d", myDate.Year(), myDate.Month(), myDate.Day())
	if Crqcode == "" {
		py := pinyin.NewArgs()
		strs := pinyin.Pinyin(Crname, py)
		//strs := py.Convert(Crname)
		//fmt.Println(strs)
		for _, v := range strs {
			Crqcode += string(v[0][0])
		}
		Crqcode = strings.ToUpper(Crqcode)
	}
	if Uid == "" {
		sql := `select cast(max(uid) as int) as uid from crcustomer`
		rows, err := AppDB.Query(sql)
		if err != nil {
			fmt.Printf("Query %v", err)
		}
		defer rows.Close()
		var ct Customer
		for rows.Next() {
			err = rows.ScanStructByName(&ct)
			if err == nil {
				val, err := strconv.Atoi(ct.Uid)
				val++
				newid := fmt.Sprintf("%012d", val)
				sql := `insert into crcustomer (uid,crID,crType,crname,crqcode,crkehlx,crMobile,crZip,crMemo) values ('%s','%s',1,'%s','%s','%s','%s','%s','%s')`
				sql = fmt.Sprintf(sql, newid, Storeid, Crname, Crqcode, Kehlxid, Mobile, Crzip, Crmemo)
				fmt.Println(sql)
				_, err = AppDB.Exec(sql)
				if err == nil {
					sql := `insert into crperson (uid,crTitle,crSex,crBirthday,crMarriage,crIdentity,crHobby) values ('%s','%s','%s','%s','%s','%s','%s')`
					sql = fmt.Sprintf(sql, newid, Crtitle, Crsex, Crbirthday, Crmarriage, Cridentity, Crhobby)
					fmt.Println(sql)
					_, err = AppDB.Exec(sql)
					if err != nil {
						fmt.Printf("QueryErr %v", sql)
					} else {
						ctx.JSON(200, &ContextResult{Success: true, Data: newid + "," + Crqcode + "," + Crbirthday})

					}
				}
			} else {
				fmt.Printf("**********%v*********", err)
			}
		}
	} else {
		sql := `update crcustomer set crname='%s',crQCode='%s',crKehlx='%s',crMobile='%s',crZip='%s',crMemo='%s' where uid='%s'`
		sql = fmt.Sprintf(sql, Crname, Crqcode, Kehlxid, Mobile, Crzip, Crmemo, Uid)
		fmt.Println(sql)
		_, err := AppDB.Exec(sql)
		if err == nil {
			sql := `update crperson set crTitle='%s',crSex='%s',crBirthday='%s',crMarriage='%s',crIdentity='%s',crHobby='%s' where uid='%s'`
			sql = fmt.Sprintf(sql, Crtitle, Crsex, Crbirthday, Crmarriage, Cridentity, Crhobby, Uid)
			fmt.Println(sql)
			_, err = AppDB.Exec(sql)
			if err != nil {
				fmt.Printf("QueryErr %v", sql)
			} else {
				ctx.JSON(200, &ContextResult{Success: true, Data: Crbirthday})

			}
		}
		//		if err != nil {
		//			fmt.Printf("QueryErr %v", sql)
		//		} else {
		//			ctx.JSON(200, &ContextResult{Success: true, Data: ""})

		//		}
	}
}

func OnCardTotal(ctx *Context) {
	startDateStr := ctx.Query("StartDate")
	finishDateStr := ctx.Query("FinishDate")
	KH := ctx.Query("KH")
	//KH := ctx.Query("KH") + "%"
	//Show0 := ctx.QueryInt("Show0")
	CardType := ctx.Query("CardType")
	var sqltext, sqltext2, sqltext3 string
	if CardType == "" {
		sqltext = "where 1=1"
	} else {
		strs := strings.Split(CardType, ",")
		for _, j := range strs {
			sqltext = sqltext + ` or huiyk_leixid='` + j + `'`
		}
		sqltext = strings.Replace(sqltext, `or`, `where (`, 1) + ")"
		fmt.Println(sqltext)
	}
	//startDate := String2Time(TimeLayout, startDateStr)
	//finishDate := String2Time(TimeLayout, finishDateStr)

	//startDateStr = fmt.Sprintf("%04d-%02d-%02d", startDate.Year(), startDate.Month(), startDate.Day())
	//finishDateStr = fmt.Sprintf("%04d-%02d-%02d", finishDate.Year(), finishDate.Month(), finishDate.Day())
	//fmt.Println(startDateStr, finishDateStr)
	CardPoint := ctx.Query("CardPoint")
	if CardPoint != "" {
		sqltext2 = `and left(huiyk_danwmc,4)='` + CardPoint + `'`
	}
	if KH != "" {
		sqltext3 = `and huiyk_id like '` + KH + `%'`
	}

	fmt.Printf("%v %v %v %v %v %v", startDateStr, finishDateStr, KH, CardType, CardPoint)
	sql := `select huiyk_id as huiykid, isnull(crname,'') as crname,cardtype,credit,debit,balance,acbalance,
	isnull(crmobile,'') as crmobile,case huiyk_zhuangt when '01' then '启用' when '04' then '挂失' when '03' then '作废' end as huiykzhuangt,huiyk_fakrq as huiykfakrq,
	huiyk_jiezrq as huiykjiezrq from (SELECT HUIYK.huiyk_id,crCustomer.crName,
	ssCate.Name as cardtype,SUM(BizFolio.Credit) AS credit,
	SUM(BizFolio.Debit) AS debit,crAccount.acBalance,crCustomer.crMobile,
	HUIYK.huiyk_zhuangt,HUIYK.huiyk_fakrq, HUIYK.huiyk_jiezrq,crAccount.UID 
	FROM HUIYK INNER JOIN crAccount ON HUIYK.huiyk_danwid = crAccount.UID 
	INNER JOIN BizFolio ON crAccount.UID = BizFolio.Account INNER JOIN ssCate 
	ON HUIYK.huiyk_leixid = ssCate.UID LEFT OUTER JOIN crCustomer ON 
	HUIYK.huiyk_Gerid = crCustomer.UID %s and sysdate between '%s' and '%s' %s %s GROUP BY HUIYK.huiyk_id, crCustomer.crName,
	ssCate.Name,crAccount.acBalance, crCustomer.crMobile, HUIYK.huiyk_zhuangt,
	HUIYK.huiyk_fakrq, HUIYK.huiyk_jiezrq,crAccount.UID) a,(SELECT BizFolio.Account,
	BizFolio.Balance FROM BizFolio INNER JOIN (SELECT MAX(UID) AS uid, Account 
	FROM BizFolio GROUP BY Account) m ON BizFolio.UID = m.uid) b where a.uid=b.account`
	//sql := `select huiyk_id`
	//sql := `select huiyk_id as huiykid from huiyk`
	sql = fmt.Sprintf(sql, sqltext, startDateStr, finishDateStr, sqltext2, sqltext3)
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
	//fmt.Printf("999999999 %#v", cts)
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

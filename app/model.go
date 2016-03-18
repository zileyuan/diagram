package app

type Overview struct {
	Cscode string  `json:"cscode"`
	Csname string  `json:"csname"`
	Rbtamt float64 `json:"rbtamt"`
	Renj   float64 `json:"renj"`
	Rens   int     `json:"rens"`
	Samt   float64 `json:"samt"`
	Xsamt  float64 `json:"xsamt"`
	Zhuoj  float64 `json:"zhuoj"`
	Zhuos  int     `json:"zhuos"`
}

type CardTotal struct {
	Huiykid      string  `json:"huiykid"`      //卡号
	Crname       string  `json:"crname"`       //持卡人
	Cardtype     string  `json:"cardtype"`     //卡类型
	Credit       float64 `json:"credit"`       //本期储值
	Debit        float64 `json:"debit"`        //本期消费
	Balance      float64 `json:"balance"`      //本期余额
	Acbalance    float64 `json:"acbalance"`    //最终余额
	Crmobile     string  `json:"crmobile"`     //手机号
	Huiykzhuangt string  `json:"huiykzhuangt"` //卡状态
	Huiykfakrq   string  `json:"huiykfakrq"`   //发卡日期
	Huiykjiezrq  string  `json:"huiykjiezrq"`  //截止日期
}

type Cardtype struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type Storecode struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type Custtype struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type Customer struct {
	Uid        string `json:"uid"`
	Store      string `json:"store"`
	Crname     string `json:"crname"`
	Crtitle    string `json:"crtitle"`
	Kehlx      string `json:"kehlx"`
	Crsex      string `json:"crsex"`
	Mobile     string `json:"mobile"`
	Crbirthday string `json:"crbirthday"`
}

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
	Huiyk_id      string  `json:"huiyk_id"`      //卡号
	Huiyk_germc   string  `json:"huiyk_germc"`   //持卡人
	Name          string  `json:"name"`          //卡类型
	Chuz          float64 `json:"chuz"`          //本期储值
	Xiaof         float64 `json:"xiaof"`         //本期消费
	Zxye          float64 `json:"zxye"`          //期初金额
	Zdye          float64 `json:"zdye"`          //本期余额
	Zzye          float64 `json:"zzye"`          //最终余额
	Crmobile      string  `json:"crmobile"`      //手机号
	Huiyk_zhuangt string  `json:"huiyk_zhuangt"` //卡状态
	Huiyk_fakrq   string  `json:"huiyk_fakrq"`   //发卡日期
	Huiyk_jiezrq  string  `json:"huiyk_jiezrq"`  //截止日期
}

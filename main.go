package main

import (
	"fmt"
	"net/http"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/pongo2"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"github.com/zileyuan/diagram/app"
)

func newMacaron() *macaron.Macaron {
	m := macaron.New()

	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("static"))
	m.Use(pongo2.Pongoer(pongo2.Options{
		Directory:  "views",
		IndentJSON: macaron.Env != macaron.PROD,
		IndentXML:  macaron.Env != macaron.PROD,
	}))
	m.Use(cache.Cacher())
	m.Use(session.Sessioner())
	return m
}

func newInstance() *macaron.Macaron {
	m := newMacaron()
	//路由跳转
	m.Get("/", app.DoIndex)
	m.Get("/index", app.DoIndex)
	m.Get("/page1", app.DoPage1)
	m.Get("/page2", app.DoPage2)
	//服务路由
	m.Any("/login", app.OnLogin)
	m.Any("/overview", app.OnOverview)
	m.Any("/cardtotal", app.OnCardTotal)
	m.Any("/customer", app.OnCustomer)
	m.Any("/cardtype", app.OnCardType)
	m.Any("/storecode", app.OnStoreCode)
	m.Any("/custtype", app.OnCustType)
	m.Any("/updcust", app.OnUpdCust)
	return m
}

func main() {
	fmt.Printf("-- Start %v Service --\n", app.AppName)
	m := newInstance()
	listenAddr := fmt.Sprintf("0.0.0.0:%d", app.HttpPort)
	fmt.Printf("-- Listen %v --\n", listenAddr)
	err := http.ListenAndServe(listenAddr, m)
	if err != nil {
		app.AppLog.Fatalf("Fail to start server: %v", err)
	}
}

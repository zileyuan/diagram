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
	m.Any("/overview", app.DoOverview)
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
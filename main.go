package main

import (
	"fmt"
	"net/http"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/pongo2"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
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
	m.Get("/", DoIndex)
	m.Get("/index", DoIndex)
	m.Any("/overview", DoOverview)
	return m
}

func main() {
	fmt.Printf("-- Start %v Service --\n", AppName)
	m := newInstance()
	listenAddr := fmt.Sprintf("0.0.0.0:%d", HttpPort)
	fmt.Printf("-- Listen %v --\n", listenAddr)
	err := http.ListenAndServe(listenAddr, m)
	if err != nil {
		AppLog.Fatalf("Fail to start server: %v", err)
	}
}

//func main() {
//	m := macaron.Classic()
//	m.Get("/", func() string {
//		return "Hello world!"
//	})
//	m.Run()
//}

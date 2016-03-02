package main

import (
	"fmt"
	"net/http"

	"github.com/go-macaron/macaron"
	"github.com/macaron-contrib/cache"
	"github.com/macaron-contrib/pongo2"
	"github.com/macaron-contrib/session"
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
	return m
}

func main() {
	fmt.Printf("-- Start %v Service --\n", AppName)
	m := newInstance()
	listenAddr := fmt.Sprintf("0.0.0.0:%d", HttpPort)
	fmt.Printf("-- Listen %v --\n", listenAddr)
	err := http.ListenAndServe(listenAddr, m)
	if err != nil {
		Log.Fatalf("Fail to start server: %v", err)
	}
}

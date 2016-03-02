package main

import "github.com/go-macaron/macaron"

func DoIndex(ctx *macaron.Context) string {
	return "Hello world!"
	//ctx.HTML(200, "show1")
}
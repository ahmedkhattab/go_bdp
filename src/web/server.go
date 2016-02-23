package main

import (
	"github.com/hoisie/web"
)

func hello(val string) string {
	return "hello " + val
}

func main() {
	web.Get("/hello/(.*)", hello)
	web.Run("localhost:9999")
}

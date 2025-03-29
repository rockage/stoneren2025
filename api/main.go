package main

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

func main() {

	crs := cors.New(cors.Options{ //crs相当于一个中间件，允许所有主机通过
		AllowedOrigins:   []string{"*"}, //
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})
	app := iris.New()
	app.Logger().SetLevel("debug")
	index := app.Party("/", crs) //所有请求先过crs中间件

	index.Get("/threads", Threads)
	index.Get("/comments", Comments)
	index.Get("/users", Users)

	// 监听 5050 端口
	app.Listen(":5050")
}

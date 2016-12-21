package server

import (
	"log"
	//	_ "./rfid"
	"github.com/changvvb/rfid/model"
	"github.com/iris-contrib/middleware/basicauth"
	"github.com/kataras/go-template/html"
	"github.com/kataras/iris"
)

type User struct {
	model.User
	UserName string
}

type Server struct {
	server *iris.Framework
}

func New() *Server {
	return &Server{server: iris.New()}
}

func (s *Server) Run() {

	server := s.server

	server.Config.IsDevelopment = true

	//server.StaticServe("./static", "/res")
	server.StaticServe("/home/changvvb/go/src/rfid/static", "/res")
	server.UseTemplate(html.New()).Directory("/home/changvvb/go/src/rfid/templates", ".html")

	server.Get("/", func(ctx *iris.Context) {
		userName := ctx.Session().GetString("username")
		if userName != "" {
			ctx.MustRender("index.html", User{UserName: userName})
			return
		}
		ctx.MustRender("index.html", nil)
	})

	server.Get("logintest/:userName", func(ctx *iris.Context) {
		userName := ctx.Param("userName")
		ctx.MustRender("index.html", struct{ UserName string }{UserName: userName})
	})

	server.Get("/login", func(ctx *iris.Context) {
		ctx.MustRender("login.html", nil)
	})

	server.Get("/logout", func(ctx *iris.Context) {
		printLog(ctx)
		ctx.Session().Clear()
		ctx.MustRender("index.html", nil)
	})

	server.Post("/login", func(ctx *iris.Context) {

		password := ctx.FormValueString("password")
		userName := ctx.FormValueString("username")
		log.Println(userName)
		log.Println(password)

		/* if password == "changvvb" && userName == "changvvb" { */
		// log.Println("login Success")
		// ctx.Session().Set("username", userName)
		// ctx.Redirect("/")
		/* } */

		if model.Exam(userName, password) {
			ctx.Session().Set("username", userName)
			ctx.Redirect("/")
			log.Println("login Success")
		} else {
			ctx.MustRender("login.html", struct{ LoginError bool }{true})
		}
	})

	authConfig := basicauth.Config{
		Users:      map[string]string{"changvvb": "changvvb"},
		Realm:      "Authorization Required", // if you don't set it it's "Authorization Required"
		ContextKey: "user",                   // if you don't set it it's "user"
		Expires:    0,                        //time.Duration(30) * time.Minute,
	}

	authentication := basicauth.New(authConfig)
	type Admins struct {
		UserName string
		Admins   []model.Admin
	}
	needAuth := server.Party("/admin", authentication)
	{
		needAuth.Get("/", func(ctx *iris.Context) {
			username := ctx.GetString("user") //  the Contextkey from the authConfig
			printLog(ctx, username)
			admins := []model.Admin{{AdminName: "changvvb", AdminPassword: "changvvb"}, {AdminName: "haha", AdminPassword: "hehe"}}
			ctx.MustRender("admin.html", Admins{username, admins})

		})

		needAuth.Get("/profile", func(ctx *iris.Context) {
			username := ctx.GetString("mycustomkey") //  the Contextkey from the authConfig
			ctx.Write("Hello authenticated user: %s from localhost:8080/secret/profile ", username)

		})

		needAuth.Get("/settings", func(ctx *iris.Context) {

		})

		needAuth.Get("/logout", func(ctx *iris.Context) {
			ctx.Session().Clear()
		})

		needAuth.Get("/savedata", func(ctx *iris.Context) {
			adminName := ctx.URLParam("name")
			password := ctx.URLParam("pass")
			log.Println("/admin/savedata", adminName, password)
			model.AddAdmin(adminName, password)
		})

	}

	server.Listen(":8080")
}

func printLog(ctx *iris.Context, v ...interface{}) {
	log.Println(ctx.PathString, v)
}

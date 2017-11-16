package server

import (
	"fmt"
	"log"
	"strconv"
	"time"
	//"time"
	"github.com/changvvb/rfid/rfid"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

type Card struct {
	// ID   string
	Name string //0 0
	Age  int    //0 1
	Sex  int    //0 2  0:woman,1:man
	Tel  string //0 3
}

var wsMap map[string]websocket.Connection

var ch chan ([]byte)
var flag bool
var cardID string

func init() {
	ch = make(chan []byte, 1)
	wsMap = make(map[string]websocket.Connection)
	// rfid.Connect("/dev/ttyUSB0")
	rfid.SearchCardCallBack = func(b []byte) {
		log.Println("search OK")
		if len(b) != 4 {
			cardID = "00000000"
			return
		}
		s := fmt.Sprintf("%X", b)
		if cardID != s {
			cardID = s
			for _, v := range wsMap {
				v.EmitMessage([]byte(s))
			}
		}
	}
	// rfid.AutoSearch14443()
}

func (c *Card) WriteToCard() {
	rfid.Auth14443(1)
	rfid.Write14443(1, 0, []byte(c.Name))
	rfid.Write14443(1, 1, []byte(c.Tel))
	rfid.Write14443(1, 2, []byte{byte(c.Sex), byte(c.Age)})

	// rfid.Auth14443(2)
	// rfid.Write14443(2, 0, []byte{byte(c.Sex)})
	// rfid.Write14443(2, 1, []byte(c.Tel))
}

type Server struct {
	server *iris.Application
}

func New() *Server {
	return &Server{server: iris.New()}
}

func (s *Server) Run() {

	server := s.server

	server.StaticServe("./templates", "/static")
	temp := iris.HTML("./templates", ".html")
	server.RegisterView(temp)

	// rfid.Connect("/dev/ttyUSB0")
	server.Use(func(ctx iris.Context) {

		if ctx.Path() == "/listdevices" || ctx.Path() == "/connect" || ctx.Path() == "/readpage" {
			ctx.Next()
			return
		}

		if rfid.BoolReady() {
			ctx.Next()
		} else {
			ctx.Text("Forbidden !!")
		}
	})

	server.Get("/readpage", func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		err := ctx.View("read.html")
		log.Println(err)
	})

	server.Get("/", func(ctx iris.Context) {
		ctx.WriteString("dsfasdfasdf")
	})

	server.Get("/readall", func(ctx iris.Context) {
		if rfid.BoolAutoSearch {
			rfid.StopAutoSearch14443()
			time.Sleep(time.Second / 5)
			defer rfid.AutoSearch14443()
		}
		var i, j byte
		for i = 0; i < 20; i++ {
			rfid.Auth14443(i)
			for j = 0; j < 3; j++ {
				b, _ := rfid.Read14443(i, j)
				ctx.Writef("%d,%d:%X\n", i, j, b)
			}
		}
	})

	server.Get("/read", func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		if rfid.BoolAutoSearch {
			rfid.StopAutoSearch14443()
			time.Sleep(time.Second / 5)
			defer rfid.AutoSearch14443()
		}
		card := Card{}
		rfid.Auth14443(1)
		card.Name = rfid.ReadString14443(1, 0)
		card.Tel = rfid.ReadString14443(1, 1)
		if b, err := rfid.Read14443(1, 2); err == nil {
			card.Age = int(b[1])
			card.Sex = int(b[0])
		}

		log.Printf("%+v\n", card)
		ctx.JSON(card)
	})

	server.Get("/write", func(ctx iris.Context) {
		ctx.View("write.html")
	})

	server.Get("/serialport", func(ctx iris.Context) {
		ctx.View("serialport.html")
	})

	server.Post("/write", func(ctx iris.Context) {
		if rfid.BoolAutoSearch {
			rfid.StopAutoSearch14443()
			time.Sleep(time.Second / 5)
			defer rfid.AutoSearch14443()
		}

		card := Card{}
		card.Name = ctx.FormValue("name")
		if age, err := strconv.Atoi(ctx.FormValue("age")); err != nil {
			return
		} else {
			card.Age = age
		}

		if sex, err := strconv.Atoi(ctx.FormValue("sex")); err != nil {
			return
		} else {
			card.Sex = sex
		}
		card.Tel = ctx.FormValue("tel")
		log.Printf("%+v\n", card)
		card.WriteToCard()

	})

	server.Get("/autosearch", func(ctx iris.Context) {
		rfid.AutoSearch14443()
	})

	server.Get("/stopautosearch", func(ctx iris.Context) {
		rfid.StopAutoSearch14443()
	})

	server.Get("/connect", func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		port := ctx.URLParam("port")
		log.Println("port", port)
		err := rfid.Connect(port)
		if err == nil {
			ctx.Writef("Connect successfully")
		} else {
			ctx.Writef("Error: %s", err.Error())
		}

	})

	server.Post("/connect", func(ctx iris.Context) {
	})

	server.Get("/listdevices", func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		s := ctx.URLParam("flag")
		if s == "1" {
			ctx.JSON(rfid.ListDevice(true))
		} else {
			ctx.JSON(rfid.ListDevice(false))
		}
	})

	server.Get("/websocket", func(ctx iris.Context) {
		ctx.View("websocket.html")
		rfid.AutoSearch14443()
	})

	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})

	server.Get("/ws", ws.Handler())

	ws.OnConnection(func(c websocket.Connection) {

		wsMap[c.ID()] = c

		c.Join("room")

		c.On("chat", func(message string) {
			log.Println("Websocket", message)
		})

		c.OnMessage(func(b []byte) {
			log.Println(string(b))
			// c.EmitMessage(b)
		})

		c.OnDisconnect(func() {
			log.Println("delete  ...... ")
			delete(wsMap, c.ID())
		})

		c.OnError(func(err string) {
			log.Println("Error: ", err)
		})

	})

	server.Run(iris.Addr(":8080"))
}

func printLog(ctx iris.Context, v ...interface{}) {
	log.Println(ctx.Path, v)
}

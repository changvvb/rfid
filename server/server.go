package server

import (
	"fmt"
	"iotWeb/model"
	"log"
	"strconv"
	"time"
	//"time"

	"github.com/changvvb/rfid/rfid"
	"github.com/kataras/iris"
	//	_ "./rfid"
)

type User struct {
	model.User
	UserName string
}

type Card struct {
	ID   string
	Name string //0 0
	Age  int    //0 1
	Sex  int    //0 2  0:woman,1:man
	Tel  string //0 3
}

var wsMap map[string]iris.WebsocketConnection

var ch chan ([]byte)
var flag bool
var cardID string

func init() {
	ch = make(chan []byte, 1)
	wsMap = make(map[string]iris.WebsocketConnection)
	rfid.Connect("/dev/ttyUSB0")
	rfid.SearchCardCallBack = func(b []byte) {
		log.Println("search OK")
		if len(b) != 4 {
			cardID = "00000000"
			return
		}
		log.Println("length:---------->", len(b))
		s := fmt.Sprintf("%X", b)
		if cardID != s {
			cardID = s
			for _, v := range wsMap {
				v.EmitMessage([]byte(s))
			}
		}
	}

	rfid.AutoSearch14443()
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
	server *iris.Framework
}

func New() *Server {
	return &Server{server: iris.New()}
}

func (s *Server) Run() {

	server := s.server
	server.Config.IsDevelopment = true

	rfid.Connect("/dev/ttyUSB0")
	server.UseFunc(func(ctx *iris.Context) {
		if rfid.BoolReady() {
			ctx.Next()
		} else {
			ctx.Text(403, "Forbidden !!")
		}
	})

	server.Get("/read", func(ctx *iris.Context) {
		rfid.StopAutoSearch14443()
		time.Sleep(time.Second / 2)
		defer rfid.AutoSearch14443()
		card := Card{}
		rfid.Auth14443(1)
		card.Name = rfid.ReadString14443(1, 0)
		card.Tel = rfid.ReadString14443(1, 1)
		if b, err := rfid.Read14443(1, 2); err == nil {
			card.Age = int(b[1])
			card.Sex = int(b[0])
		}
		/*
			rfid.Auth14443(2)
			if b, err := rfid.Read14443(2, 0); err == nil {

				card.Sex = int(b[0])
			}
			card.Tel = rfid.ReadString14443(2, 1)
		*/
		log.Printf("%+v\n", card)
		ctx.JSON(iris.StatusOK, card)
	})

	server.Get("/write", func(ctx *iris.Context) {
		// html := `<html>
		// <body>
		// <form method="POST" action="/write">
		// <input type="text" name="name"></input>
		// <input type="text" name="age"></input>
		// <input type="text" name="sex"></input>
		// <input type="text" name="tel"></input>
		// <input type="submit" ></input>
		// </form>
		// </body>
		// </html>
		// `

		// html = `<html><head></head><body></body></html>`

		// ctx.RenderTemplateSource(iris.StatusOK, html, nil)
		ctx.Render("write.html", nil)
	})

	server.Post("/write", func(ctx *iris.Context) {
		rfid.StopAutoSearch14443()
		time.Sleep(time.Second / 2)
		defer rfid.AutoSearch14443()

		card := Card{}
		card.Name = ctx.FormValueString("name")
		if age, err := strconv.Atoi(ctx.FormValueString("age")); err != nil {
			return
		} else {
			card.Age = age
		}

		if sex, err := strconv.Atoi(ctx.FormValueString("sex")); err != nil {
			return
		} else {
			card.Sex = sex
		}
		card.Tel = ctx.FormValueString("tel")
		log.Printf("%+v\n", card)
		card.WriteToCard()

	})

	server.Get("/connect", func(ctx *iris.Context) {
		html := `<html>
		<body>
		<form method="POST" action="/write">
		<input type="text" name="name"></input>
		<input type="text" name="age"></input>
		<input type="text" name="sex"></input>
		<input type="text" name="tel"></input>
		<input type="submit" ></input>
		</form>
		</body>
		</html>	
		`
		ctx.RenderTemplateSource(iris.StatusOK, html, nil)

	})

	server.Post("/connect", func(ctx *iris.Context) {
	})

	server.Get("/listdevices", func(ctx *iris.Context) {
		s := ctx.URLParam("flag")
		if s == "1" {
			ctx.JSON(iris.StatusOK, rfid.ListDevice(true))
		} else {
			ctx.JSON(iris.StatusOK, rfid.ListDevice(false))
		}
	})

	server.Get("/websocket", func(ctx *iris.Context) {
		ctx.MustRender("websocket.html", nil)
		rfid.AutoSearch14443()
	})

	server.Config.Websocket.Endpoint = "/ws"

	server.Websocket.OnConnection(func(c iris.WebsocketConnection) {

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

		//     go func() {
		//         card := ""
		//         for {
		//             if flag && rfid.BoolAutoSearch {
		//                 s := fmt.Sprintf("%X", cardID)
		//                 b := []byte(s)
		//                 if card != s {
		//                     card = s
		//                     c.EmitMessage(b)
		//                     // rfid.StopAutoSearch14443()
		//                     flag = false
		//                     time.Sleep(time.Second / 10)
		//                 }
		//             }
		//         }
		//     }()
	})

	server.Listen(":8080")
}

func printLog(ctx *iris.Context, v ...interface{}) {
	log.Println(ctx.PathString, v)
}

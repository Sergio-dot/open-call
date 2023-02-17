package server

import (
	"encoding/gob"
	"flag"
	"log"
	"os"
	"time"

	"github.com/Sergio-dot/open-call/internal/auth"
	"github.com/Sergio-dot/open-call/internal/handlers"
	"github.com/Sergio-dot/open-call/internal/models"
	w "github.com/Sergio-dot/open-call/pkg/webrtc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// application flags
var (
	addr = flag.String("addr", ":"+os.Getenv("PORT"), "")
	cert = flag.String("cert", "", "")
	key  = flag.String("key", "", "")
)

func Run() error {
	gob.Register(models.User{})
	gob.Register(time.Time{})

	flag.Parse()

	if *addr == ":" {
		*addr = ":8080"
	}

	// initialize a logger
	dbLogger := log.New(os.Stdout, "\r\n", 0)

	// connect to PostgreSQL database
	db, err := gorm.Open("postgres", "user=postgres password=root dbname=opencall sslmode=disable")
	if err != nil {
		return err
	}
	defer func(db *gorm.DB) {
		err = db.Close()
		if err != nil {
			log.Fatal("Couldn't close database properly. Killing application")
		}
	}(db)
	db.LogMode(true)
	db.SetLogger(dbLogger)

	// pass database connection to handlers
	handlers.DB = db

	// setup fiber configs
	engine := html.New("./views", ".tmpl")
	app := fiber.New(fiber.Config{Views: engine})
	store := session.New()

	// middlewares
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(auth.Session(*store))

	// pass session manager to handlers
	handlers.Store = store

	// non-protected routes
	app.Get("/", handlers.Home)
	app.Post("/login", handlers.Login)
	app.Post("/signup", handlers.SignUp)
	app.Get("/login/google", handlers.GoogleLogin)
	app.Get("/login/google/callback", handlers.GoogleCallback)

	// protected routes
	group := app.Group("/", auth.Authentication)
	group.Get("/logout", handlers.Logout)
	group.Get("/dashboard", handlers.Dashboard)
	group.Get("/user/update/:id", handlers.UpdateUser)
	group.Get("/room/create", handlers.RoomCreate)
	group.Get("/room/:uuid", handlers.Room)
	group.Get("/room/:uuid/websocket", websocket.New(handlers.RoomWebsocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))
	group.Get("/room/:uuid/chat", handlers.RoomChat)
	group.Get("/room/:uuid/chat/websocket", websocket.New(handlers.RoomChatWebsocket))
	group.Get("/room/:uuid/viewer/websocket", websocket.New(handlers.RoomViewerWebsocket))
	group.Get("/stream/:suuid", handlers.Stream)
	group.Get("/stream/:suuid/websocket", websocket.New(handlers.StreamWebsocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))
	group.Get("/stream/:suuid/chat/websocket", websocket.New(handlers.StreamChatWebsocket))
	group.Get("/stream/:suuid/viewer/websocket", websocket.New(handlers.StreamViewerWebsocket))

	// static files server
	app.Static("/", "./assets")

	// initialize room and stream maps
	w.Rooms = make(map[string]*w.Room)
	w.Streams = make(map[string]*w.Room)

	go dispatchKeyFrames()

	if *cert != "" {
		return app.ListenTLS(*addr, *cert, *key)
	}
	return app.Listen(*addr)
}

func dispatchKeyFrames() {
	for range time.NewTicker(time.Second * 3).C {
		for _, room := range w.Rooms {
			room.Peers.DispatchKeyFrame()
		}
	}
}

package main

import (
	"context"
	"flag"
	"log"

	"github.com/Jimbo8702/goreservation/api"
	"github.com/Jimbo8702/goreservation/api/middleware"
	"github.com/Jimbo8702/goreservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config =fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
    },
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		// store initialization
		userStore 	 = db.NewMongoUserStore(client)
		hotelStore 	 = db.NewMongoHotelStore(client)
		roomStore 	 = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
		store = &db.Store{
			User: userStore,
			Hotel: hotelStore,
			Room: roomStore,
			Booking: bookingStore,
		}
		// handler initialization
		authHandler	 = api.NewAuthHandler(userStore)
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		roomHandler  = api.NewRoomHandler(store)
		//app initialization
		app 		 = fiber.New(config)
		auth 		 = app.Group("api")
		apiv1 		 = app.Group("api/v1", middleware.JWTAuthentication(userStore))
	)

	// auth 
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Versioned API routes
	// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	//booking handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	app.Listen(*listenAddr)
}



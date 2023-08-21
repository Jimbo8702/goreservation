package main

import (
	"context"
	"log"
	"os"

	"github.com/Jimbo8702/goreservation/api"
	"github.com/Jimbo8702/goreservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//configuration
// 1. MongoDB endpoint
// 2. ListenAddress of our HTTP server
// 3. JWT secret
// 4. MongoDBName

var config = fiber.Config{
    ErrorHandler: api.ErrorHandler,
}

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
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
		authHandler	   = api.NewAuthHandler(userStore)
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		//app initialization
		app 		 = fiber.New(config)
		auth 		 = app.Group("api")
		apiv1 		 = app.Group("api/v1", api.JWTAuthentication(userStore))
		admin 		 = apiv1.Group("/admin", api.AdminAuth)
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

	//rooms handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	//bookings handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	// admin.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
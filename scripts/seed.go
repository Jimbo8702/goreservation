package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Jimbo8702/goreservation/api"
	"github.com/Jimbo8702/goreservation/db"
	"github.com/Jimbo8702/goreservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	var err error 
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := db.Store{
		User: db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room: db.NewMongoRoomStore(client, hotelStore),
		Hotel: hotelStore,
	}

	user := fixtures.AddUser(&store, "james", "foo", false)
	fmt.Println("james ->", api.CreateTokenFromUser(user))

	admin := fixtures.AddUser(&store, "admin", "admin", true)
	fmt.Println("admin ->",api.CreateTokenFromUser(admin))

	hotel := fixtures.AddHotel(&store, "Bellucia", "France", 3, nil)
	fmt.Println("hotel ->", hotel.ID)

	room := fixtures.AddRoom(&store, "large", true, 100.99, hotel.ID)
	fmt.Println("room ->", room.ID)

	booking := fixtures.AddBooking(&store, user.ID, room.ID, time.Now(), time.Now().AddDate(0,0,5))
	fmt.Println("booking ->", booking.ID)
}
